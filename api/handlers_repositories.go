package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/ecr-api/ecr"
	"github.com/YaleSpinup/ecr-api/iam"
	"github.com/YaleSpinup/ecr-api/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// RepositoriesCreateHandler is the http handler for creating a repository
func (s *server) RepositoriesCreateHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	group := vars["group"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
		"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, err))
		return
	}

	req := RepositoryCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into create repository input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	orch := newEcrOrchestrator(
		ecr.New(ecr.WithSession(session.Session)),
		s.org,
	)

	resp, err := orch.repositoryCreate(r.Context(), account, group, &req)
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to create repository"))
		return
	}

	j, err := json.Marshal(resp)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to marshal response from the ecr service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// RepositoriesListHandler is the http handler for listing repositories
func (s *server) RepositoriesListHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
		"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
		"arn:aws:iam::aws:policy/ResourceGroupsandTagEditorReadOnlyAccess",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	var repos []string
	if group, ok := vars["group"]; ok {
		service := resourcegroupstaggingapi.New(
			resourcegroupstaggingapi.WithSession(session.Session),
		)

		// build up tag filters starting with the org
		tagFilters := []*resourcegroupstaggingapi.TagFilter{
			{
				Key:   "spinup:org",
				Value: []string{s.org},
			},
			{
				Key:   "spinup:spaceid",
				Value: []string{group},
			},
		}

		out, err := service.GetResourcesWithTags(r.Context(), []string{"ecr"}, tagFilters)
		if err != nil {
			handleError(w, errors.Wrap(err, "failed to create repository"))
			return
		}

		log.Debugf("got output from resourcegroups tagging api %s", awsutil.Prettify(out))

		repos = make([]string, 0, len(out))
		for _, repo := range out {
			a, err := arn.Parse(aws.StringValue(repo.ResourceARN))
			if err != nil {
				msg := fmt.Sprintf("failed to parse ARN %s: %s", repo, err)
				handleError(w, errors.Wrap(err, msg))
				return
			}

			prefix := fmt.Sprintf("repository/%s/", group)
			rid := strings.TrimPrefix(a.Resource, prefix)
			repos = append(repos, rid)
		}
	} else {
		service := ecr.New(
			ecr.WithSession(session.Session),
		)

		var err error
		repos, err = service.ListRepositories(r.Context())
		if err != nil {
			handleError(w, err)
			return
		}
	}

	j, err := json.Marshal(repos)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to marshal response from the ecr service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// RepositoriesShowHandler gets the details about an individual repository
func (s *server) RepositoriesShowHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]
	group := vars["group"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		s.orgPolicy,
		"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	orch := newEcrOrchestrator(
		ecr.New(ecr.WithSession(session.Session)),
		s.org,
	)

	resp, err := orch.repositoryDetails(r.Context(), account, group, name)
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to get repository details"))
		return
	}

	j, err := json.Marshal(resp)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to marshal response from the ecr service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// RepositoriesUpdateHandler handles updating a repository
func (s *server) RepositoriesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]
	group := vars["group"]

	req := RepositoryUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into update repository input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		s.orgPolicy,
		"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
	}

	orch := newEcrOrchestrator(
		ecr.New(ecr.WithSession(session.Session)),
		s.org,
	)

	resp, err := orch.repositoryUpdate(r.Context(), account, group, name, &req)
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to update repository"))
		return
	}

	j, err := json.Marshal(resp)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to marshal response from the ecr service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// RepositoriesDeleteHandler deletes a repository
func (s *server) RepositoriesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]
	group := vars["group"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	policy, err := s.repositoryDeletePolicy(s.org)
	if err != nil {
		handleError(w, apierror.New(apierror.ErrInternalError, "failed to generate policy", err))
		return
	}

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		policy,
		"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
	}

	orch := newEcrOrchestrator(
		ecr.New(ecr.WithSession(session.Session)),
		s.org,
	)

	resp, err := orch.repositoryDelete(r.Context(), account, group, name)
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to create repository"))
		return
	}

	iamOrch := newIamOrchestrator(
		iam.New(iam.WithSession(session.Session)),
		s.org,
	)

	users, err := iamOrch.repositoryUserDeleteAll(r.Context(), name, group)
	if err != nil {
		handleError(w, err)
		return
	}

	response := struct {
		RepositoryResponse
		Users []string
	}{*resp, users}

	j, err := json.Marshal(response)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to marshal response from the ecr service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// ScanRepositoriesListHandler Scans all repositories
func (s *server) ScanRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
		"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	service := ecr.New(
		ecr.WithSession(session.Session),
	)

	repositories, err := service.ListRepositories(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	for _, repository := range repositories {
		images, err := service.ListImages(r.Context(), repository)
		if err != nil {
			handleError(w, err)
			return
		}
		for _, image := range images {
			err = service.ScanImage(r.Context(), image, repository)
			if err != nil {
				handleError(w, err)
				return
			}
		}

	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("scan initiated"))
}
