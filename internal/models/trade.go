package models

import "time"

type Trade struct {
	Symbol string
	Side string // BUY or SELL
	Price float64
	Quantity float64
	Timestamp time.Time
	RealizedPnL float64
}
