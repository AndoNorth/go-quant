package engine

import (
    "fmt"
    "sync"
    "time"

    "github.com/AndoNorth/go-quant/internal/models"
    "github.com/AndoNorth/go-quant/internal/storage"
)

type TradeSimulator struct {
    mu            sync.Mutex
    position      float64
    entryPrice    float64
    realizedPnL   float64
    unrealizedPnL float64
    feeRate       float64
    trades        []models.Trade
    store         *storage.SQLiteStore
}

func NewTradeSimulator(store *storage.SQLiteStore) *TradeSimulator {
    return &TradeSimulator{
        trades:  make([]models.Trade, 0),
        store:   store,
        feeRate: 0.001, // 0.1% fee
    }
}

func (t *TradeSimulator) ExecuteTrade(symbol, side string, price float64, qty float64) {
    t.mu.Lock()
    defer t.mu.Unlock()

    fee := price * qty * t.feeRate
    pnl := 0.0

    switch side {
    case "BUY":
        if t.position < 0 { // closing short
            pnl = (t.entryPrice - price) * (-t.position)
            t.realizedPnL += pnl - fee
            t.position = 0
        } else { // open or add to long
            t.entryPrice = ((t.entryPrice * t.position) + (price * qty)) / (t.position + qty)
            t.position += qty
            t.realizedPnL -= fee
        }

    case "SELL":
        if t.position > 0 { // closing long
            pnl = (price - t.entryPrice) * t.position
            t.realizedPnL += pnl - fee
            t.position = 0
        } else { // open or add to short
            t.entryPrice = ((t.entryPrice * (-t.position)) + (price * qty)) / ((-t.position) + qty)
            t.position -= qty
            t.realizedPnL -= fee
        }
    }

    trade := models.Trade{
        Symbol:      symbol,
        Side:        side,
        Price:       price,
        Quantity:    qty,
        Timestamp:   time.Now(),
        RealizedPnL: t.realizedPnL,
    }

    t.trades = append(t.trades, trade)
    _ = t.store.SaveTrade(trade)

    fmt.Printf("Executed %s %.4f %s @ %.2f | RealizedPnL: %.2f | Fee: %.4f\n",
        side, qty, symbol, price, t.realizedPnL, fee)
}

// Call this on every tick even if no trade executes
func (t *TradeSimulator) UpdateUnrealizedPnL(currentPrice float64) {
    t.mu.Lock()
    defer t.mu.Unlock()
    t.unrealizedPnL = (currentPrice - t.entryPrice) * t.position
}

// Total account PnL (realized + unrealized)
func (t *TradeSimulator) GetTotalPnL() float64 {
    t.mu.Lock()
    defer t.mu.Unlock()
    return t.realizedPnL + t.unrealizedPnL
}

func (t *TradeSimulator) GetPnL() float64 {
    t.mu.Lock()
    defer t.mu.Unlock()
    return t.realizedPnL
}
