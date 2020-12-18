package iam

type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

type StatementEntry struct {
	Effect    string
	Action    []string
	Resource  string
	Condition Condition `json:",omitempty"`
}

// Condition maps a condition operator to the condition-key/condition-value statement
// ie. "{ "StringEquals" : { "aws:username" : "johndoe" }}"
// for more information, see https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition.html
type Condition map[string]ConditionStatement

// ConditionStatement maps condition-key to condition-value
// ie. "{ "aws:username" : "johndoe" }"
type ConditionStatement map[string]string
