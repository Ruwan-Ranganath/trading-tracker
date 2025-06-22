<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { Chart, type ChartConfiguration } from 'chart.js/auto';
  import { api } from '$lib/services/api';
  import type { PnLEntry } from '$lib/types/api';
  
  // Declare Chart on window for any libraries that might need it
  if (typeof window !== 'undefined') {
    (window as any).Chart = Chart;
  }
  
  // Component state
  let loading = true;
  let error: string | null = null;
  let pnlData: PnLEntry[] = [];
  let totalTrades = 0;
  let totalPnL = 0;
  let winRate = 0;
  
  // Chart instance and canvas ref
  let chart: Chart | null = null;
  let chartCanvas: HTMLCanvasElement | null = null;

  onMount(async () => {
    // Load data and render chart
    try {
      // Fetch PnL data
      pnlData = await api.getPnL();
      
      // Calculate metrics (mock data for now)
      totalTrades = 42;
      totalPnL = pnlData.reduce((sum, entry) => sum + entry.pnl, 0);
      winRate = 0.72; // 72% win rate
      
      // Wait for the next tick to ensure DOM is updated
      await tick();
      
      // Wait for the next frame to ensure canvas is in the DOM
      await new Promise(resolve => requestAnimationFrame(resolve));
      
      // Render chart
      renderChart(pnlData);
      
      // Handle window resize
      const handleResize = () => {
        if (chart) {
          chart.resize();
        }
      };
      
      window.addEventListener('resize', handleResize);
      
      // Cleanup
      return () => {
        window.removeEventListener('resize', handleResize);
        if (chart) {
          chart.destroy();
          chart = null;
        }
      };
      
    } catch (e: unknown) {
      const errorMessage = e instanceof Error ? e.message : 'An unknown error occurred';
      error = errorMessage;
      console.error('Error loading data:', e);
    } finally {
      loading = false;
    }
  });

  function renderChart(data: PnLEntry[]) {
    // Wait for the next tick to ensure the canvas is in the DOM
    requestAnimationFrame(() => {
      if (!chartCanvas) {
        console.error('Canvas element not found');
        return;
      }
      
      const ctx = chartCanvas.getContext('2d');
      if (!ctx) {
        console.error('Could not get 2D context');
        return;
      }
      
      // Destroy existing chart if it exists
      if (chart) {
        chart.destroy();
        chart = null;
      }

      // Ensure we have data to display
      if (!data || data.length === 0) {
        console.warn('No data available for chart');
        return;
      }

      try {
        const chartConfig: ChartConfiguration = {
          type: 'line',
          data: {
            labels: data.map(item => item.date),
            datasets: [{
              label: 'PnL (USDT)',
              data: data.map(item => item.pnl),
              borderColor: 'rgb(59, 130, 246)',
              backgroundColor: 'rgba(59, 130, 246, 0.1)',
              tension: 0.3,
              fill: true
            }]
          },
          options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: {
              duration: 0 // Disable animations for better performance
            },
            plugins: {
              title: {
                display: true,
                text: 'PnL (Last 7 Days)'
              },
              tooltip: {
                mode: 'index',
                intersect: false
              }
            },
            scales: {
              y: {
                beginAtZero: false,
                title: {
                  display: true,
                  text: 'USDT'
                }
              },
              x: {
                title: {
                  display: true,
                  text: 'Date'
                }
              }
            }
          }
        };

        chart = new Chart(ctx, chartConfig);
      } catch (e) {
        console.error('Failed to create chart:', e);
      }
    });
  }
</script>

<div class="container mx-auto p-4">
  <h1 class="text-3xl font-bold mb-6">Trading Dashboard</h1>
  
  {#if loading}
    <div class="text-center py-8">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500 mx-auto"></div>
      <p class="mt-2">Loading dashboard...</p>
    </div>
  {:else if error}
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">Error:</strong>
      <span class="block sm:inline"> {error}</span>
    </div>
  {:else}
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <!-- Total Trades -->
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h3 class="text-gray-500 text-sm font-medium">Total Trades</h3>
        <p class="text-3xl font-bold">{totalTrades}</p>
        <p class="text-sm text-gray-500 mt-1">All time</p>
      </div>
      
      <!-- Total PnL -->
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h3 class="text-gray-500 text-sm font-medium">Total PnL</h3>
        <p class={`text-3xl font-bold ${totalPnL >= 0 ? 'text-green-600' : 'text-red-600'}`}>
          {totalPnL >= 0 ? '+' : ''}{totalPnL.toFixed(2)} USDT
        </p>
        <p class="text-sm text-gray-500 mt-1">All time</p>
      </div>
      
      <!-- Win Rate -->
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h3 class="text-gray-500 text-sm font-medium">Win Rate</h3>
        <p class="text-3xl font-bold">{(winRate * 100).toFixed(0)}%</p>
        <p class="text-sm text-gray-500 mt-1">Based on {totalTrades} trades</p>
      </div>
    </div>
    
    <!-- PnL Chart -->
    <div class="bg-white p-6 rounded-lg shadow-md mb-8">
      <h2 class="text-xl font-semibold mb-4">Profit & Loss (Last 7 Days)</h2>
      <div class="h-80 w-full relative">
        <canvas 
          bind:this={chartCanvas}
          id="pnlChart"
          width="100%"
          height="100%"
          aria-label="Profit and Loss Chart"
        ></canvas>
      </div>
    </div>
    
    <!-- Quick Actions -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-xl font-semibold mb-4">Recent Trades</h2>
        <p class="text-gray-500 mb-4">View and manage your recent trades</p>
        <a href="/trades" class="inline-block bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md transition-colors">
          View All Trades
        </a>
      </div>
      
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-xl font-semibold mb-4">Detailed Analytics</h2>
        <p class="text-gray-500 mb-4">Dive deeper into your trading performance</p>
        <a href="/pnl" class="inline-block bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md transition-colors">
          View Analytics
        </a>
      </div>
    </div>
  {/if}
</div>
