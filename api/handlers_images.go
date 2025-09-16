package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/ecr-api/ecr"
	awsecr "github.com/aws/aws-sdk-go/service/ecr"
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

	// First, get the image details
	imageID := &awsecr.ImageIdentifier{ImageTag: &tag}
	images, err := service.GetImages(r.Context(), repository, imageID)
	if err != nil {
		handleError(w, err)
		return
	}

	// Create response structure
	type ImageTagResponse struct {
		ImageDetail   *awsecr.ImageDetail        `json:"imageDetail,omitempty"`
		ScanFindings  *awsecr.ImageScanFindings  `json:"scanFindings,omitempty"`
		ScanError     string                  `json:"scanError,omitempty"`
	}

	response := ImageTagResponse{}

	// Add image detail if found
	if len(images) > 0 {
		response.ImageDetail = images[0]
	}

	// Try to get scan findings, but don't fail if they're not available
	findings, err := service.GetImageScanFindings(r.Context(), repository, tag)
	if err != nil {
		// Log the error but don't fail the request
		response.ScanError = fmt.Sprintf("Unable to retrieve scan findings: %v", err)
	} else {
		response.ScanFindings = findings
	}

	j, err := json.Marshal(response)
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
