package iam

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

var testTime = time.Now()
var testPastTime = time.Unix(rand.Int63n(time.Now().Unix()), 0)

// mockIAMClient is a fake IAM client
type mockIAMClient struct {
	iamiface.IAMAPI
	t   *testing.T
	err error
}

func newMockIAMClient(t *testing.T, err error) iamiface.IAMAPI {
	return &mockIAMClient{
		t:   t,
		err: err,
	}
}

func TestNewSession(t *testing.T) {
	client := New()
	to := reflect.TypeOf(client).String()
	if to != "iam.IAM" {
		t.Errorf("expected type to be iam.IAM, got %s", to)
	}
}

var testPolicyDocMap map[string]PolicyDocument = map[string]PolicyDocument{
	"v1": {
		Version: "v1",
	},
	"v2": {
		Version: "v2",
	},
	"oneStatement": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid: "oneSID",
			},
		},
	},
	"oneStatementa": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid: "oneSIDa",
			},
		},
	},
	"twoStatement": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid: "oneSID",
			},
			{
				Sid: "twoSID",
			},
		},
	},
	"twoStatementa": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid: "oneSIDa",
			},
			{
				Sid: "twoSID",
			},
		},
	},
	"oneStatementAllow": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:    "oneSID",
				Effect: "Allow",
			},
		},
	},
	"oneStatementDeny": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:    "oneSID",
				Effect: "Deny",
			},
		},
	},
	"oneStatementOneResource": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1"},
			},
		},
	},
	"oneStatementOneResourcea": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1a"},
			},
		},
	},
	"oneStatementOnePrincipal": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:       "oneSID",
				Effect:    "Allow",
				Principal: Principal{"foo": []string{"bar"}},
			},
		},
	},
	"oneStatementOnePrincipala": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:       "oneSID",
				Effect:    "Allow",
				Principal: Principal{"fooa": []string{"bara"}},
			},
		},
	},
	"oneStatementOne": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1", "resource2"},
			},
		},
	},
	"oneStatementTwoResourcea": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1a", "resource2a"},
			},
		},
	},
	"oneStatementOneAction": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1"},
				Action:   []string{"thing1"},
			},
		},
	},
	"oneStatementOneActiona": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1"},
				Action:   []string{"thing1a"},
			},
		},
	},
	"oneStatementTwoAction": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1"},
				Action:   []string{"thing1", "thing2"},
			},
		},
	},
	"oneStatementTwoActiona": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1a"},
				Action:   []string{"thing1a", "thing2a"},
			},
		},
	},
	"oneConditionOneOperatorOneStatementString": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1a"},
				Action:   []string{"thing1a"},
				Condition: Condition{
					"operator1": ConditionStatement{
						"key1": Value{"stringvalue1"},
					},
				},
			},
		},
	},
	"oneConditionOneOperatorOneStatementStringa": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1a"},
				Action:   []string{"thing1a"},
				Condition: Condition{
					"operator1": ConditionStatement{
						"key1": Value{"stringvalue1a"},
					},
				},
			},
		},
	},
	"oneConditionOneOperatorTwoStatementString": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1a"},
				Action:   []string{"thing1a"},
				Condition: Condition{
					"operator1": ConditionStatement{
						"key1": Value{"stringvalue1"},
						"key2": Value{"stringvalue2"},
					},
				},
			},
		},
	},
	"oneConditionTwoOperatorOneStatementString": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1a"},
				Action:   []string{"thing1a"},
				Condition: Condition{
					"operator1": ConditionStatement{
						"key1": Value{"stringvalue1"},
					},
					"operator2": ConditionStatement{
						"key1": Value{"stringvalue1"},
					},
				},
			},
		},
	},
	"oneConditionTwoOperatorAOneStatementString": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1a"},
				Action:   []string{"thing1a"},
				Condition: Condition{
					"operator1a": ConditionStatement{
						"key1": Value{"stringvalue1"},
					},
					"operator2a": ConditionStatement{
						"key1": Value{"stringvalue1"},
					},
				},
			},
		},
	},
	"oneConditionTwoOperatorTwoStatementString": {
		Version: "v1",
		Statement: []StatementEntry{
			{
				Sid:      "oneSID",
				Effect:   "Allow",
				Resource: []string{"resource1a"},
				Action:   []string{"thing1a"},
				Condition: Condition{
					"operator1": ConditionStatement{
						"key1": Value{"stringvalue1"},
						"key2": Value{"stringvalue2"},
					},
					"operator2": ConditionStatement{
						"key1": Value{"stringvalue1"},
						"key2": Value{"stringvalue2"},
					},
				},
			},
		},
	},
}

func TestPolicyDeepEqual(t *testing.T) {
	type args struct {
		p1 PolicyDocument
		p2 PolicyDocument
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "different version",
			args: args{
				p1: testPolicyDocMap["v1"],
				p2: testPolicyDocMap["v2"],
			},
			want: false,
		},
		{
			name: "same version",
			args: args{
				p1: testPolicyDocMap["v1"],
				p2: testPolicyDocMap["v1"],
			},
			want: true,
		},
		{
			name: "one statement, SIDs equal",
			args: args{
				p1: testPolicyDocMap["oneStatement"],
				p2: testPolicyDocMap["oneStatement"],
			},
			want: true,
		},
		{
			name: "one statement, SIDs different",
			args: args{
				p1: testPolicyDocMap["oneStatement"],
				p2: testPolicyDocMap["oneStatementa"],
			},
			want: false,
		},
		{
			name: "two statements, SIDs equal",
			args: args{
				p1: testPolicyDocMap["twoStatement"],
				p2: testPolicyDocMap["twoStatement"],
			},
			want: true,
		},
		{
			name: "two statement, SIDs different",
			args: args{
				p1: testPolicyDocMap["twoStatement"],
				p2: testPolicyDocMap["twoStatementa"],
			},
			want: false,
		},
		{
			name: "one statement, SIDs equal, same effect",
			args: args{
				p1: testPolicyDocMap["oneStatementAllow"],
				p2: testPolicyDocMap["oneStatementAllow"],
			},
			want: true,
		},
		{
			name: "one statement, SIDs equal, different effect",
			args: args{
				p1: testPolicyDocMap["oneStatementAllow"],
				p2: testPolicyDocMap["oneStatementDeny"],
			},
			want: false,
		},
		{
			name: "one statement, SIDs equal, same effect, one principal, same",
			args: args{
				p1: testPolicyDocMap["oneStatementOnePrincipal"],
				p2: testPolicyDocMap["oneStatementOnePrincipal"],
			},
			want: true,
		},
		{
			name: "one statement, SIDs equal, same effect, one principal, different",
			args: args{
				p1: testPolicyDocMap["oneStatementOnePrincipal"],
				p2: testPolicyDocMap["oneStatementOnePrincipala"],
			},
			want: false,
		},
		{
			name: "one statement, SIDs equal, same effect, one resource, same",
			args: args{
				p1: testPolicyDocMap["oneStatementOneResource"],
				p2: testPolicyDocMap["oneStatementOneResource"],
			},
			want: true,
		},
		{
			name: "one statement, SIDs equal, same effect, one resource, different",
			args: args{
				p1: testPolicyDocMap["oneStatementOneResource"],
				p2: testPolicyDocMap["oneStatementOneResourcea"],
			},
			want: false,
		},
		{
			name: "one statement, SIDs equal, same effect, different # of resources",
			args: args{
				p1: testPolicyDocMap["oneStatementOneResource"],
				p2: testPolicyDocMap["oneStatementTwoResource"],
			},
			want: false,
		},
		{
			name: "one statement, SIDs equal, same effect, same # of resources, different resources",
			args: args{
				p1: testPolicyDocMap["oneStatementTwoResource"],
				p2: testPolicyDocMap["oneStatementTwoResourcea"],
			},
			want: false,
		},

		{
			name: "one statement, SIDs equal, same effect, one action, same",
			args: args{
				p1: testPolicyDocMap["oneStatementOneAction"],
				p2: testPolicyDocMap["oneStatementOneAction"],
			},
			want: true,
		},
		{
			name: "one statement, SIDs equal, same effect, one action, different",
			args: args{
				p1: testPolicyDocMap["oneStatementOneAction"],
				p2: testPolicyDocMap["oneStatementOneActiona"],
			},
			want: false,
		},
		{
			name: "one statement, SIDs equal, same effect, different # of actions",
			args: args{
				p1: testPolicyDocMap["oneStatementOneAction"],
				p2: testPolicyDocMap["oneStatementTwoAction"],
			},
			want: false,
		},
		{
			name: "one statement, SIDs equal, same effect, same # of actions, different",
			args: args{
				p1: testPolicyDocMap["oneStatementTwoAction"],
				p2: testPolicyDocMap["oneStatementTwoActiona"],
			},
			want: false,
		},
		{
			name: "one condition, one operator, one statement string, same",
			args: args{
				p1: testPolicyDocMap["oneConditionOneOperatorOneStatementString"],
				p2: testPolicyDocMap["oneConditionOneOperatorOneStatementString"],
			},
			want: true,
		},
		{
			name: "one condition, one operator, one statement string, different value",
			args: args{
				p1: testPolicyDocMap["oneConditionOneOperatorOneStatementString"],
				p2: testPolicyDocMap["oneConditionOneOperatorOneStatementStringa"],
			},
			want: false,
		},
		{
			name: "one condition, two operator, one statement string, same",
			args: args{
				p1: testPolicyDocMap["oneConditionTwoOperatorOneStatementString"],
				p2: testPolicyDocMap["oneConditionTwoOperatorOneStatementString"],
			},
			want: true,
		},
		{
			name: "one condition, two operator, one statement string, different",
			args: args{
				p1: testPolicyDocMap["oneConditionTwoOperatorOneStatementString"],
				p2: testPolicyDocMap["oneConditionTwoOperatorOneStatementStringa"],
			},
			want: false,
		},
		{
			name: "one condition, two operator, one statement string, different",
			args: args{
				p1: testPolicyDocMap["oneConditionTwoOperatorOneStatementString"],
				p2: testPolicyDocMap["oneConditionTwoOperatorOneStatementStringa"],
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PolicyDeepEqual(tt.args.p1, tt.args.p2); got != tt.want {
				t.Errorf("PolicyDeepEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
