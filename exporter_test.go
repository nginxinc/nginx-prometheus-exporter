package main

import (
	"errors"
	"flag"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestCreateClientWithRetries(t *testing.T) {
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

			got, err := createClientWithRetries(getClient, tt.args.retries, &tt.args.retryInterval)

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

func TestCreatePositiveDurationFlag(t *testing.T) {
	type args struct {
		dur    time.Duration
		key    string
		helper string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"CreatePositiveDurationFlag creates a positiveDuration flag",
			args{
				5 * time.Millisecond,
				"key",
				"helper",
			},
			false,
		},
		{
			"CreatePositiveDurationFlag returns an error",
			args{
				-5 * time.Millisecond,
				"neg_key",
				"helper",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createPositiveDurationFlag(&tt.args.dur, tt.args.key, tt.args.helper)

			err := flag.CommandLine.Parse(os.Args[1:])
			if err != nil {
				t.Error(err)
			}

			// Test if flag got added succesfully
			testFlag := flag.Lookup(tt.args.key)
			if testFlag == nil {
				t.Errorf("Got: nil. Expected: %v flag to be found", tt.args.key)
			}
		})
	}
}

func TestPositiveDurationSet(t *testing.T) {
	type fields struct {
		Duration time.Duration
	}
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"PositiveDurationSet updates flag.Value",
			fields{
				5 * time.Millisecond,
			},
			args{
				"15ms",
			},
			false,
		},
		{
			"PositiveDurationSet returns error for trying update to negative value",
			fields{
				5 * time.Millisecond,
			},
			args{
				"-15ms",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pd := &positiveDuration{
				Duration: tt.fields.Duration,
			}
			if err := pd.Set(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("positiveDuration.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
