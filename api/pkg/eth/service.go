package eth

//go:generate mockgen -source=service.go -destination service_mock.go -package eth

import (
	"bsc-fees/pkg/config"
	"bsc-fees/pkg/net"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
)
// This service can be used for both bscscan and ethscan
type TxGetter interface {
	GetAccountTransactions(account string) ([]Transaction, error)
}

type providerService struct {
	client net.Client
	cfg    config.Provider
}

func NewProviderService(c net.Client, cfg config.Provider) TxGetter {
	return &providerService{
		client: c,
		cfg:    cfg}
}

func (s *providerService) GetAccountTransactions(address string) ([]Transaction, error) {
	u, err := url.Parse(synthesizeProviderUri(s.cfg, address))
	if err != nil {
		return []Transaction{}, err
	}

	resp, err := s.client.SendRequest(u, map[string]string{})
	if err != nil {
		return []Transaction{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusAccepted {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, resp.Body)
		if err != nil {
			return []Transaction{}, fmt.Errorf("failed Request with '%d'. and not parsable body", resp.StatusCode)
		} else {
			return []Transaction{}, fmt.Errorf("failed Request with '%d'. and body: '%s' ", resp.StatusCode, buf.String())
		}
	}

	e := Envelope{}
	if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
		return []Transaction{}, err
	}

	if e.Status != 1 {
		return []Transaction{},
			fmt.Errorf("server did not like our request with message: %s and result: %s", e.Message, string(e.Result))
	}

	// decode the result
	var txs []NormalTx
	if err := json.Unmarshal(e.Result, &txs); err != nil {
		return []Transaction{}, err
	}

	transactions := processTransactions(txs)

	return transactions, nil

}

func synthesizeProviderUri(cfg config.Provider, address string) string {
	return fmt.Sprintf("%s?module=%s&action=%s&address=%s&apikey=%s", cfg.BaseURI, cfg.Module, cfg.Action, address, cfg.ApiKey)
}

func processTransactions(normalTxs []NormalTx) []Transaction {

	txs := make([]Transaction, 0, len(normalTxs))
	for _, ntx := range normalTxs {
		tx := Transaction{
			Time:     time.Time(ntx.TimeStamp),
			GasPrice: big.Int(ntx.GasPrice),
			GasUsed:  ntx.GasUsed,
		}
		txs = append(txs, tx)
	}
	return txs
}
