package main

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestCreateClientWithRetries(t *testing.T) {
	t.Parallel()

	type args struct {
		client        interface{}
		err           error
		retries       uint
		retryInterval time.Duration
	}

	tests := []struct {
		name            string
		args            args
		expectedRetries int
		want            interface{}
		wantErr         bool
	}{
		{
			name: "getClient returns a valid client",
			args: args{
				client: "client",
				err:    nil,
			},
			expectedRetries: 0,
			want:            "client",
			wantErr:         false,
		},
		{
			name: "getClient returns an error after no retries",
			args: args{
				client: nil,
				err:    errors.New("error"),
			},
			expectedRetries: 0,
			want:            nil,
			wantErr:         true,
		},
		{
			name: "getClient returns an error after retries",
			args: args{
				client:        nil,
				err:           errors.New("error"),
				retries:       3,
				retryInterval: time.Millisecond * 1,
			},
			expectedRetries: 3,
			want:            nil,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invocations := 0
			getClient := func() (interface{}, error) {
				invocations++
				return tt.args.client, tt.args.err
			}

			got, err := createClientWithRetries(getClient, tt.args.retries, tt.args.retryInterval)

			actualRetries := invocations - 1

			if actualRetries != tt.expectedRetries {
				t.Errorf("createClientWithRetries() got %v retries, expected %v", actualRetries, tt.expectedRetries)
				return
			} else if (err != nil) != tt.wantErr {
				t.Errorf("createClientWithRetries() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil && tt.wantErr {
				return
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createClientWithRetries() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsePositiveDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		testInput string
		want      positiveDuration
		wantErr   bool
	}{
		{
			"ParsePositiveDuration returns a positiveDuration",
			"15ms",
			positiveDuration{15 * time.Millisecond},
			false,
		},
		{
			"ParsePositiveDuration returns error for trying to parse negative value",
			"-15ms",
			positiveDuration{},
			true,
		},
		{
			"ParsePositiveDuration returns error for trying to parse empty string",
			"",
			positiveDuration{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePositiveDuration(tt.testInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePositiveDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePositiveDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseUnixSocketAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		testInput       string
		wantSocketPath  string
		wantRequestPath string
		wantErr         bool
	}{
		{
			"Normal unix socket address",
			"unix:/path/to/socket",
			"/path/to/socket",
			"",
			false,
		},
		{
			"Normal unix socket address with location",
			"unix:/path/to/socket:/with/location",
			"/path/to/socket",
			"/with/location",
			false,
		},
		{
			"Unix socket address with trailing ",
			"unix:/trailing/path:",
			"/trailing/path",
			"",
			false,
		},
		{
			"Unix socket address with too many colons",
			"unix:/too:/many:colons:",
			"",
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			socketPath, requestPath, err := parseUnixSocketAddress(tt.testInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseUnixSocketAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(socketPath, tt.wantSocketPath) {
				t.Errorf("socket path: parseUnixSocketAddress() = %v, want %v", socketPath, tt.wantSocketPath)
			}
			if !reflect.DeepEqual(requestPath, tt.wantRequestPath) {
				t.Errorf("request path: parseUnixSocketAddress() = %v, want %v", requestPath, tt.wantRequestPath)
			}
		})
	}
}

func TestParseConstLabels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		labels  string
		want    constLabel
		wantErr bool
	}{
		{
			name:    "Const labels with no labels",
			labels:  "",
			want:    constLabel{},
			wantErr: false,
		},
		{
			name:    "Const labels with one label with valid format",
			labels:  "label=valid",
			want:    constLabel{labels: map[string]string{"label": "valid"}},
			wantErr: false,
		},
		{
			name:    "Const labels with one label with invalid format",
			labels:  "label: invalid",
			want:    constLabel{},
			wantErr: true,
		},
		{
			name:    "Const labels with invalid format for multiple labels",
			labels:  "label=valid,,label2=wrongformat",
			want:    constLabel{},
			wantErr: true,
		},
		{
			name:    "Const labels with multiple labels, one label with invalid format",
			labels:  "label=valid,label2:wrongformat",
			want:    constLabel{},
			wantErr: true,
		},
		{
			name:    "Const labels with label name containing invalid char",
			labels:  "l bel=invalid",
			want:    constLabel{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConstLabels(tt.labels)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConstLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConstLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}
