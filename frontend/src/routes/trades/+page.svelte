<script lang="ts">
  import { onMount } from 'svelte';
  import type { Trade } from '$lib/types/api';
  import { api } from '$lib/services/api';

  let trades: Trade[] = [];
  let loading = true;
  let error: string | null = null;

  onMount(async () => {
    try {
      trades = await api.getTrades();
    } catch (e) {
      error = e instanceof Error ? e.message : 'Failed to load trades';
    } finally {
      loading = false;
    }
  });

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const updateNote = async (tradeId: string, event: Event) => {
    const target = event.target as HTMLInputElement;
    const note = target.value;
    
    try {
      await api.updateTradeNote(tradeId, note);
    } catch (e) {
      console.error('Failed to update note:', e);
    }
  };
</script>

<div class="container mx-auto p-4">
  <h1 class="text-2xl font-bold mb-6">Trade History</h1>
  
  {#if loading}
    <div class="text-center py-8">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500 mx-auto"></div>
      <p class="mt-2">Loading trades...</p>
    </div>
  {:else if error}
    <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">Error:</strong>
      <span class="block sm:inline"> {error}</span>
    </div>
  {:else}
    <div class="overflow-x-auto">
      <table class="min-w-full bg-white">
        <thead>
          <tr class="bg-gray-100">
            <th class="py-2 px-4 text-left">Time</th>
            <th class="py-2 px-4 text-left">Symbol</th>
            <th class="py-2 px-4 text-right">Side</th>
            <th class="py-2 px-4 text-right">Price</th>
            <th class="py-2 px-4 text-right">Quantity</th>
            <th class="py-2 px-4 text-right">Total</th>
            <th class="py-2 px-4 text-left">Note</th>
          </tr>
        </thead>
        <tbody>
          {#each trades as trade (trade.id)}
            <tr class="border-t">
              <td class="py-2 px-4">{formatDate(trade.time)}</td>
              <td class="py-2 px-4 font-mono">{trade.symbol}</td>
              <td class="py-2 px-4 text-right">
                <span class={`px-2 py-1 rounded ${trade.isBuyer ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
                  {trade.isBuyer ? 'BUY' : 'SELL'}
                </span>
              </td>
              <td class="py-2 px-4 text-right font-mono">{trade.price.toFixed(2)}</td>
              <td class="py-2 px-4 text-right font-mono">{trade.quantity.toFixed(4)}</td>
              <td class="py-2 px-4 text-right font-mono">{trade.quoteQty.toFixed(2)}</td>
              <td class="py-2 px-4">
                <input
                  type="text"
                  class="border rounded px-2 py-1 w-full"
                  value={trade.note || ''}
                  on:change={(e) => updateNote(trade.id, e)}
                  placeholder="Add a note..."
                />
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</div>
