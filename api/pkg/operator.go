package pkg

import (
	"bsc-fees/pkg/binance"
	"bsc-fees/pkg/eth"
	"bsc-fees/pkg/config"
	"bsc-fees/pkg/net"
	"context"
	"fmt"
)

const concurrentRoutinesNum = 4

type Operator interface {
	Calculate(account string) (Response, error)
}

type operator struct {
	ctx            context.Context
	cfg            *config.Config // does it need to be a reference?
	bscService     eth.TxGetter
	ethService     eth.TxGetter
	binanceService binance.RateGetter
	calculator     Calculator
}

func NewOperator(ctx context.Context, cfg config.Config) Operator {
	bscClient := net.NewClient(cfg)
	binanceClient := net.NewClient(cfg)
	bscService := eth.NewProviderService(bscClient, cfg.Bsc)
	ethService := eth.NewProviderService(bscClient, cfg.Eth)
	binanceService := binance.NewBinanceService(binanceClient, cfg.Binance)
	calculator := NewCalculator(cfg.Options)
	return &operator{
		ctx:            ctx,
		cfg:            &cfg,
		bscService:     bscService,
		ethService:     ethService,
		binanceService: binanceService,
		calculator:     calculator,
	}
}

func (o *operator) Calculate(account string) (Response, error) {
	bscFeesCh := make(chan float64)
	ethFeesCh := make(chan float64)
	bnbusdCurrRateCh := make(chan float64)
	ethusdCurrRateCh := make(chan float64)
	errorCh := make(chan error, concurrentRoutinesNum)

	go func() {
		fees, err := calculateBscFees(o, account)
		if err != nil {
			errorCh <- err
			return
		}
		bscFeesCh <- fees
	}()

	go func() {
		fees, err := calculateEthFees(o, account)
		if err != nil {
			errorCh <- err
			return
		}
		ethFeesCh <- fees
	}()

	go func() {
		currentRate, err := GetLatestCurrencyRate(o, o.cfg.Binance.BnbusdCurrencySymbol)
		if err != nil {
			errorCh <- err
			return
		}
		bnbusdCurrRateCh <- currentRate
	}()

	go func() {
		currentRate, err := GetLatestCurrencyRate(o, o.cfg.Binance.EthusdCurrencySymbol)
		if err != nil {
			errorCh <- err
			return
		}
		ethusdCurrRateCh <- currentRate
	}()

	var bscFees float64
	var ethFees float64
	var bnbusdCurrentRate float64
	var ethusdCurrentRate float64
	var err error
	for i := 0; i < concurrentRoutinesNum; i++ {
		select {
		case bscFees = <-bscFeesCh:
		case ethFees = <-ethFeesCh:
		case bnbusdCurrentRate = <-bnbusdCurrRateCh:
		case ethusdCurrentRate = <-ethusdCurrRateCh:
		case err = <-errorCh:
		}
	}

	if err != nil {
		return Response{}, err
	}

	return Response{
		BnbusdPrice: bnbusdCurrentRate,
		EthusdPrice: ethusdCurrentRate,
		BscFees:     bscFees,
		EthFees:     ethFees,
	}, nil
}

func calculateBscFees(o *operator, account string) (float64, error) {
	transactions, err := o.bscService.GetAccountTransactions(account)
	if err != nil {
		return 0, fmt.Errorf("something went terribly wrong with BSC/ETH networking: %w", err)
	}
	fees := o.calculator.CalculateFees(transactions)

	f, _ := fees.Float64()
	return f, nil
}

func calculateEthFees(o *operator, account string) (float64, error) {
	transactions, err := o.ethService.GetAccountTransactions(account)
	if err != nil {
		return 0, fmt.Errorf("something went terribly wrong with ETH networking: %w", err)
	}
	fees := o.calculator.CalculateFees(transactions)

	f, _ := fees.Float64()
	return f, nil
}

func GetLatestCurrencyRate(o *operator, symbol string) (float64, error) {
	currentRate, err := o.binanceService.GetLatestCurrencyRate(symbol)
	if err != nil {
		return 0, fmt.Errorf("something went terribly wrong with Binance networking: %w", err)
	}
	return currentRate, nil
}
