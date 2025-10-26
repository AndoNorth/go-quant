package datafeed

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/AndoNorth/go-quant/internal/models"
    "nhooyr.io/websocket"
)

type BinanceFeed struct {
    Symbols []string
}

func NewBinanceFeed(symbols []string) *BinanceFeed {
    return &BinanceFeed{Symbols: symbols}
}

func (b *BinanceFeed) Start(ctx context.Context, out chan<- models.Tick) error {
    for _, sym := range b.Symbols {
        go b.subscribeSymbol(ctx, sym, out)
    }
    return nil
}

func (b *BinanceFeed) subscribeSymbol(ctx context.Context, symbol string, out chan<- models.Tick) {
    url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@trade", toLower(symbol))
    c, _, err := websocket.Dial(ctx, url, nil)
    if err != nil {
        log.Printf("error connecting to Binance: %v", err)
        return
    }
    defer c.Close(websocket.StatusNormalClosure, "done")

    for {
        _, data, err := c.Read(ctx)
        if err != nil {
            log.Printf("[%s] read error: %v", symbol, err)
            return
        }

        var msg struct {
            Price string `json:"p"`
            Vol   string `json:"q"`
            Ts    int64  `json:"T"`
        }
        if err := json.Unmarshal(data, &msg); err != nil {
            continue
        }

        price, _ := parseFloat(msg.Price)
        vol, _ := parseFloat(msg.Vol)

        tick := models.Tick{
            Symbol:    symbol,
            Price:     price,
            Volume:    vol,
            Timestamp: time.UnixMilli(msg.Ts),
        }

        select {
        case out <- tick:
        case <-ctx.Done():
            return
        }
    }
}

// --- helpers ---
func toLower(sym string) string {
    return fmt.Sprintf("%s", lowerNoSlash(sym))
}

func lowerNoSlash(s string) string {
    out := ""
    for _, c := range s {
        if c != '/' {
            out += string(toLowerRune(c))
        }
    }
    return out
}

func toLowerRune(r rune) rune {
    if r >= 'A' && r <= 'Z' {
        return r + 32
    }
    return r
}

func parseFloat(s string) (float64, error) {
    var f float64
    _, err := fmt.Sscan(s, &f)
    return f, err
}

