package pkg

import (
	"bsc-fees/pkg/binance"
	"bsc-fees/pkg/bsc"
	"bsc-fees/pkg/config"
	"bsc-fees/pkg/net"
	"context"
	"fmt"
)

type Operator interface {
	Calculate(account string) (Response, error)
}

type operator struct {
	ctx            context.Context
	cfg            *config.Config // does it need to be a reference?
	bscService     bsc.TxGetter
	binanceService binance.RateGetter
	calculator     Calculator
}

func NewOperator(ctx context.Context, cfg config.Config) Operator {
	bscClient := net.NewClient(cfg)
	binanceClient := net.NewClient(cfg)
	bscService := bsc.NewBscService(bscClient, cfg.Bsc)
	binanceService := binance.NewBinanceService(binanceClient, cfg.Binance)
	calculator := NewCalculator(cfg.Options)
	return &operator{
		ctx:            ctx,
		cfg:            &cfg,
		bscService:     bscService,
		binanceService: binanceService,
		calculator:     calculator,
	}
}

func (o *operator) Calculate(account string) (Response, error) {
	feesCh := make(chan float64)
	currRateCh := make(chan float64)
	errorCh := make(chan error, 2)

	go func() {
		fees, err := calculateFees(o, account)
		if err != nil {
			errorCh <- err
			return
		}
		feesCh <- fees
	}()

	go func() {
		currentRate, err := GetLatestCurrencyRate(o)
		if err != nil {
			errorCh <- err
			return
		}
		currRateCh <- currentRate
	}()

	var fees float64
	var currentRate float64
	var err error
	for i := 0; i < 2; i++ {
		select {
		case fees = <-feesCh:
		case currentRate = <-currRateCh:
		case err = <-errorCh:
		}
	}

	if err != nil {
		return Response{}, err
	}

	return Response{
		BnbusdPrice: currentRate,
		Fees:        fees,
	}, nil
}

func calculateFees(o *operator, account string) (float64, error) {
	transactions, err := o.bscService.GetAccountTransactions(account)
	if err != nil {
		return 0, fmt.Errorf("something went terribly wrong with BSC networking: %w", err)
	}
	fees := o.calculator.CalculateFees(transactions)

	f, _ := fees.Float64()
	return f, nil
}

func GetLatestCurrencyRate(o *operator) (float64, error) {
	currentRate, err := o.binanceService.GetLatestCurrencyRate()
	if err != nil {
		return 0, fmt.Errorf("something went terribly wrong with Binance networking: %w", err)
	}
	return currentRate, nil
}
