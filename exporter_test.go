package main

import (
	"log"
	"net/http"
	"testing"
	"time"

	plusclient "github.com/nginxinc/nginx-plus-go-sdk/client"
	"github.com/nginxinc/nginx-prometheus-exporter/client"
)

func TestCreateClientWithRetries(t *testing.T) {
	type MockedNginxClient struct {
		apiEndpoint string
		httpClient  *http.Client
	}
	type args struct {
		httpClient    *http.Client
		scrapeURI     string
		clientName    string
		nginxPlus     bool
		retries       int
		retryInterval time.Duration
	}

	httpClient := &http.Client{}

	tests := []struct {
		name    string
		args    args
		want    MockedNginxClient
		wantErr bool
	}{
		{
			"Nginx Client, valid uri",
			args{
				httpClient: httpClient,
				clientName: "Nginx Client",
				scrapeURI:  "http://demo.nginx.com/stub_status",
				nginxPlus:  false,
			},
			MockedNginxClient{
				apiEndpoint: "http://demo.nginx.com/stub_status",
				httpClient:  httpClient,
			},
			false,
		},
		{
			"Nginx Plus Client, valid uri",
			args{
				httpClient: httpClient,
				clientName: "Nginx Plus Client",
				scrapeURI:  "http://demo.nginx.com/api",
				nginxPlus:  true,
			},
			MockedNginxClient{
				apiEndpoint: "http://demo.nginx.com/api",
				httpClient:  httpClient,
			},
			false,
		},
		{
			"Nginx Client, invalid uri",
			args{
				httpClient: httpClient,
				clientName: "Nginx Client",
				scrapeURI:  "http://TYPOdemo.nginx.com/stub_status",
				nginxPlus:  false,
			},
			MockedNginxClient{
				apiEndpoint: "http://TYPOdemo.nginx.com/stub_status",
				httpClient:  httpClient,
			},
			true,
		},
		{
			"Nginx Client, invalid uri, retries",
			args{
				httpClient:    httpClient,
				clientName:    "Nginx Client",
				scrapeURI:     "http://TYPOdemo.nginx.com/stub_status",
				nginxPlus:     false,
				retries:       2,
				retryInterval: 1 * time.Second,
			},
			MockedNginxClient{
				apiEndpoint: "http://TYPOdemo.nginx.com/stub_status",
				httpClient:  httpClient,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createClientWithRetries(*tt.args.httpClient, tt.args.scrapeURI, tt.args.clientName, tt.args.nginxPlus, tt.args.retries, tt.args.retryInterval)
			log.Printf("Error %v, Expected %v", err, tt.wantErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("createClientWithRetries() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr {
				return // error returned as wanted
			}

			if tt.args.nginxPlus {
				cl := got.(*plusclient.NginxClient)
				if _, err := cl.GetStats(); err != nil {
					t.Errorf("Failed to create NginxPlusClient: %v", err)
				}
			} else {
				cl := got.(*client.NginxClient)
				if _, err := cl.GetStubStats(); err != nil {
					t.Errorf("Failed to create NginxClient: %v", err)
				}
			}
		})
	}
}
