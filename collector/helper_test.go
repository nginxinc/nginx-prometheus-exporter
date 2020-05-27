package collector

import (
	"reflect"
	"testing"
)

func TestMergeLabelValues(t *testing.T) {
	type args struct {
		a map[string]string
		b map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			"base case",
			args{
				map[string]string{"poda": "labela"},
				map[string]string{"podb": "labelb"},
			},
			map[string]string{"poda": "labela", "podb": "labelb"},
		},
		{
			"overwrite",
			args{
				map[string]string{"poda": "labela"},
				map[string]string{"poda": "labelb"},
			},
			map[string]string{"poda": "labelb"},
		},
		{
			"merge reversed order",
			args{
				map[string]string{"podb": "labelb"},
				map[string]string{"poda": "labela"},
			},
			map[string]string{"poda": "labela", "podb": "labelb"},
		},
		{
			"nil merge",
			args{
				map[string]string{"podb": "labelb"},
				nil,
			},
			map[string]string{"podb": "labelb"},
		},
		{
			"nil merge reversed order",
			args{
				nil,
				map[string]string{"podb": "labelb"},
			},
			map[string]string{"podb": "labelb"},
		},
		{
			"two nils",
			args{
				nil,
				nil,
			},
			map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MergeLabelValues(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeLabelValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveBlankLabelNames(t *testing.T) {
	type args struct {
		labels []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Remove blank label",
			args{
				[]string{""},
			},
			[]string{},
		},
		{
			"Remove blank label from slice",
			args{
				[]string{"label", "", "label3"},
			},
			[]string{"label", "label3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveBlankLabelNames(tt.args.labels...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveBlankLabelNames() = %v, want %v", got, tt.want)
			}
		})
	}
}
