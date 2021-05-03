package bsc

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

type TxGetter interface {
	GetAccountTransactions(account string) ([]Transaction, error)
}

type bscService struct {
	client net.Client
	cfg    config.BSC
}

func NewBscService(c net.Client, cfg config.BSC) TxGetter {
	return &bscService{
		client: c,
		cfg:    cfg}
}

func (s *bscService) GetAccountTransactions(address string) ([]Transaction, error) {
	u, err := url.Parse(synthesizeBSCUri(s.cfg, address))
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

func synthesizeBSCUri(cfg config.BSC, address string) string {
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
