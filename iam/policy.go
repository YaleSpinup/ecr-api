package iam

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/YaleSpinup/apierror"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	log "github.com/sirupsen/logrus"
)

type PolicyDocument struct {
	// 2012-10-17 or 2008-10-17 old policies, do NOT use this for new policies
	Version   string           `json:"Version"`
	Id        string           `json:"Id,omitempty"`
	Statement []StatementEntry `json:"Statement"`
}

type StatementEntry struct {
	Sid          string    `json:"Sid,omitempty"`          // statement ID, service specific
	Effect       string    `json:"Effect"`                 // Allow or Deny
	Principal    Principal `json:"Principal,omitempty"`    // principal that is allowed or denied
	NotPrincipal Principal `json:"NotPrincipal,omitempty"` // exception to a list of principals
	Action       Value     `json:"Action"`                 // allowed or denied action
	NotAction    Value     `json:"NotAction,omitempty"`    // matches everything except
	Resource     Value     `json:"Resource,omitempty"`     // object or objects that the statement covers
	NotResource  Value     `json:"NotResource,omitempty"`  // matches everything except
	Condition    Condition `json:"Condition,omitempty"`    // conditions for when a policy is in effect
}

type Principal map[string]Value
type Condition map[string]ConditionStatement
type ConditionStatement map[string]Value

func PolicyDeepEqual(p1, p2 PolicyDocument) bool {
	if p1.Version != p2.Version {
		log.Debugf("policy version %s is not the same as %s", p1.Version, p2.Version)
		return false
	}

	if len(p1.Statement) != len(p2.Statement) {
		log.Debugf("policy statement length %d is not the same as %d", len(p1.Statement), len(p2.Statement))
		return false
	}

	for _, p1s := range p1.Statement {
		log.Debugf("looking for statement matching %+v", p1s)

		var statementEqual bool
		for _, p2s := range p2.Statement {
			log.Debugf("comparing statement %+v", p2s)

			// compare policies with the same SID
			if p1s.Sid != p2s.Sid {
				log.Debugf("SID %s doesn't match %s", p1s.Sid, p2s.Sid)
				continue
			}

			// if the Effect is different
			if p1s.Effect != p2s.Effect {
				log.Debugf("Effect %s doesn't match %s", p1s.Effect, p2s.Effect)
				return false
			}

			// if the Principals are different
			if !p1s.Principal.Equal(p2s.Principal) {
				log.Debugf("Principal %s doesn't match %s", p1s.Principal, p2s.Principal)
				return false
			}

			if !p1s.Resource.Equal(p2s.Resource) {
				log.Debugf("Resource list %+v doesn't match %+v", p1s.Resource, p2s.Resource)
				return false
			}

			// if the actions are different
			if !p1s.Action.Equal(p2s.Action) {
				log.Debugf("Action list %+v doesn't match %+v", p1s.Action, p2s.Action)
				return false
			}

			// if the conditions are different
			if !p1s.Condition.Equal(p2s.Condition) {
				log.Debugf("Condition %+v doesn't match %+v", p1s.Condition, p2s.Condition)
				return false
			}

			statementEqual = true
		}

		if statementEqual {
			continue
		}

		return false
	}

	return true
}

type Value []string

// UnmarshalJSON unmarshalls IAM values, converting everything to []string to avoid casting
func (value *Value) UnmarshalJSON(b []byte) error {
	var raw interface{}
	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	var p []string
	//  value can be string or []string, convert everything to []string
	switch v := raw.(type) {
	case string:
		p = []string{v}
	case []interface{}:
		var items []string
		for _, item := range v {
			items = append(items, fmt.Sprintf("%v", item))
		}
		p = items
	default:
		return fmt.Errorf("invalid %s value element: allowed is only string or []string", value)
	}

	*value = p
	return nil
}

func (v Value) Equal(v1 Value) bool {
	if len(v) != len(v1) {
		return false
	}

	for _, i1 := range v {
		var found bool
		for _, i2 := range v1 {
			if i1 == i2 {
				found = true
				break
			}
		}

		if found {
			continue
		}

		return false
	}

	return true
}

func (p Principal) Equal(p1 Principal) bool {
	if len(p) != len(p1) {
		return false
	}

	for k, v := range p {
		v1, ok := p1[k]
		if !ok {
			return false
		}

		if !v.Equal(v1) {
			return false
		}
	}

	return true
}

func (c Condition) Equal(c1 Condition) bool {
	if len(c) != len(c1) {
		return false
	}

	for k, cs := range c {
		cs1, ok := c1[k]
		if !ok {
			return false
		}

		if !cs.Equal(cs1) {
			return false
		}
	}

	return true
}

func (c ConditionStatement) Equal(c1 ConditionStatement) bool {
	if len(c) != len(c1) {
		return false
	}

	for k, v := range c {
		v1, ok := c1[k]
		if !ok {
			return false
		}

		if !v.Equal(v1) {
			return false
		}
	}

	return true
}

func (i *IAM) GetPolicyByName(ctx context.Context, name, path string) (*iam.Policy, error) {
	if name == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if path == "" {
		path = "/"
	}

	var policy *iam.Policy
	err := i.Service.ListPoliciesPagesWithContext(ctx, &iam.ListPoliciesInput{PathPrefix: aws.String(path)}, func(page *iam.ListPoliciesOutput, lastPage bool) bool {
		for _, p := range page.Policies {
			if aws.StringValue(p.PolicyName) == name {
				policy = p
				log.Debugf("found policy with name %s and path %s", name, path)
				return false
			}
		}
		return true
	})

	if err != nil {
		return nil, ErrCode("failed finding policy", err)
	}

	if policy == nil {
		return nil, apierror.New(apierror.ErrNotFound, "policy not found", nil)
	}

	if aws.StringValue(policy.Path) != path {
		msg := fmt.Sprintf("policy %s found, but not in path %s (actual path %s)", name, path, aws.StringValue(policy.Path))
		return nil, apierror.New(apierror.ErrNotFound, msg, nil)
	}

	return policy, nil
}

func (i *IAM) GetDefaultPolicyVersion(ctx context.Context, arn, version string) (*iam.PolicyVersion, error) {
	if arn == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("getting policy document for %s", arn)

	out, err := i.Service.GetPolicyVersionWithContext(ctx, &iam.GetPolicyVersionInput{
		PolicyArn: aws.String(arn),
		VersionId: aws.String(version),
	})

	if err != nil {
		return nil, ErrCode("failed to get policy version", err)
	}

	log.Debugf("got output from getting policy version: %+v", out)

	return out.PolicyVersion, nil
}

func (i *IAM) WaitForPolicy(ctx context.Context, policyArn string) error {
	if policyArn == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if err := i.Service.WaitUntilPolicyExistsWithContext(ctx, &iam.GetPolicyInput{
		PolicyArn: aws.String(policyArn),
	}); err != nil {
		return ErrCode("failed waiting for policy to create", err)
	}

	return nil
}

func (i *IAM) CreatePolicy(ctx context.Context, name, path, policyDoc string) (*iam.Policy, error) {
	if name == "" || policyDoc == "" {
		return nil, apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	if path == "" {
		path = "/"
	}

	log.Infof("creating policy %s in path %s", name, path)

	out, err := i.Service.CreatePolicyWithContext(ctx, &iam.CreatePolicyInput{
		Path:           aws.String(path),
		PolicyDocument: aws.String(policyDoc),
		PolicyName:     aws.String(name),
	})

	if err != nil {
		return nil, ErrCode("failed to create policy", err)
	}

	log.Debugf("got output from create policy:  %+v", out)

	return out.Policy, nil
}

func (i *IAM) UpdatePolicy(ctx context.Context, arn, policyDoc string) error {
	if arn == "" || policyDoc == "" {
		return apierror.New(apierror.ErrBadRequest, "invalid input", nil)
	}

	log.Infof("updating policy %s by creating new version and setting as default", arn)

	out, err := i.Service.CreatePolicyVersionWithContext(ctx, &iam.CreatePolicyVersionInput{
		PolicyArn:      aws.String(arn),
		PolicyDocument: aws.String(policyDoc),
		SetAsDefault:   aws.Bool(true),
	})

	if err != nil {
		return ErrCode("failed to update policy", err)
	}

	log.Debugf("got output from create policy version: %+v", out)

	return nil
}
