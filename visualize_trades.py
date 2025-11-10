import sqlite3
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

# Connect to SQLite DB
conn = sqlite3.connect("trades.db")
df = pd.read_sql("SELECT * FROM trades ORDER BY timestamp;", conn)
conn.close()

if df.empty:
    print("No trades found in database.")
    exit()

# Convert timestamp to datetime
df["timestamp"] = pd.to_datetime(df["timestamp"])
df["price"] = df["price"].astype(float)
df["quantity"] = df["quantity"].astype(float)
df["realized_pnl"] = df["realized_pnl"].astype(float)

# --- Derived columns ---
fee_rate = 0.001  # 0.1% per trade
df["fee"] = df["price"] * df["quantity"] * fee_rate

# Adjust realized PnL for fees
df["fee_adjusted_pnl"] = df["realized_pnl"] - df["fee"].cumsum()

# Compute cumulative (realized) PnL curve
df["cum_pnl"] = df["realized_pnl"]

# Compute PnL changes per trade
df["pnl_change"] = df["realized_pnl"].diff().fillna(0)

# --- Performance Metrics ---
returns = df["pnl_change"]

total_pnl = df["realized_pnl"].iloc[-1]
total_fees = df["fee"].sum()
avg_pnl = returns.mean()
std_pnl = returns.std()
sharpe = (avg_pnl / std_pnl * np.sqrt(252)) if std_pnl > 0 else 0

# Win rate
wins = (returns > 0).sum()
losses = (returns <= 0).sum()
win_rate = wins / (wins + losses) if (wins + losses) > 0 else 0

# Max Drawdown
df["cum_max"] = df["cum_pnl"].cummax()
df["drawdown"] = df["cum_pnl"] - df["cum_max"]
max_drawdown = df["drawdown"].min()

# --- Print summary ---
print("\n=== Strategy Performance Summary ===")
print(f"Total Trades:         {len(df)}")
print(f"Total Realized PnL:   {total_pnl:.2f}")
print(f"Total Fees Paid:      {total_fees:.2f}")
print(f"Net Fee-Adj. PnL:     {total_pnl - total_fees:.2f}")
print(f"Average PnL/trade:    {avg_pnl:.4f}")
print(f"PnL Std Dev:          {std_pnl:.4f}")
print(f"Win Rate:             {win_rate*100:.2f}%")
print(f"Sharpe Ratio:         {sharpe:.2f}")
print(f"Max Drawdown:         {max_drawdown:.2f}")
print("====================================\n")

# --- Visualization ---
plt.figure(figsize=(10, 6))

# PnL curve
plt.plot(df["timestamp"], df["cum_pnl"], label="Cumulative Realized PnL", color="blue")

# Fee-adjusted PnL curve
plt.plot(df["timestamp"], df["fee_adjusted_pnl"], label="Fee-Adjusted PnL", color="orange", linestyle="--")

# Drawdown shaded area
plt.fill_between(df["timestamp"], df["cum_pnl"], df["cum_max"], color="red", alpha=0.1, label="Drawdown")

# Trade markers
buy_points = df[df["side"] == "BUY"]
sell_points = df[df["side"] == "SELL"]
plt.scatter(buy_points["timestamp"], buy_points["price"], marker="^", color="green", label="BUY", alpha=0.7)
plt.scatter(sell_points["timestamp"], sell_points["price"], marker="v", color="red", label="SELL", alpha=0.7)

plt.title("Strategy PnL and Trade Visualization")
plt.xlabel("Time")
plt.ylabel("PnL (USD)")
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.show()
