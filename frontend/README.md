# Trading Tracker Frontend

This is the frontend for the Trading Tracker application, built with SvelteKit and TypeScript.

## Features

- View trade history with filtering and sorting
- Track PnL (Profit and Loss) with interactive charts
- Add notes to trades
- Responsive design that works on desktop and mobile

## Prerequisites

- Node.js 16+ and npm 8+
- Backend API server (see backend README for setup)

## Getting Started

1. Install dependencies:

```bash
npm install
```

2. Copy the example environment file:

```bash
cp .env.example .env
```

3. Update the environment variables in `.env` to match your backend API URL.

## Development

Start the development server:

```bash
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```bash
npm run build
```

You can preview the production build with `npm run preview`.

## Testing

Run the test suite:

```bash
npm test
```

## Environment Variables

- `PUBLIC_API_URL`: The URL of the backend API (default: http://localhost:8080/api)

```bash
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Building

To create a production version of your app:

```bash
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.
