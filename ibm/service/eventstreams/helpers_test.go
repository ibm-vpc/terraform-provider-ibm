// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package eventstreams

import (
	"testing"
)

// TestGetAdminURL verifies that getAdminURL resolves the REST endpoint from
// instance Extensions under either the Classic key ("kafka_http_url") or the
// Enterprise Gen2 key ("dataservices.connection.rest_url"), and returns a
// descriptive error when neither is present.
func TestGetAdminURL(t *testing.T) {
	const classicURL = "https://admin.123.messagehub.cloud.ibm.com:443"
	const gen2URL = "https://gen2-admin.messagehub.cloud.ibm.com:443"

	tests := []struct {
		description string
		extensions  map[string]interface{}
		wantURL     string
		wantErr     bool
	}{
		{
			description: "Classic/Enterprise: kafka_http_url present",
			extensions:  map[string]interface{}{"kafka_http_url": classicURL},
			wantURL:     classicURL,
		},
		{
			description: "Enterprise Gen2: dataservices.connection.rest_url present",
			extensions:  map[string]interface{}{"dataservices.connection.rest_url": gen2URL},
			wantURL:     gen2URL,
		},
		{
			description: "Both keys present: kafka_http_url takes precedence",
			extensions: map[string]interface{}{
				"kafka_http_url":                  classicURL,
				"dataservices.connection.rest_url": gen2URL,
			},
			wantURL: classicURL,
		},
		{
			description: "Neither key present: error returned",
			extensions:  map[string]interface{}{"some_other_key": "value"},
			wantErr:     true,
		},
		{
			description: "Empty extensions map: error returned",
			extensions:  map[string]interface{}{},
			wantErr:     true,
		},
		{
			description: "kafka_http_url is empty string: falls back to gen2 key",
			extensions: map[string]interface{}{
				"kafka_http_url":                  "",
				"dataservices.connection.rest_url": gen2URL,
			},
			wantURL: gen2URL,
		},
		{
			description: "kafka_http_url wrong type (non-string): falls back to gen2 key",
			extensions: map[string]interface{}{
				"kafka_http_url":                  42,
				"dataservices.connection.rest_url": gen2URL,
			},
			wantURL: gen2URL,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			got, err := getAdminURL(tc.extensions)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error but got URL %q", got)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if got != tc.wantURL {
				t.Errorf("got URL %q, want %q", got, tc.wantURL)
			}
		})
	}
}
