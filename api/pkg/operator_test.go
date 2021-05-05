package pkg

import (
	"bsc-fees/pkg/binance"
	"bsc-fees/pkg/config"
	"bsc-fees/pkg/eth"
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func Test_operator_Calculate(t *testing.T) {

	type fields struct {
		bscService     eth.TxGetter
		ethService     eth.TxGetter
		binanceService binance.RateGetter
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"everything went wrong", fields{
			bscService:     new(erroneousTxGetter),
			ethService:     new(erroneousTxGetter),
			binanceService: new(erroneousRateGetter),
		}, true},
		{"Symbol rate failed", fields{
			bscService:     new(successfulTxGetter),
			ethService:     new(successfulTxGetter),
			binanceService: new(erroneousRateGetter),
		}, true},
		{"Both eth and bsc services failed", fields{
			bscService:     new(erroneousTxGetter),
			ethService:     new(erroneousTxGetter),
			binanceService: new(successfulRateGetter),
		}, true},
		{"Only bsc service failed", fields{
			bscService:     new(erroneousTxGetter),
			ethService:     new(successfulTxGetter),
			binanceService: new(successfulRateGetter),
		}, false},
		{"Only eth service failed", fields{
			bscService:     new(successfulTxGetter),
			ethService:     new(erroneousTxGetter),
			binanceService: new(successfulRateGetter),
		}, false},
		{"Everything went according to the plan!", fields{
			bscService:     new(successfulTxGetter),
			ethService:     new(successfulTxGetter),
			binanceService: new(successfulRateGetter),
		}, false},
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Could not read config file")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := NewManualOperator(
				context.Background(),
				cfg,
				tt.fields.bscService,
				tt.fields.ethService,
				tt.fields.binanceService)
			_, err := o.Calculate("0x000000000000000000000000000000000000000000")
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

type erroneousTxGetter struct {
}

func (e *erroneousTxGetter) GetAccountTransactions(address string) ([]eth.Transaction, error) {
	return nil, fmt.Errorf("there was an error")
}

type successfulTxGetter struct {
}

func (s *successfulTxGetter) GetAccountTransactions(address string) ([]eth.Transaction, error) {
	txs := []eth.Transaction{
		{Time: time.Now(), GasPrice: *big.NewInt(12), GasUsed: 14},
	}
	return txs, nil
}

type erroneousRateGetter struct {
}

func (e *erroneousRateGetter) GetCurrencyRates(times []time.Time, symbol string) (map[time.Time]float64, error) {
	return nil, fmt.Errorf("there was an error")
}

func (e *erroneousRateGetter) GetLatestCurrencyRate(symbol string) (float64, error) {
	return 0, fmt.Errorf("there was an error")
}

type successfulRateGetter struct {
}

func (s *successfulRateGetter) GetCurrencyRates(times []time.Time, symbol string) (map[time.Time]float64, error) {
	return map[time.Time]float64{}, nil
}

func (s *successfulRateGetter) GetLatestCurrencyRate(symbol string) (float64, error) {
	return 15, nil
}
