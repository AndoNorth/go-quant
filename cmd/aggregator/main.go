package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/AndoNorth/go-quant/internal/datafeed"
    "github.com/AndoNorth/go-quant/internal/engine"
    "github.com/AndoNorth/go-quant/internal/models"
    "github.com/AndoNorth/go-quant/internal/strategy"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-sigs
        fmt.Println("\nShutting down...")
        cancel()
    }()

    out := make(chan models.Tick, 100)
    feed := datafeed.NewBinanceFeed([]string{"BTCUSDT"})
    feed.Start(ctx, out)

    strat := strategy.NewMeanReversion(5, 20)
    sim := engine.NewTradeSimulator()

    for {
        select {
        case <-ctx.Done():
            fmt.Printf("\nFinal Realized PnL: %.2f\n", sim.GetPnL())
            return
        case tick := <-out:
            signal := strat.OnTick(tick)
            if signal != "" {
                sim.ExecuteTrade(tick.Symbol, signal, tick.Price, 0.01)
            }
        }
    }
}

