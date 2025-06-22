# Trading Tracker - Binance PnL Journal

A full-stack trading journal and PnL tracker for Binance Futures.

## Features

- 📊 Track your Binance Futures trades automatically
- 📈 View daily/weekly PnL and performance metrics
- ✏️ Add notes and labels to trades
- 🚀 Modern web interface with SvelteKit
- ☁️ Ready for cloud deployment with Kubernetes

## Tech Stack

- **Backend**: Go 1.22
- **Frontend**: SvelteKit
- **Deployment**: Helm, Kubernetes
- **CI/CD**: GitHub Actions

## Development

### Prerequisites

- Go 1.22+
- Node.js 18+
- Docker (for containerization)

### Getting Started

1. Clone the repository
2. Set up the backend:
   ```bash
   cd backend
   go mod download
   ```
3. Set up the frontend:
   ```bash
   cd frontend
   npm install
   ```
4. Copy `.env.example` to `.env` and configure your environment variables
5. Run the development servers

## License

MIT
