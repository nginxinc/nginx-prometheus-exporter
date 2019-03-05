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
	type args struct {
		scrapeURI     string
		nginxPlus     bool
		retries       int
		retryInterval time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Nginx Client, valid uri",
			args{
				scrapeURI: "http://demo.nginx.com/stub_status",
				nginxPlus: false,
			},
			false,
		},
		{
			"Nginx Plus Client, valid uri",
			args{
				scrapeURI: "http://demo.nginx.com/api",
				nginxPlus: true,
			},
			false,
		},
		{
			"Nginx Client, invalid uri",
			args{
				scrapeURI: "http://TYPO.nginx.com/stub_status",
				nginxPlus: false,
			},
			true,
		},
		{
			"Nginx Plus Client, invalid uri, retries",
			args{
				scrapeURI:     "http://TYPO.nginx.com/api",
				nginxPlus:     true,
				retries:       2,
				retryInterval: 100 * time.Millisecond,
			},
			true,
		},
	}

	httpClient := &http.Client{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createClientWithRetries(*httpClient, tt.args.scrapeURI, "clientname", tt.args.nginxPlus, tt.args.retries, tt.args.retryInterval)
			log.Printf("Error %v, Expected %v", err, tt.wantErr)
			if (err != nil) != tt.wantErr {
				t.Errorf("createClientWithRetries() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr {
				return
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
