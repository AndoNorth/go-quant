# go-quant

## Project Overview
```
go-quant/
├── cmd/
│   ├── aggregator/         # Entry point for market data aggregator
│   │   └── main.go
│   ├── backtester/         # Entry point for offline backtesting
│   │   └── main.go
│   └── strategy_runner/    # Entry point for live strategy
│       └── main.go
│
├── internal/
│   ├── datafeed/           # Exchange/WebSocket connectors
│   │   ├── binance.go
│   │   └── polygon.go
│   │
│   ├── models/             # Core data structures
│   │   ├── tick.go
│   │   ├── order.go
│   │   └── trade.go
│   │
│   ├── strategy/           # Strategies implement a common interface
│   │   ├── interface.go
│   │   ├── mean_reversion.go
│   │   └── momentum.go
│   │
│   ├── engine/             # Core orchestration logic
│   │   ├── aggregator.go
│   │   ├── backtester.go
│   │   ├── runner.go
│   │   └── metrics.go
│   │
│   ├── storage/            # Data persistence (SQLite, CSV, etc.)
│   │   ├── sqlite.go
│   │   ├── csv.go
│   │   └── interface.go
│   │
│   └── utils/              # Shared helpers (logging, config)
│       ├── config.go
│       └── logger.go
│
├── pkg/                    # Public reusable packages (optional)
│   └── indicators/         # Moving averages, RSI, etc.
│       ├── sma.go
│       ├── ema.go
│       └── bollinger.go
│
├── configs/                # YAML/TOML config files
│   ├── aggregator.yaml
│   ├── strategy.yaml
│   └── db.yaml
│
├── scripts/                # Utility scripts (data import, setup)
│   └── import_data.sh
│
├── testdata/               # Sample CSV price data for backtesting
│   └── BTCUSDT_1h.csv
│
├── go.mod
├── go.sum
└── README.md
```
