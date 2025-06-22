import type { Trade, PnLEntry, ApiResponse } from '$lib/types/api';

const API_BASE_URL = import.meta.env.PUBLIC_API_URL || 'http://localhost:5000/api';

export const api = {
  async getTrades(): Promise<Trade[]> {
    const response = await fetch(`${API_BASE_URL}/trades`);
    if (!response.ok) {
      throw new Error('Failed to fetch trades');
    }
    const result: ApiResponse<Trade[]> = await response.json();
    if (result.error) {
      throw new Error(result.error);
    }
    return result.data || [];
  },

  async getPnL(): Promise<PnLEntry[]> {
    const response = await fetch(`${API_BASE_URL}/pnl`);
    if (!response.ok) {
      throw new Error('Failed to fetch PnL');
    }
    const result: ApiResponse<PnLEntry[]> = await response.json();
    if (result.error) {
      throw new Error(result.error);
    }
    return result.data || [];
  },

  async updateTradeNote(tradeId: string, note: string): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/trades/${tradeId}/note`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ note })
    });
    
    if (!response.ok) {
      const result: ApiResponse<null> = await response.json().catch(() => ({}));
      throw new Error(result.error || 'Failed to update trade note');
    }
    
    const result: ApiResponse<null> = await response.json();
    if (result.error) {
      throw new Error(result.error);
    }
  },
};
