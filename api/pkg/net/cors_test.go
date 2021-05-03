package net

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var bscfeesHeaders = map[string]string{
	"Access-Control-Allow-Origin":  "https://bscfees.com",
	"Access-Control-Allow-Methods": "GET OPTIONS",
	"Access-Control-Allow-Headers": "Accept-Content",
}

var wwwBscfeesHeaders = map[string]string{
	"Access-Control-Allow-Origin":  "https://www.bscfees.com",
	"Access-Control-Allow-Methods": "GET OPTIONS",
	"Access-Control-Allow-Headers": "Accept-Content",
}

func TestSetupCORSHeaders(t *testing.T) {

	tests := []struct {
		name   string
		env    string
		origin string
		want   map[string]string
	}{
		{"dev env should allow all origins", "dev", "nobody cares", Headers["dev"]},
		{"prod should not allow unknown origin", "prod", "https://not-known-origin.com", Headers["prod"]},
		{"prod should replay https://bscfees.com", "prod", "https://bscfees.com", bscfeesHeaders},
		{"prod should replay https://www.bscfees.com", "prod", "https://www.bscfees.com", wwwBscfeesHeaders},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := SetupCORSHeaders(tt.env, tt.origin)
			assert.Equal(t, headers, tt.want)
		})
	}
}
