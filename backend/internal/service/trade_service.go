package service

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/withops/trading-tracker/internal/binance"
)

type TradeService struct {
	binanceClient *binance.Client
}

func NewTradeService(binanceClient *binance.Client) *TradeService {
	return &TradeService{
		binanceClient: binanceClient,
	}
}

type Trade struct {
	ID              string    `json:"id"`
	Symbol          string    `json:"symbol"`
	Price           float64   `json:"price"`
	Quantity        float64   `json:"quantity"`
	QuoteQty       float64   `json:"quoteQty"`
	Commission     float64   `json:"commission"`
	CommissionAsset string    `json:"commissionAsset"`
	Time           time.Time `json:"time"`
	IsBuyer        bool      `json:"isBuyer"`
	IsMaker        bool      `json:"isMaker"`
	Note           string    `json:"note,omitempty"`
}

func (s *TradeService) GetTrades(ctx context.Context, symbol string, limit int) ([]Trade, error) {
	// First try to get recent trades to verify API access
	recentTrades, err := s.binanceClient.GetRecentTrades(ctx, symbol, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent trades (public endpoint): %w", err)
	}

	log.Printf("Recent trades (public): %+v", recentTrades)

	// Then get account trades
	binanceTrades, err := s.binanceClient.GetAccountTrades(ctx, symbol, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get account trades: %w", err)
	}

	log.Printf("Account trades: %d found", len(binanceTrades))

	trades := make([]Trade, 0, len(binanceTrades))
	for _, t := range binanceTrades {
		trades = append(trades, Trade{
			ID:              strconv.FormatInt(t.ID, 10),
			Symbol:          t.Symbol,
			Price:           t.Price,
			Quantity:        t.Qty,
			QuoteQty:       t.QuoteQty,
			Commission:     t.Commission,
			CommissionAsset: t.CommissionAsset,
			Time:           time.Unix(0, t.Time*int64(time.Millisecond)),
			IsBuyer:        t.IsBuyer,
			IsMaker:        t.IsMaker,
		})
	}

	return trades, nil
}

func (s *TradeService) GetPnL(ctx context.Context, symbol string, days int) ([]PnLEntry, error) {
	if days <= 0 {
		days = 7 // Default to 7 days
	}

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -days)

	trades, err := s.binanceClient.GetTrades(ctx, binance.GetTradesOptions{
		Symbol:    symbol,
		StartTime: &startTime,
		EndTime:   &endTime,
		Limit:     1000,
	})

	if err != nil {
		return nil, err
	}

	// Group trades by day and calculate PnL
	pnlByDay := make(map[string]float64)
	for _, trade := range trades {
		tradeTime := time.Unix(0, trade.Time*int64(time.Millisecond))
		day := tradeTime.Format("2006-01-02")

		// Calculate PnL for this trade (simplified)
		// In a real implementation, you'd need to consider position sizing, fees, etc.
		pnl := 0.0
		if trade.IsBuyer {
			pnl = -trade.QuoteQty // Negative for buys
		} else {
			pnl = trade.QuoteQty // Positive for sells
		}

		pnlByDay[day] += pnl
	}

	// Convert to sorted slice
	var pnlData []PnLEntry
	for date, pnl := range pnlByDay {
		t, _ := time.Parse("2006-01-02", date)
		pnlData = append(pnlData, PnLEntry{
			Date: t,
			PnL:  pnl,
		})
	}

	// Sort by date
	sort.Slice(pnlData, func(i, j int) bool {
		return pnlData[i].Date.Before(pnlData[j].Date)
	})

	return pnlData, nil
}

type PnLEntry struct {
	Date time.Time `json:"date"`
	PnL  float64   `json:"pnl"`
}
