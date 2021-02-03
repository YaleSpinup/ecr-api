package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/ecr-api/iam"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// UsersCreateHandler creates a user that can access the repository.  It first checks if the shared policy exists
// in the account and creates it/updates it as needed.
func (s *server) UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	name := vars["name"]
	group := vars["group"]

	req := RepositoryUserCreateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into create user input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	if req.UserName == "" {
		handleError(w, apierror.New(apierror.ErrBadRequest, "username is required", nil))
		return
	}

	if len(req.Groups) == 0 {
		handleError(w, apierror.New(apierror.ErrBadRequest, "at least 1 group is required", nil))
		return
	}

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	// IAM doesn't support resource tags, so we can't pass the s.orgPolicy here
	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	orch := &iamOrchestrator{
		client: iam.New(
			iam.WithSession(session.Session),
		),
		org: s.org,
	}

	groupName, err := orch.prepareAccount(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	out, err := orch.repositoryUserCreate(r.Context(), name, group, groupName, &req)
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(out)
	if err != nil {
		log.Errorf("cannot marshal reasponse(%v) into JSON: %s", out, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UserListHandler lists the users within a prefix.  Specific repositories, group name and
// all managed users with an empty prefix are acceptable and expand the scope of the list.
func (s *server) UsersListHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	group := vars["group"]
	name := vars["name"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	// IAM doesn't support resource tags, so we can't pass the s.orgPolicy here
	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	orch := &iamOrchestrator{
		client: iam.New(
			iam.WithSession(session.Session),
		),
		org: s.org,
	}

	output, err := orch.listRepositoryUsers(r.Context(), group, name)
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(output)
	if err != nil {
		log.Errorf("cannot marshal reasponse(%v) into JSON: %s", output, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UsersShowHandler gets the information about a repository user.
func (s *server) UsersShowHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	group := vars["group"]
	name := vars["name"]
	user := vars["user"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	// IAM doesn't support resource tags, so we can't pass the s.orgPolicy here
	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	orch := &iamOrchestrator{
		client: iam.New(
			iam.WithSession(session.Session),
		),
		org: s.org,
	}

	output, err := orch.getRepositoryUser(r.Context(), group, name, user)
	if err != nil {
		handleError(w, err)
		return
	}

	j, err := json.Marshal(output)
	if err != nil {
		log.Errorf("cannot marshal reasponse(%v) into JSON: %s", output, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UsersUpdateHandler updates a repository user
func (s *server) UsersUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	group := vars["group"]
	name := vars["name"]
	userName := vars["user"]

	req := RepositoryUserUpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		msg := fmt.Sprintf("cannot decode body into update repository user input: %s", err)
		handleError(w, apierror.New(apierror.ErrBadRequest, msg, err))
		return
	}

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	// IAM doesn't support resource tags, so we can't pass the s.orgPolicy here
	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	orch := &iamOrchestrator{
		client: iam.New(
			iam.WithSession(session.Session),
		),
		org: s.org,
	}

	resp, err := orch.repositoryUserUpdate(r.Context(), name, group, userName, &req)
	if err != nil {
		handleError(w, errors.Wrap(err, "failed to update repository user"))
		return
	}

	j, err := json.Marshal(resp)
	if err != nil {
		handleError(w, errors.Wrap(err, "unable to marshal response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// UsersDeleteHandler deletes a repository user
func (s *server) UsersDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w = LogWriter{w}
	vars := mux.Vars(r)
	account := vars["account"]
	group := vars["group"]
	name := vars["name"]
	userName := vars["user"]

	role := fmt.Sprintf("arn:aws:iam::%s:role/%s", account, s.session.RoleName)

	// IAM doesn't support resource tags, so we can't pass the s.orgPolicy here
	session, err := s.assumeRole(
		r.Context(),
		s.session.ExternalID,
		role,
		"",
	)
	if err != nil {
		msg := fmt.Sprintf("failed to assume role in account: %s", account)
		handleError(w, apierror.New(apierror.ErrForbidden, msg, nil))
		return
	}

	orch := &iamOrchestrator{
		client: iam.New(
			iam.WithSession(session.Session),
		),
		org: s.org,
	}

	if err := orch.repositoryUserDelete(r.Context(), name, group, userName); err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
