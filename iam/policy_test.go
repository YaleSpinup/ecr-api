package iam

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
)

func TestIAM_GetPolicyByName(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *iam.Policy
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.GetPolicyByName(tt.args.ctx, tt.args.name, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.GetPolicyByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.GetPolicyByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_GetDefaultPolicyVersion(t *testing.T) {
	type args struct {
		ctx     context.Context
		arn     string
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    *iam.PolicyVersion
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.GetDefaultPolicyVersion(tt.args.ctx, tt.args.arn, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.GetDefaultPolicyVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.GetDefaultPolicyVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_WaitForPolicy(t *testing.T) {

	type args struct {
		ctx       context.Context
		policyArn string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			if err := i.WaitForPolicy(tt.args.ctx, tt.args.policyArn); (err != nil) != tt.wantErr {
				t.Errorf("IAM.WaitForPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIAM_CreatePolicy(t *testing.T) {
	type args struct {
		ctx       context.Context
		name      string
		path      string
		policyDoc string
	}
	tests := []struct {
		name    string
		args    args
		want    *iam.Policy
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			got, err := i.CreatePolicy(tt.args.ctx, tt.args.name, tt.args.path, tt.args.policyDoc)
			if (err != nil) != tt.wantErr {
				t.Errorf("IAM.CreatePolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IAM.CreatePolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIAM_UpdatePolicy(t *testing.T) {
	type args struct {
		ctx       context.Context
		arn       string
		policyDoc string
	}
	tests := []struct {
		name    string
		args    args
		err     error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &IAM{Service: newMockIAMClient(t, tt.err)}
			if err := i.UpdatePolicy(tt.args.ctx, tt.args.arn, tt.args.policyDoc); (err != nil) != tt.wantErr {
				t.Errorf("IAM.UpdatePolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConditionStatement_Equal(t *testing.T) {
	type args struct {
		c1 ConditionStatement
	}
	tests := []struct {
		name string
		c    ConditionStatement
		args args
		want bool
	}{
		{
			name: "empty conditions",
			c:    ConditionStatement{},
			args: args{},
			want: true,
		},
		{
			name: "empty c1, set c2",
			c:    ConditionStatement{},
			args: args{
				c1: ConditionStatement{
					"foo": Value{"bar"},
				},
			},
			want: false,
		},
		{
			name: "empty c1, set c2",
			c: ConditionStatement{
				"foo": Value{"bar"},
			},
			args: args{
				c1: ConditionStatement{},
			},
			want: false,
		},
		{
			name: "same key, different values",
			c: ConditionStatement{
				"key1": Value{"somevalue1"},
			},
			args: args{
				c1: ConditionStatement{
					"key1": Value{"somevalue2"},
				},
			},
			want: false,
		},
		{
			name: "same key, same values",
			c: ConditionStatement{
				"key1": Value{"somevalue1"},
			},
			args: args{
				c1: ConditionStatement{
					"key1": Value{"somevalue1"},
				},
			},
			want: true,
		},
		{
			name: "same keys and values",
			c: ConditionStatement{
				"key1": []string{"somevalue1"},
				"key2": []string{"somevalue2"},
				"key3": []string{"somevalue3"},
				"key4": []string{"somevalue4"},
				"key5": []string{"somevalue5"},
			},
			args: args{
				c1: ConditionStatement{
					"key1": []string{"somevalue1"},
					"key2": []string{"somevalue2"},
					"key3": []string{"somevalue3"},
					"key4": []string{"somevalue4"},
					"key5": []string{"somevalue5"},
				},
			},
			want: true,
		},
		{
			name: "same keys and values, different order",
			c: ConditionStatement{
				"key5": []string{"somevalue5"},
				"key4": []string{"somevalue4"},
				"key3": []string{"somevalue3"},
				"key2": []string{"somevalue2"},
				"key1": []string{"somevalue1"},
			},
			args: args{
				c1: ConditionStatement{
					"key1": []string{"somevalue1"},
					"key2": []string{"somevalue2"},
					"key3": []string{"somevalue3"},
					"key4": []string{"somevalue4"},
					"key5": []string{"somevalue5"},
				},
			},
			want: true,
		},
		{
			name: "same length, different keys",
			c: ConditionStatement{
				"key1": Value{"stringvalue1"},
				"key2": Value{"stringvalue1"},
			},
			args: args{
				c1: ConditionStatement{
					"key3": Value{"stringvalue1"},
					"key4": Value{"stringvalue1"},
				},
			},
			want: false,
		},
		{
			name: "same keys, different order value",
			c: ConditionStatement{
				"key1": []string{"somevalue1", "somevalue2", "somevalue3", "somevalue4", "somevalue5"},
			},
			args: args{
				c1: ConditionStatement{
					"key1": []string{"somevalue5", "somevalue4", "somevalue3", "somevalue2", "somevalue1"},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Equal(tt.args.c1); got != tt.want {
				t.Errorf("ConditionStatement.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCondition_Equal(t *testing.T) {
	type args struct {
		c1 Condition
	}
	tests := []struct {
		name string
		c    Condition
		args args
		want bool
	}{
		{
			name: "empty conditions",
			c:    Condition{},
			args: args{
				c1: Condition{},
			},
			want: true,
		},
		{
			name: "empty conditions, c1 not empty",
			c:    Condition{},
			args: args{
				c1: Condition{
					"key1": ConditionStatement{},
				},
			},
			want: false,
		},
		{
			name: "not empty conditions, c1 empty",
			c: Condition{
				"key1": ConditionStatement{},
			},
			args: args{
				c1: Condition{},
			},
			want: false,
		},
		{
			name: "same length, equal keys",
			c: Condition{
				"key1": ConditionStatement{},
				"key2": ConditionStatement{},
				"key3": ConditionStatement{},
			},
			args: args{
				c1: Condition{
					"key1": ConditionStatement{},
					"key2": ConditionStatement{},
					"key3": ConditionStatement{},
				},
			},
			want: true,
		},
		{
			name: "same length, equal keys",
			c: Condition{
				"key1": ConditionStatement{},
				"key2": ConditionStatement{},
				"key3": ConditionStatement{},
			},
			args: args{
				c1: Condition{
					"key4": ConditionStatement{},
					"key5": ConditionStatement{},
					"key6": ConditionStatement{},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Equal(tt.args.c1); got != tt.want {
				t.Errorf("Condition.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrincipal_Equal(t *testing.T) {
	type args struct {
		p1 Principal
	}
	tests := []struct {
		name string
		p    Principal
		args args
		want bool
	}{
		{
			name: "empty principals",
			p:    Principal{},
			args: args{
				p1: Principal{},
			},
			want: true,
		},
		{
			name: "empty principal, not empty p1",
			p:    Principal{},
			args: args{
				p1: Principal{
					"key1": Value{"value1"},
				},
			},
			want: false,
		},
		{
			name: "not empty principal, empty p1",
			p: Principal{
				"key1": Value{"value1"},
			},
			args: args{
				p1: Principal{},
			},
			want: false,
		},
		{
			name: "equal principals, 1 key and value",
			p: Principal{
				"key1": Value{"value1"},
			},
			args: args{
				p1: Principal{
					"key1": Value{"value1"},
				},
			},
			want: true,
		},
		{
			name: "different principal values",
			p: Principal{
				"key1": Value{"value1"},
			},
			args: args{
				p1: Principal{
					"key1": Value{"value2"},
				},
			},
			want: false,
		},
		{
			name: "equal principals many key and values",
			p: Principal{
				"key1": Value{"value1"},
				"key2": Value{"value21", "value22", "value23"},
				"key3": Value{"value31", "value32"},
			},
			args: args{
				p1: Principal{
					"key1": Value{"value1"},
					"key2": Value{"value21", "value22", "value23"},
					"key3": Value{"value31", "value32"},
				},
			},
			want: true,
		},
		{
			name: "same keys different order",
			p: Principal{
				"key1": Value{"value1"},
				"key2": Value{"value21", "value22", "value23"},
				"key3": Value{"value31", "value32"},
			},
			args: args{
				p1: Principal{
					"key3": Value{"value31", "value32"},
					"key2": Value{"value21", "value22", "value23"},
					"key1": Value{"value1"},
				},
			},
			want: true,
		},
		{
			name: "same keys different value order",
			p: Principal{
				"key1": Value{"value1", "value2", "value3"},
			},
			args: args{
				p1: Principal{
					"key1": Value{"value3", "value2", "value1"},
				},
			},
			want: true,
		},
		{
			name: "same number, different keys",
			p: Principal{
				"key1": Value{"value1"},
				"key2": Value{"value1"},
			},
			args: args{
				p1: Principal{
					"key3": Value{"value1"},
					"key4": Value{"value1"},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Equal(tt.args.p1); got != tt.want {
				t.Errorf("Principal.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
