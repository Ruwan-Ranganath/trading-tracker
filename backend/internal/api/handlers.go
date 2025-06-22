package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/withops/trading-tracker/internal/service"
)

type Handler struct {
	tradeService *service.TradeService
}

func NewHandler(tradeService *service.TradeService) *Handler {
	return &Handler{
		tradeService: tradeService,
	}
}

// GetTrades returns a list of trades
func (h *Handler) GetTrades(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		symbol = "BTCUSDT" // Default symbol
	}

	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	trades, err := h.tradeService.GetTrades(c.Request.Context(), symbol, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch trades: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": trades,
	})
}

// UpdateTradeNote updates the note for a trade
func (h *Handler) UpdateTradeNote(c *gin.Context) {
	tradeID := c.Param("id")
	if tradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "trade ID is required"})
		return
	}

	var request struct {
		Note string `json:"note"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// In a real implementation, you would update the note in your database
	// For now, we'll just return a success message

	c.JSON(http.StatusOK, gin.H{
		"message": "Trade note updated successfully",
	})
}

// GetPnL returns PnL data
func (h *Handler) GetPnL(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		symbol = "BTCUSDT" // Default symbol
	}

	days := 7
	if daysStr := c.Query("days"); daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	pnlData, err := h.tradeService.GetPnL(c.Request.Context(), symbol, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch PnL data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": pnlData,
	})
}
