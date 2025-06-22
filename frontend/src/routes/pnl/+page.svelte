<script lang="ts">
  import { onMount } from 'svelte';
  import { Chart } from 'chart.js/auto';
  import type { PnLEntry } from '$lib/types/api';
  import { api } from '$lib/services/api';

  let loading = true;
  let error: string | null = null;
  let chart: Chart | null = null;

  onMount(() => {
    (async () => {
      try {
        const data = await api.getPnL();
        renderChart(data);
      } catch (e) {
        error = e instanceof Error ? e.message : 'An unknown error occurred';
      } finally {
        loading = false;
      }
    })();

    return () => {
      if (chart) {
        chart.destroy();
      }
    };
  });

  function renderChart(pnlData: PnLEntry[]) {
    const ctx = document.getElementById('pnlChart') as HTMLCanvasElement;
    
    if (chart) {
      chart.destroy();
    }

    chart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: pnlData.map(item => item.date),
        datasets: [{
          label: 'Daily PnL (USDT)',
          data: pnlData.map(item => item.pnl),
          borderColor: 'rgb(59, 130, 246)',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          tension: 0.3,
          fill: true
        }]
      },
      options: {
        responsive: true,
        plugins: {
          title: {
            display: true,
            text: 'Profit & Loss (Last 7 Days)'
          },
          tooltip: {
            mode: 'index',
            intersect: false
          }
        },
        scales: {
          y: {
            beginAtZero: true,
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
    });
  }
</script>

<div class="container mx-auto p-4">
  <h1 class="text-2xl font-bold mb-6">Profit & Loss</h1>
  
  {#if loading}
    <div class="text-center py-8">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500 mx-auto"></div>
      <p class="mt-2">Loading PnL data...</p>
    </div>
  {:else if error}
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">Error:</strong>
      <span class="block sm:inline"> {error}</span>
    </div>
  {:else}
    <div class="bg-white p-6 rounded-lg shadow-md">
      <div class="h-96">
        <canvas id="pnlChart"></canvas>
      </div>
    </div>
  {/if}
</div>
