package strategy

import "github.com/AndoNorth/go-quant/internal/models"

type Strategy interface {
    OnTick(t models.Tick) (signal string) // returns "BUY", "SELL", or ""
}

