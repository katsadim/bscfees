package binance

import (
	"bsc-fees/pkg/config"
	http2 "bsc-fees/pkg/net"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type RateGetter interface {
	GetCurrencyRates(times []time.Time) (map[time.Time]float64, error)
	GetLatestCurrencyRate() (float64, error)
}

type binanceService struct {
	client http2.Client
	cfg    config.Binance
}

func NewBinanceService(c http2.Client, cfg config.Binance) RateGetter {
	return &binanceService{
		client: c,
		cfg:    cfg}
}

func (bs *binanceService) GetCurrencyRates(times []time.Time) (map[time.Time]float64, error) {

	res := make(map[time.Time]float64, len(times))

	for i, t := range times {
		rate, err := getCurrencyForTime(t, bs.cfg, bs.client)
		if err != nil {
			return map[time.Time]float64{}, fmt.Errorf("request #%d failed with err: %w", i, err)
		}
		res[t] = rate
	}

	return res, nil
}

func (bs *binanceService) GetLatestCurrencyRate() (float64, error) {
	u, err := url.Parse(synthesizeCurrentPriceUri(bs.cfg))
	if err != nil {
		return 0.0, err
	}

	resp, err := bs.client.SendRequest(u, map[string]string{"X-MBX-APIKEY": bs.cfg.ApiKey})
	if err != nil {
		return 0.0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusAccepted {
		return 0.0, fmt.Errorf("expected 2XX but instead got a %d", resp.StatusCode)
	}

	price := CurrentAveragePrice{}
	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		return 0.0, err
	}

	return price.Price, nil

}

func getCurrencyForTime(t time.Time, cfg config.Binance, client http2.Client) (float64, error) {
	u, err := url.Parse(synthesizeCandlesUri(cfg, t))
	if err != nil {
		return 0.0, err
	}

	resp, err := client.SendRequest(u, map[string]string{"X-MBX-APIKEY": cfg.ApiKey})
	if err != nil {
		return 0.0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusAccepted {
		return 0.0, fmt.Errorf("expected 2XX but instead got a %d", resp.StatusCode)
	}

	candles := CandleSticks{}
	if err := json.NewDecoder(resp.Body).Decode(&candles); err != nil {
		return 0.0, err
	}

	if len(candles) == 0 {
		return 0.0, fmt.Errorf("no candles returned")
	}

	candle := pickNearestCandle(t, candles)

	highStr := candle[HighIndex].(string)
	high, err := strconv.ParseFloat(highStr, 16)
	if err != nil {
		return 0.0, fmt.Errorf("parse of float failed: %v", err)
	}

	lowStr := candle[LowIndex].(string)
	low, err := strconv.ParseFloat(lowStr, 16)
	if err != nil {
		return 0.0, fmt.Errorf("parse of float failed: %v", err)
	}

	avg := (high + low) / 2

	return avg, nil
}

// pickNearestCandle selects the candle that includes the time passed as argument
func pickNearestCandle(t time.Time, candles CandleSticks) CandleStick {
	convertedTime := t.UTC().UnixNano() / int64(time.Millisecond)
	ret := candles[0]
	for _, candle := range candles {
		if convertedTime >= int64(candle[OpenTimeIndex].(float64)) && convertedTime <= int64(candle[CloseTimeIndex].(float64)) {
			ret = candle
		}
	}
	return ret
}

func synthesizeCandlesUri(cfg config.Binance, t time.Time) string {
	startTime := t.Add(-30*time.Second).UTC().UnixNano() / int64(time.Millisecond)
	endTime := t.Add(30*time.Second).UTC().UnixNano() / int64(time.Millisecond)

	return fmt.Sprintf("%s%s?symbol=%s&startTime=%d&endTime=%d&interval=1m", cfg.BaseURI, cfg.CandlesEndpoint, cfg.CurrencySymbol, startTime, endTime)
}

func synthesizeCurrentPriceUri(cfg config.Binance) string {
	return fmt.Sprintf("%s%s?symbol=%s", cfg.BaseURI, cfg.CurrentPriceEndpoint, cfg.CurrencySymbol)
}
