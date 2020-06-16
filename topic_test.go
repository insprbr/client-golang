package client

import (
	"reflect"
	"testing"
)

func Test_fromTopic(t *testing.T) {
	type args struct {
		in0 string
		env *Env
	}
	tests := []struct {
		name string
		args args
		want channel
	}{
		{
			"No prefix topic",
			args{
				"chimera_namespace_name",
				testEnvVars("name", "namespace", "", "", "", "", ""),
			},
			channel{
				name:      "name",
				namespace: "namespace",
			},
		},
		{
			"Prefix topic",
			args{
				"chimera_prefix_namespace_name",
				testEnvVars("name", "namespace", "", "", "", "", "prefix"),
			},
			channel{
				name:      "name",
				namespace: "namespace",
				prefix:    "prefix",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envs = tt.args.env
			if got := fromTopic(tt.args.in0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fromTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toTopic(t *testing.T) {
	type args struct {
		ch  string
		env *Env
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"With prefix",
			args{
				"name",
				testEnvVars("name", "namespace", "", "", "", "", "prefix"),
			},
			"chimera-prefix-namespace-name",
		},
		{
			"Without prefix",
			args{
				"name",
				testEnvVars("name", "namespace", "", "", "", "", ""),
			},
			"chimera-namespace-name",
		},
	}
	for _, tt := range tests {
		envs = tt.args.env
		t.Run(tt.name, func(t *testing.T) {
			if got := toTopic(tt.args.ch); got != tt.want {
				t.Errorf("toTopic() = %v, want %v", got, tt.want)
			}
		})
	}
}
