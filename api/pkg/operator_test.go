package pkg

import (
	"bsc-fees/pkg/binance"
	"bsc-fees/pkg/config"
	"bsc-fees/pkg/eth"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func Test_operator_Calculate(t *testing.T) {
	type fields struct {
		bscServiceError     bool
		ethServiceError     bool
		binanceServiceError bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"everything went wrong", fields{
			bscServiceError:     true,
			ethServiceError:     true,
			binanceServiceError: true,
		}, true},
		{"Symbol rate failed", fields{
			bscServiceError:     false,
			ethServiceError:     false,
			binanceServiceError: true,
		}, true},
		{"Both eth and bsc services failed", fields{
			bscServiceError:     true,
			ethServiceError:     true,
			binanceServiceError: false,
		}, true},
		{"Only bsc service failed", fields{
			bscServiceError:     true,
			ethServiceError:     false,
			binanceServiceError: false,
		}, false},
		{"Only eth service failed", fields{
			bscServiceError:     false,
			ethServiceError:     true,
			binanceServiceError: false,
		}, false},
		{"Everything went according to the plan!", fields{
			bscServiceError:     false,
			ethServiceError:     false,
			binanceServiceError: false,
		}, false},
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Could not read config file")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockCtrl := gomock.NewController(t)
			mockBscTxGetter := eth.NewMockTxGetter(mockCtrl)
			setupTxGetterMock(mockBscTxGetter, tt.fields.bscServiceError)

			mockEthTxGetter := eth.NewMockTxGetter(mockCtrl)
			setupTxGetterMock(mockEthTxGetter, tt.fields.ethServiceError)

			mockRateGetter := binance.NewMockRateGetter(mockCtrl)
			setupRateGetterMock(mockRateGetter, tt.fields.binanceServiceError)

			o := NewManualOperator(
				context.Background(),
				cfg,
				mockBscTxGetter,
				mockEthTxGetter,
				mockRateGetter)
			_, err := o.Calculate("0x000000000000000000000000000000000000000000")
			if (err != nil) != tt.wantErr {
				assert.Fail(t, "Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			mockCtrl.Finish()
		})
	}
}

func setupTxGetterMock(m *eth.MockTxGetter, returnError bool) {
	if returnError {
		m.EXPECT().
			GetAccountTransactions(gomock.Any()).
			Return(nil, fmt.Errorf("there was an error"))
	} else {
		txs := []eth.Transaction{
			{Time: time.Now(), GasPrice: *big.NewInt(12), GasUsed: 14},
		}
		m.EXPECT().
			GetAccountTransactions(gomock.Any()).
			Return(txs, nil)
	}
}

func setupRateGetterMock(m *binance.MockRateGetter, returnError bool) {
	if returnError {
		m.EXPECT().
			GetLatestCurrencyRate(gomock.Any()).
			Return(0.0, fmt.Errorf("there was an error")).
			Times(2)
	} else {
		m.EXPECT().
			GetLatestCurrencyRate(gomock.Any()).
			Return(15.0, nil).
			Times(2)
	}
}
