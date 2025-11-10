package storage

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "modernc.org/sqlite"
    "github.com/AndoNorth/go-quant/internal/models"
)

type SQLiteStore struct {
    db *sql.DB
}

func NewSQLiteStore(path string) *SQLiteStore {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        fmt.Println("Creating new SQLite database:", path)
    }

    db, err := sql.Open("sqlite", path)
    if err != nil {
        log.Fatalf("failed to open SQLite DB: %v", err)
    }

    store := &SQLiteStore{db: db}
    store.initSchema()
    return store
}

func (s *SQLiteStore) initSchema() {
    _, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS trades (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            symbol TEXT,
            side TEXT,
            price REAL,
            quantity REAL,
            timestamp TEXT,
            realized_pnl REAL
        );
    `)
    if err != nil {
        log.Fatalf("failed to create schema: %v", err)
    }
}

func (s *SQLiteStore) SaveTrade(t models.Trade) error {
    _, err := s.db.Exec(`
        INSERT INTO trades (symbol, side, price, quantity, timestamp, realized_pnl)
        VALUES (?, ?, ?, ?, ?, ?);
    `, t.Symbol, t.Side, t.Price, t.Quantity, t.Timestamp.Format("2006-01-02T15:04:05Z07:00"), t.RealizedPnL)
    return err
}

func (s *SQLiteStore) GetAllTrades() ([]models.Trade, error) {
    rows, err := s.db.Query("SELECT symbol, side, price, quantity, timestamp, realized_pnl FROM trades ORDER BY id;")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    trades := []models.Trade{}
    for rows.Next() {
        var t models.Trade
        var ts string
        err := rows.Scan(&t.Symbol, &t.Side, &t.Price, &t.Quantity, &ts, &t.RealizedPnL)
        if err != nil {
            return nil, err
        }
        trades = append(trades, t)
    }
    return trades, nil
}
