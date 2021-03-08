package api

import (
	"encoding/json"
	"fmt"

	"github.com/YaleSpinup/ecr-api/iam"
	log "github.com/sirupsen/logrus"
)

// orgTagAccessPolicy generates the org tag conditional policy to be passed inline when assuming a role
func orgTagAccessPolicy(org string) (string, error) {
	log.Debugf("generating org policy document")

	policy := iam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iam.StatementEntry{
			{
				Effect:   "Allow",
				Action:   []string{"*"},
				Resource: []string{"*"},
				Condition: iam.Condition{
					"StringEquals": iam.ConditionStatement{
						"aws:ResourceTag/spinup:org": []string{org},
					},
				},
			},
		},
	}

	j, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func (s *server) repositoryUserCreatePolicy() (string, error) {
	policy := &iam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iam.StatementEntry{
			{
				Sid:    "CreateRepositoryUser",
				Effect: "Allow",
				Action: []string{
					"iam:CreatePolicy",
					"iam:UntagUser",
					"iam:GetPolicyVersion",
					"iam:AddUserToGroup",
					"iam:GetPolicy",
					"iam:ListAttachedGroupPolicies",
					"iam:ListGroupPolicies",
					"iam:AttachGroupPolicy",
					"iam:GetUser",
					"iam:CreatePolicyVersion",
					"iam:CreateUser",
					"iam:GetGroup",
					"iam:CreateGroup",
					"iam:TagUser",
				},
				Resource: []string{
					"arn:aws:iam::*:group/*",
					fmt.Sprintf("arn:aws:iam::*:policy/spinup/%s/*", s.org),
					fmt.Sprintf("arn:aws:iam::*:user/spinup/%s/*", s.org),
				},
			},
			{
				Sid:    "ListRepositoryUserPolicies",
				Effect: "Allow",
				Action: []string{
					"iam:ListPolicies",
				},
				Resource: []string{"*"},
			},
		},
	}

	j, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func (s *server) repositoryUserDeletePolicy() (string, error) {
	policy := &iam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iam.StatementEntry{
			{
				Sid:    "DeleteRepositoryUser",
				Effect: "Allow",
				Action: []string{
					"iam:DeleteAccessKey",
					"iam:RemoveUserFromGroup",
					"iam:ListAccessKeys",
					"iam:ListGroupsForUser",
					"iam:DeleteUser",
					"iam:GetUser",
				},
				Resource: []string{
					fmt.Sprintf("arn:aws:iam::*:user/spinup/%s/*", s.org),
					fmt.Sprintf("arn:aws:iam::*:group/spinup/%s/SpinupECRAdminGroup-%s", s.org, s.org),
				},
			},
		},
	}

	j, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

func (s *server) repositoryUserUpdatePolicy() (string, error) {
	policy := &iam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iam.StatementEntry{
			{
				Sid:    "UpdateRepositoryUser",
				Effect: "Allow",
				Action: []string{
					"iam:UntagUser",
					"iam:DeleteAccessKey",
					"iam:RemoveUserFromGroup",
					"iam:TagUser",
					"iam:CreateAccessKey",
					"iam:ListAccessKeys",
				},
				Resource: []string{
					fmt.Sprintf("arn:aws:iam::*:user/spinup/%s/*", s.org),
					fmt.Sprintf("arn:aws:iam::*:group/spinup/%s/SpinupECRAdminGroup-%s", s.org, s.org),
				},
			},
		},
	}

	j, err := json.Marshal(policy)
	if err != nil {
		return "", err
	}

	return string(j), nil
}

// repositoryPolicy accepts a list of groups and returns the policy to allow ecr access
// for resources in the same org/group as well as any passed groups
func repositoryPolicy(groups []string) (string, error) {
	groupConditions := append([]string{"${aws:ResourceTag/spinup:spaceid}"}, groups...)

	log.Debugf("generating policy text from groups %+v", groups)

	policy := iam.PolicyDocument{
		Version: "2012-10-17",
		Statement: []iam.StatementEntry{
			{
				Sid:    "AllowPullImagesFromSpaceAndOrg",
				Effect: "Allow",
				Action: []string{
					"ecr:GetAuthorizationToken",
					"ecr:BatchCheckLayerAvailability",
					"ecr:GetDownloadUrlForLayer",
					"ecr:BatchGetImage",
				},
				Principal: "*",
				Condition: iam.Condition{
					"StringEqualsIgnoreCase": iam.ConditionStatement{
						"aws:PrincipalTag/spinup:org":     "${aws:ResourceTag/spinup:org}",
						"aws:PrincipalTag/spinup:spaceid": groupConditions,
					},
				},
			},
		},
	}

	policyDoc, err := json.Marshal(policy)
	if err != nil {
		log.Errorf("failed to generate repository policy documentfor %s", err)
		return "", err
	}

	log.Debugf("returning policy document from groups: %s", string(policyDoc))

	return string(policyDoc), nil
}

// repositoryGroupsFromPolicy returns the list of groups from the repository policy string
func repositoryGroupsFromPolicy(policy string) ([]string, error) {
	if policy == "" {
		return []string{}, nil
	}

	log.Debugf("getting groups from policy text: %s", policy)

	policyDoc := iam.PolicyDocument{}
	if err := json.Unmarshal([]byte(policy), &policyDoc); err != nil {
		return nil, err
	}

	groups := []string{}

	// for all of the statements in our policy
	for _, statement := range policyDoc.Statement {
		// if we aren't dealing with the policy we set, continue to the next statement
		if statement.Sid != "AllowPullImagesFromSpaceAndOrg" {
			continue
		}

		conditionStatement, ok := statement.Condition["StringEqualsIgnoreCase"]
		if !ok {
			continue
		}

		for k, v := range conditionStatement {
			// look for the condition on the spaceid tag
			if k != "aws:PrincipalTag/spinup:spaceid" {
				log.Debugf("resource policy condition tag key '%s' is not 'aws:PrincipalTag/spinup:spaceid', continuing", k)
				continue
			}

			// should be a list of strings unless there are no
			// additional groups added to the list
			list, ok := v.([]interface{})
			if !ok {
				log.Debugf("resource policy condition value '%+v' is not a list, continuing", v)
				continue
			}

			// collect the spaceid tags and add to the list of groups
			for _, g := range list {
				// values should all be strings
				gv, ok := g.(string)
				if !ok {
					log.Warnf("tag value '%v' is not a string", g)
					continue
				}

				// ignore the "same space" group
				if gv == "${aws:ResourceTag/spinup:spaceid}" {
					continue
				}

				groups = append(groups, gv)
			}
		}
	}

	log.Debugf("returning groups list: %v", groups)

	return groups, nil
}
