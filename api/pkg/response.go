package pkg

type Response struct {
	BnbusdPrice float64 `json:"bnbusdPrice"`
	EthusdPrice float64 `json:"ethusdPrice"`
	BscFees     float64 `json:"bscFees"`
	EthFees     float64 `json:"ethFees"`
}
