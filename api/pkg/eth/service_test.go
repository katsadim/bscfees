package eth

import (
	"bsc-fees/pkg/config"
	"bsc-fees/pkg/net"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var firstTX Transaction
var secondTX Transaction

func TestMain(m *testing.M) {
	var firstTime = Time{}
	var secondTime = Time{}
	if err := firstTime.UnmarshalText([]byte("1614835689")); err != nil {
		log.Fatal("Failed unmarshal")
	}
	if err := secondTime.UnmarshalText([]byte("1615833729")); err != nil {
		log.Fatal("Failed unmarshal")
	}
	firstTX = Transaction{Time: firstTime.Time(), GasPrice: *big.NewInt(11000000000), GasUsed: 21000}
	secondTX = Transaction{Time: secondTime.Time(), GasPrice: *big.NewInt(10000000000), GasUsed: 23000}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func Test_service_GetAccountTransactions_Responses(t *testing.T) {

	tests := []struct {
		name    string
		want    []Transaction
		wantErr bool
	}{
		{"one_transaction", []Transaction{firstTX}, false},
		{"two_transactions", []Transaction{firstTX, secondTX}, false},
		{"no_transactions", []Transaction{}, false},
		{"error_status", []Transaction{}, true},
		{"invalid_json", []Transaction{}, true},
		{"invalid_results", []Transaction{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := readFile(tt.name + ".json")
			require.Nil(t, err)

			mockClient := new(net.MockClient)
			s := NewProviderService(
				mockClient,
				config.Provider{},
			)

			mockClient.On(
				"SendRequest",
				mock.Anything,
				mock.Anything,
			).Return(
				&http.Response{Body: file, StatusCode: http.StatusOK},
				nil)

			got, err := s.GetAccountTransactions("dummy addr")

			if !assert.Equal(t, tt.wantErr, err != nil, "Want error: "+strconv.FormatBool(tt.wantErr)) && err != nil {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_service_GetAccountTransactions_InvalidRequest(t *testing.T) {
	mockClient := new(net.MockClient)
	s := NewProviderService(
		mockClient,
		config.Provider{},
	)

	mockClient.On(
		"SendRequest",
		mock.Anything,
		mock.Anything,
	).Return(
		nil,
		fmt.Errorf("something went terribly wrong"))

	_, err := s.GetAccountTransactions("dummy addr")

	assert.NotNil(t, err)
}

func readFile(filename string) (io.ReadCloser, error) {
	content, err := ioutil.ReadFile(filepath.Join("testdata", filename))

	if err != nil {
		return nil, fmt.Errorf("failed reading file: %s", err)
	}

	return ioutil.NopCloser(bytes.NewReader(content)), nil
}
