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
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, nil))
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
