package handler

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleRequest_address(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{"should exist", args{}, 400},
		{"should not be empty", args{address: ""}, 400},
		{"should start with 0x", args{address: "0m036ACb8567E497994C9115002Ecf78794EFaa48B"}, 400},
		{"length should be 42", args{address: "0x036ACb8567E497994C91150"}, 400},
		{"should be hexadecimal", args{address: "0x036ACb8567E497994C9115002Ecf7879ZZZZZZZZ"}, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleRequest(createApiGwRequest(tt.args.address))

			assert.Nil(t, err)

			assert.Equal(t, tt.wantStatusCode, got.StatusCode)
		})
	}
}

func createApiGwRequest(address string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{"address": address},
	}
}
