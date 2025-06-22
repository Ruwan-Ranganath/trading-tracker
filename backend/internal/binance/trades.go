package binance

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Trade represents a single trade from Binance
type Trade struct {
	ID              int64   `json:"id"`
	Symbol          string  `json:"symbol"`
	OrderID         int64   `json:"orderId"`
	Price           float64 `json:"price,string"`
	Qty            float64 `json:"qty,string"`
	QuoteQty       float64 `json:"quoteQty,string"`
	Commission     float64 `json:"commission,string"`
	CommissionAsset string  `json:"commissionAsset"`
	Time           int64   `json:"time"`
	IsBuyer        bool    `json:"isBuyer"`
	IsMaker        bool    `json:"isMaker"`
	IsBestMatch    bool    `json:"isBestMatch"`
}

// GetTradesOptions defines the parameters for fetching trades
type GetTradesOptions struct {
	Symbol    string
	StartTime *time.Time
	EndTime   *time.Time
	Limit     int
}

// GetTrades fetches the user's trades for a specific symbol
func (c *Client) GetTrades(ctx context.Context, opts GetTradesOptions) ([]Trade, error) {
	if opts.Symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	if opts.Limit <= 0 || opts.Limit > 1000 {
		opts.Limit = 1000
	}

	params := url.Values{}
	params.Add("symbol", opts.Symbol)
	params.Add("limit", strconv.Itoa(opts.Limit))

	if opts.StartTime != nil {
		params.Add("startTime", strconv.FormatInt(opts.StartTime.UnixMilli(), 10))
	}

	if opts.EndTime != nil {
		params.Add("endTime", strconv.FormatInt(opts.EndTime.UnixMilli(), 10))
	}

	// Add signature for private endpoint
	err := c.signRequest(params)
	if err != nil {
		return nil, fmt.Errorf("error signing request: %w", err)
	}

	endpoint := "/fapi/v1/userTrades" + "?" + params.Encode()

	var trades []Trade
	err = c.doRequest(ctx, "GET", endpoint, nil, &trades)
	if err != nil {
		return nil, fmt.Errorf("error fetching trades: %w", err)
	}

	return trades, nil
}

// GetAccountTrades fetches all trades for the account
func (c *Client) GetAccountTrades(ctx context.Context, symbol string, limit int) ([]Trade, error) {
	return c.GetTrades(ctx, GetTradesOptions{
		Symbol: symbol,
		Limit:  limit,
	})
}

// GetRecentTrades fetches recent trades for a symbol
func (c *Client) GetRecentTrades(ctx context.Context, symbol string, limit int) ([]Trade, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	if limit <= 0 || limit > 1000 {
		limit = 1000
	}

	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("limit", strconv.Itoa(limit))

	endpoint := "/fapi/v1/trades" + "?" + params.Encode()

	var trades []Trade
	err := c.doRequest(ctx, "GET", endpoint, nil, &trades)
	if err != nil {
		return nil, fmt.Errorf("error fetching recent trades: %w", err)
	}

	return trades, nil
}
