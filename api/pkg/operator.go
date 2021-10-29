package pkg

import (
	"bsc-fees/pkg/binance"
	"bsc-fees/pkg/config"
	"bsc-fees/pkg/eth"
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

type currencyType int

const (
	bscFees = iota
	ethFees
	bnbusd
	ethusd
)

type price struct {
	value   float64
	curType currencyType
}

func NewOperator(ctx context.Context, cfg config.Config) Operator {
	bscClient := net.NewClient(cfg)
	ethClient := net.NewClient(cfg)
	binanceClient := net.NewClient(cfg)
	bscService := eth.NewProviderService(bscClient, cfg.Bsc)
	ethService := eth.NewProviderService(ethClient, cfg.Eth)
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

func NewManualOperator(
	ctx context.Context,
	cfg config.Config,
	bscService eth.TxGetter,
	ethService eth.TxGetter,
	binanceService binance.RateGetter) Operator {
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
	currencyCh := make(chan price)
	errorCh := make(chan error, concurrentRoutinesNum)

	go func() {
		fees, err := calculateBscFees(o, account)
		if err != nil {
			errorCh <- err
			return
		}
		currencyCh <- price{
			value:   fees,
			curType: bscFees,
		}
	}()

	go func() {
		fees, err := calculateEthFees(o, account)
		if err != nil {
			errorCh <- err
			return
		}
		currencyCh <- price{
			value:   fees,
			curType: ethFees,
		}
	}()

	go func() {
		currencyRate, err := GetLatestCurrencyRate(o, o.cfg.Binance.BnbusdCurrencySymbol)
		if err != nil {
			errorCh <- err
			return
		}
		currencyCh <- price{
			value:   currencyRate,
			curType: bnbusd,
		}
	}()

	go func() {
		currencyRate, err := GetLatestCurrencyRate(o, o.cfg.Binance.EthusdCurrencySymbol)
		if err != nil {
			errorCh <- err
			return
		}
		currencyCh <- price{
			value:   currencyRate,
			curType: ethusd,
		}
	}()

	var bscFeesPrice float64
	var ethFeesPrice float64
	var bnbusdCurrencyRate float64
	var ethusdCurrencyRate float64
	var c price
	var err error
	for i := 0; i < concurrentRoutinesNum; i++ {
		select {
		case c = <-currencyCh:
			switch c.curType {
			case bnbusd:
				bnbusdCurrencyRate = c.value
			case ethusd:
				ethusdCurrencyRate = c.value
			case bscFees:
				bscFeesPrice = c.value
			case ethFees:
				ethFeesPrice = c.value
			}
		case err = <-errorCh:
		}
	}

	err = determineError(bscFeesPrice, ethFeesPrice, bnbusdCurrencyRate, ethusdCurrencyRate, err)

	if err != nil {
		return Response{}, err
	}

	return Response{
		BnbusdPrice: bnbusdCurrencyRate,
		EthusdPrice: ethusdCurrencyRate,
		BscFees:     bscFeesPrice,
		EthFees:     ethFeesPrice,
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

// Sometimes a wallet might be valid for eth but not for bsc and vice versa.
// Therefore we need to be resilient in such cases.
// Maybe a smarter way would be more appropriate here such as better design in channels and go routines
func determineError(bscFees float64, ethFees float64, bnbusdRate float64, ethusdRate float64, err error) error {
	if err == nil {
		return nil
	}

	if bnbusdRate == 0 || ethusdRate == 0 {
		return fmt.Errorf("could not determine currency rates: %w", err)
	}

	if bscFees == 0 && ethFees == 0 {
		return fmt.Errorf("could not determine fees: %w", err)
	}

	return nil
}
