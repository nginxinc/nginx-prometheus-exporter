package main

import (
	"reflect"
	"testing"
	"time"
)

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
