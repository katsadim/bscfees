package handler

import (
	"bsc-fees/pkg"
	"bsc-fees/pkg/config"
	"bsc-fees/pkg/logger"
	"bsc-fees/pkg/net"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

const (
	addrKey    = "address"
	addrPrefix = "0x"
	addrLength = 42
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log := logger.NewLogger()

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Error("failed to parse config file: %w", err)
		errResp := NewInternalError("internal server config error")
		return errorInAPIGatewayProxyResponse(*errResp, net.Headers["dev"])
	}

	corsHeaders := net.SetupCORSHeaders(cfg.General.Env, request.Headers["origin"])

	address, ok := request.QueryStringParameters[addrKey]
	if !ok {
		errResp := NewEmptyQueryParameterError(addrKey)
		return errorInAPIGatewayProxyResponse(*errResp, corsHeaders)
	}

	if errResp := validateAddress(address); errResp != nil {
		return errorInAPIGatewayProxyResponse(*errResp, corsHeaders)
	}

	result, err := pkg.NewOperator(ctx, cfg).Calculate(address)
	if err != nil {
		log.Error("Something went terribly wrong: %w", err)
		errResp := NewInternalError("binance/BSC/ETH error")
		return errorInAPIGatewayProxyResponse(*errResp, corsHeaders)
	}

	resp, err := json.Marshal(result)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    net.Headers[cfg.General.Env],
		Body:       string(resp),
	}, nil
}

func validateAddress(address string) *Error {
	if address == "" {
		return NewEmptyQueryParameterError(addrKey)
	}

	if !strings.HasPrefix(address, addrPrefix) {
		return NewBadRequestError(fmt.Sprintf("[%s] should start with %s", addrKey, addrPrefix))
	}

	if len(address) != addrLength {
		return NewBadRequestError(fmt.Sprintf("[%s] length should be %d", addrKey, addrLength))
	}

	if !isHex(address[2:]) {
		return NewBadRequestError(fmt.Sprintf("[%s] should be hexadecimal", addrKey))
	}

	return nil
}

func errorInAPIGatewayProxyResponse(errorResp Error, corsHeaders map[string]string) (events.APIGatewayProxyResponse, error) {
	jsonErr, err := json.Marshal(errorResp)
	if err != nil {
		return events.APIGatewayProxyResponse{}, errors.New(errorResp.Message)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: errorResp.Code,
		Headers:    corsHeaders,
		Body:       string(jsonErr),
	}, nil
}

func isHex(str string) bool {
	dat := []byte(str)
	for _, c := range dat {
		charIsHex := (c >= '0' && c <= '9') ||
			(c >= 'a' && c <= 'f') ||
			(c >= 'A' && c <= 'F')
		if !charIsHex {
			return false
		}
	}
	return true
}
