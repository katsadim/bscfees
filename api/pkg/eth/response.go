package eth

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

type Envelope struct {
	Status  int    `json:"status,string"`
	Message string `json:"message"`
	//Result  []NormalTx `json:"result"`
	Result json.RawMessage `json:"result"`
}

// NormalTx holds info from normal tx query
type NormalTx struct {
	BlockNumber       int    `json:"blockNumber,string"`
	TimeStamp         Time   `json:"timeStamp"`
	Hash              string `json:"hash"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  int    `json:"transactionIndex,string"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             BigInt `json:"value"`
	Gas               int    `json:"gas,string"`
	GasPrice          BigInt `json:"gasPrice"`
	IsError           int    `json:"isError,string"`
	TxReceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed int    `json:"cumulativeGasUsed,string"`
	GasUsed           int    `json:"gasUsed,string"`
	Confirmations     int    `json:"confirmations,string"`
}

// Time is a wrapper over time.Time to implement only unmarshalText
// for json decoding.
type Time time.Time

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (t *Time) UnmarshalText(text []byte) (err error) {
	input, err := strconv.ParseInt(string(text), 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing time: %w", err)
	}

	var timestamp = time.Unix(input, 0)
	*t = Time(timestamp)

	return nil
}

// Time returns t's time.Time form
func (t Time) Time() time.Time {
	return time.Time(t)
}

// MarshalText implements the encoding.TextMarshaler
func (t Time) MarshalText() (text []byte, err error) {
	return []byte(strconv.FormatInt(t.Time().Unix(), 10)), nil
}

// BigInt is a wrapper over big.Int to implement only unmarshalText
// for json decoding.
type BigInt big.Int

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *BigInt) UnmarshalText(text []byte) (err error) {
	var bigInt = new(big.Int)
	err = bigInt.UnmarshalText(text)
	if err != nil {
		return
	}

	*b = BigInt(*bigInt)
	return nil
}

// MarshalText implements the encoding.TextMarshaler
func (b *BigInt) MarshalText() (text []byte, err error) {
	return []byte(b.Int().String()), nil
}

// Int returns b's *big.Int form
func (b *BigInt) Int() *big.Int {
	return (*big.Int)(b)
}
