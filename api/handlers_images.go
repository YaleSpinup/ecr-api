package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/ecr-api/ecr"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// RepositoriesImageListHandler is the http handler for listing images in a repository
func (s *server) RepositoriesImageListHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]
	group := vars["group"]

	repository := fmt.Sprintf("%s/%s", group, name)

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

	service := ecr.New(
		ecr.WithSession(session.Session),
	)

	images, err := service.GetImages(r.Context(), repository)
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(images)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to marshal response from the ecr service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// RepositoriesImageTagShowHandler returns information about an image tag, notably the detailed scan findings
func (s *server) RepositoriesImageTagShowHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]
	group := vars["group"]
	tag := vars["tag"]

	repository := fmt.Sprintf("%s/%s", group, name)

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

	service := ecr.New(
		ecr.WithSession(session.Session),
	)

	findings, err := service.GetImageScanFindings(r.Context(), repository, tag)
	if err != nil {
		// Check if the error is because scan findings don't exist yet
		if aerr, ok := err.(*apierror.Error); ok && aerr.Code == apierror.ErrNotFound {
			// Return empty scan findings with a status indicating no scan available
			emptyScanResponse := map[string]interface{}{
				"imageScanStatus": map[string]string{
					"status":      "NO_SCAN_AVAILABLE",
					"description": "Image scan findings are not available for this image",
				},
				"findingSeverityCounts": map[string]int{},
				"findings":              []interface{}{},
			}
			j, err := json.Marshal(emptyScanResponse)
			if err != nil {
				handleError(w, errors.Wrap(err, "unable to marshal empty scan response"))
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(j)
			return
		}
		handleError(w, err)
		return
	}

	j, err := json.Marshal(findings)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to marshal response from the ecr service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// RepositoriesImageTagDeleteHandler deletes an image tag
func (s *server) RepositoriesImageTagDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]
	group := vars["group"]
	tag := vars["tag"]

	repository := fmt.Sprintf("%s/%s", group, name)

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)
	policy, err := s.repositoryImageDeletePolicy(account, repository)
	if err != nil {
		handleError(w, apierror.New(apierror.ErrInternalError, "failed to generate policy", err))
		return
	}

	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		policy,
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	service := ecr.New(
		ecr.WithSession(session.Session),
	)

	output, err := service.DeleteImageTag(r.Context(), repository, tag)
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(output)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to marshal response from the ecr service"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
