package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/AndoNorth/go-quant/internal/datafeed"
    "github.com/AndoNorth/go-quant/internal/models"
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

    for {
        select {
        case <-ctx.Done():
            return
        case tick := <-out:
            fmt.Printf("[%s] %.2f @ %.3f | %s\n",
                tick.Symbol, tick.Price, tick.Volume, tick.Timestamp.Format(time.RFC3339))
        }
    }
}

