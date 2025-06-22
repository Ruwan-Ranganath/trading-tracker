package model

import "time"

type Trade struct {
	ID            string    `json:"id"`
	Symbol        string    `json:"symbol"`
	Side         string    `json:"side"`
	Price        float64   `json:"price,string"`
	Quantity     float64   `json:"quantity,string"`
	QuoteQty     float64   `json:"quoteQty,string"`
	Commission   float64   `json:"commission,string"`
	CommissionAsset string  `json:"commissionAsset"`
	Time         time.Time `json:"time"`
	IsBuyer      bool      `json:"isBuyer"`
	IsMaker      bool      `json:"isMaker"`
	IsBestMatch  bool      `json:"isBestMatch"`
	Note         string    `json:"note,omitempty"`
}
