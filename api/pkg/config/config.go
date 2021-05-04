package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Config struct {
	HttpClient HTTPClient
	Bsc        Provider
	Eth        Provider
	Binance    Binance
	Options    CalculatorOptions
	General    GeneralOptions
}

type HTTPClient struct {
	Timeout   time.Duration
	KeepAlive bool
}

type Provider struct {
	BaseURI           string
	Module            string
	Action            string
	ApiKey            string
	RateLimit         int
	RateLimitDuration time.Duration
}

type Binance struct {
	BaseURI               string
	CandlesEndpoint       string
	CurrentPriceEndpoint  string
	ApiKey                string
	BnbusdCurrencySymbol  string
	EthusdCurrencySymbol  string
	RateLimit             int
	RateLimitDuration     time.Duration
}

type CalculatorOptions struct {
	Currency string
}

type GeneralOptions struct {
	Env string
}

// Load parses the application configuration file and stores the values to a struct
func Load() (Config, error) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigName("config")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	c := Config{}
	if err := viper.Unmarshal(&c); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return c, nil
}
