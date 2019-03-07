package main

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestCreateClientWithRetries(t *testing.T) {
	type args struct {
		getData       func() (interface{}, error)
		retries       int
		retryInterval time.Duration
	}

	httpClient := &http.Client{}

	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			"Valid func",
			args{
				getData: func() (interface{}, error) { return "soup", nil },
			},
			"soup",
			false,
		},
		{
			"Inavlid func",
			args{
				getData: func() (interface{}, error) { return httpClient.Get("http://FAKE.notarealwebsite.com") },
			},
			nil,
			true,
		},
		{
			"Invalid func with retries",
			args{
				getData:       func() (interface{}, error) { return httpClient.Get("http://FAKE.notarealwebsite.com") },
				retries:       3,
				retryInterval: time.Millisecond * 100,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createClientWithRetries(tt.args.getData, tt.args.retries, tt.args.retryInterval)
			if (err != nil) != tt.wantErr {
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
