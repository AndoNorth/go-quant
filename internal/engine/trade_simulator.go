package engine

import (
    "fmt"
    "sync"
    "time"

    "github.com/AndoNorth/go-quant/internal/models"
    "github.com/AndoNorth/go-quant/internal/storage"
)

type TradeSimulator struct {
    mu          sync.Mutex
    position    float64
    entryPrice  float64
    realizedPnL float64
    trades      []models.Trade
    store       *storage.SQLiteStore
}

func NewTradeSimulator(store *storage.SQLiteStore) *TradeSimulator {
    return &TradeSimulator{
        trades: make([]models.Trade, 0),
        store:  store,
    }
}

func (t *TradeSimulator) ExecuteTrade(symbol, side string, price float64, qty float64) {
    t.mu.Lock()
    defer t.mu.Unlock()

    var pnl float64
    switch side {
    case "BUY":
        // closing a short
        if t.position < 0 {
            pnl = (t.entryPrice - price) * -t.position
            t.realizedPnL += pnl
            t.position = 0
        }
        // opening a new long
        if t.position == 0 {
            t.entryPrice = price
            t.position = qty
        }

    case "SELL":
        // closing a long
        if t.position > 0 {
            pnl = (price - t.entryPrice) * t.position
            t.realizedPnL += pnl
            t.position = 0
        }
        // opening a new short
        if t.position == 0 {
            t.entryPrice = price
            t.position = -qty
        }
    }

    trade := models.Trade{
        Symbol:     symbol,
        Side:       side,
        Price:      price,
        Quantity:   qty,
        Timestamp:  time.Now(),
        RealizedPnL: t.realizedPnL,
    }
    t.trades = append(t.trades, trade)
    _ = t.store.SaveTrade(trade) // save to DB

    fmt.Printf("Executed %s %.4f %s @ %.2f | RealizedPnL: %.2f\n",
        side, qty, symbol, price, t.realizedPnL)
}

func (t *TradeSimulator) GetPnL() float64 {
    t.mu.Lock()
    defer t.mu.Unlock()
    return t.realizedPnL
}

