package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/withops/trading-tracker/internal/api"
	"github.com/withops/trading-tracker/internal/binance"
	"github.com/withops/trading-tracker/internal/service"
)

func loadEnv() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	// Load .env file if in development
	if env == "development" {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: No .env file found")
		}
	}
}

func main() {
	// Load environment variables
	loadEnv()

	// Initialize Binance client
	apiKey := os.Getenv("BINANCE_API_KEY")
	secretKey := os.Getenv("BINANCE_SECRET_KEY")

	if apiKey == "" || secretKey == "" {
		log.Fatal("BINANCE_API_KEY and BINANCE_SECRET_KEY must be set in the environment")
	}

	binanceClient := binance.NewClient(apiKey, secretKey)

	// Initialize services
	tradeService := service.NewTradeService(binanceClient)

	// Initialize API handlers
	handler := api.NewHandler(tradeService)

	// Create Gin router
	r := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	
	// Get allowed origins from environment variable or use default
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://localhost:3000"
	}
	config.AllowOrigins = strings.Split(allowedOrigins, ",")
	
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Log environment
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	log.Printf("Starting server in %s mode", env)

	// API Routes
	apiGroup := r.Group("/api")
	{
		// Health check
		apiGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})

		// Trade endpoints
		apiGroup.GET("/trades", handler.GetTrades)
		apiGroup.POST("/trades/:id/note", handler.UpdateTradeNote)
		
		// PnL endpoints
		apiGroup.GET("/pnl", handler.GetPnL)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Server starting on :%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
