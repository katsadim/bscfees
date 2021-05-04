package binance

import (
	"bsc-fees/pkg/config"
	"bsc-fees/pkg/net"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var now = time.Now()

func Test_binanceService_GetCurrencyRate_Responses(t *testing.T) {
	tests := []struct {
		name    string
		times   []time.Time
		want    map[time.Time]float64
		wantErr bool
	}{
		{"no_candles", []time.Time{now}, map[time.Time]float64{}, true},
		{"one_candle", []time.Time{now}, map[time.Time]float64{now: 0.407879005}, false},
		{"two_candles", []time.Time{time.Unix(1617225720, 0)}, map[time.Time]float64{time.Unix(1617225720, 0): 300.67275}, false}, //inside first candle
		{"two_candles", []time.Time{time.Unix(1917225720, 0)}, map[time.Time]float64{time.Unix(1917225720, 0): 300.67275}, false}, // outside any candle
		{"two_candles", []time.Time{time.Unix(1617225790, 0)}, map[time.Time]float64{time.Unix(1617225790, 0): 300.7937}, false},  // inside second candle
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := readFile(tt.name + ".json")
			require.Nil(t, err)

			mockClient := new(net.MockClient)
			s := NewBinanceService(
				mockClient,
				config.Binance{},
			)

			mockClient.On(
				"SendRequest",
				mock.Anything,
				mock.Anything,
			).Return(
				&http.Response{Body: file, StatusCode: http.StatusOK},
				nil)

			got, err := s.GetCurrencyRates(tt.times, "BNBBUSD")

			if !assert.Equal(t, tt.wantErr, err != nil, "Want error: "+strconv.FormatBool(tt.wantErr)) && err != nil {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func readFile(filename string) (io.ReadCloser, error) {
	content, err := ioutil.ReadFile(filepath.Join("testdata", filename))

	if err != nil {
		return nil, fmt.Errorf("failed reading file: %s", err)
	}

	return ioutil.NopCloser(bytes.NewReader(content)), nil
}
