package binance

const (
	OpenTimeIndex  = 0
	CloseTimeIndex = 6
	HighIndex      = 2
	LowIndex       = 3
)

type CandleSticks []CandleStick

type CandleStick []interface{}

type CurrentAveragePrice struct {
	Mins  int     `json:"mins"`
	Price float64 `json:"price,string"`
}
