package strategy

import (
    "github.com/AndoNorth/go-quant/internal/models"
)

type MeanReversion struct {
    shortWindow int
    longWindow  int
    prices      []float64
    position    string
}

func NewMeanReversion(short, long int) *MeanReversion {
    return &MeanReversion{
        shortWindow: short,
        longWindow:  long,
        prices:      []float64{},
        position:    "",
    }
}

func (s *MeanReversion) OnTick(t models.Tick) string {
    s.prices = append(s.prices, t.Price)
    if len(s.prices) < s.longWindow {
        return ""
    }

    shortAvg := avg(s.prices[len(s.prices)-s.shortWindow:])
    longAvg := avg(s.prices[len(s.prices)-s.longWindow:])

    if shortAvg < longAvg && s.position != "LONG" {
        s.position = "LONG"
        return "BUY"
    } else if shortAvg > longAvg && s.position != "SHORT" {
        s.position = "SHORT"
        return "SELL"
    }
    return ""
}

func avg(vals []float64) float64 {
    sum := 0.0
    for _, v := range vals {
        sum += v
    }
    return sum / float64(len(vals))
}

