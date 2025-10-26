package models

import "time"

type Tick struct {
    Symbol    string
    Price     float64
    Volume    float64
    Timestamp time.Time
}

