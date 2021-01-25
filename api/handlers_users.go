package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/YaleSpinup/apierror"
	"github.com/YaleSpinup/ecr-api/iam"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type iamOrchestrator struct {
	client iam.IAM
	org    string
}

var ecrAdminPolicyDoc string
var EcrAdminPolicy = iam.PolicyDocument{
	Version: "2012-10-17",
	Statement: []iam.StatementEntry{
		{
			Sid:    "AllowActionsOnRepositoriesInSpaceAndOrg",
			Effect: "Allow",
			Action: []string{
				"ecr:PutLifecyclePolicy",
				"ecr:PutImageTagMutability",
				"ecr:DescribeImageScanFindings",
				"ecr:GetDownloadUrlForLayer",
				"ecr:GetAuthorizationToken",
				"ecr:UploadLayerPart",
				"ecr:BatchDeleteImage",
				"ecr:ListImages",
				"ecr:DeleteLifecyclePolicy",
				"ecr:PutImage",
				"ecr:BatchGetImage",
				"ecr:CompleteLayerUpload",
				"ecr:DescribeImages",
				"ecr:DeleteRegistryPolicy",
				"ecr:InitiateLayerUpload",
				"ecr:BatchCheckLayerAvailability",
			},
			Resource: "*",
			Condition: iam.Condition{
				"StringEquals": iam.ConditionStatement{
					"aws:ResourceTag/spinup:org":     "${aws:PrincipalTag/spinup:org}",
					"aws:ResourceTag/spinup:spaceid": "${aws:PrincipalTag/spinup:spaceid}",
					"aws:ResourceTag/Name":           "${aws:PrincipalTag/Name}",
				},
			},
		},
		{
			Sid:      "AllowDockerLogin",
			Effect:   "Allow",
			Action:   []string{"ecr:GetAuthorizationToken"},
			Resource: "*",
		},
	},
}

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

	// start prep account for user management

	path := fmt.Sprintf("/spinup/%s/", s.org)

	policyName := fmt.Sprintf("SpinupECRAdminPolicy-%s", s.org)
	policyArn, err := orch.userCreatePolicyIfMissing(r.Context(), policyName, path)
	if err != nil {
		handleError(w, err)
		return
	}

	groupName := fmt.Sprintf("SpinupECRAdminGroup-%s", s.org)
	if err := orch.userCreateGroupIfMissing(r.Context(), groupName, path, policyArn); err != nil {
		handleError(w, err)
		return
	}

	out, err := orch.repositoryUserCreate(r.Context(), name, group, groupName, req)
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

	service := iam.New(
		iam.WithSession(session.Session),
	)

	path := fmt.Sprintf("/spinup/%s", s.org)

	if group != "" {
		path = path + fmt.Sprintf("/%s", group)
	}

	if name != "" {
		path = path + fmt.Sprintf("/%s", name)
	}

	users, err := service.ListUsers(r.Context(), path)
	if err != nil {
		handleError(w, err)
		return
	}

	ps := strings.Split(path, "/")
	if len(ps) > 2 {
		prefix := fmt.Sprintf("%s-%s-", ps[len(ps)-2], ps[len(ps)-1])

		trimmed := make([]string, 0, len(users))
		for _, u := range users {
			log.Debugf("trimming prefix '%s' from username %s", prefix, u)
			u = strings.TrimPrefix(u, prefix)
			trimmed = append(trimmed, u)
		}
		users = trimmed
	}

	j, err := json.Marshal(users)
	if err != nil {
		log.Errorf("cannot marshal reasponse(%v) into JSON: %s", users, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

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

	service := iam.New(
		iam.WithSession(session.Session),
	)

	path := fmt.Sprintf("/spinup/%s/%s/%s/", s.org, group, name)
	userName := fmt.Sprintf("%s-%s-%s", group, name, user)

	iamUser, err := service.GetUserWithPath(r.Context(), path, userName)
	if err != nil {
		handleError(w, err)
		return
	}

	keys, err := service.ListAccessKeys(r.Context(), userName)
	if err != nil {
		handleError(w, err)
		return
	}

	output := repositoryUserResponseFromIAM(iamUser, keys)

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
