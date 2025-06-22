export interface Trade {
  id: string;
  symbol: string;
  side: string;
  price: number;
  quantity: number;
  quoteQty: number;
  commission: number;
  commissionAsset: string;
  time: string;
  isBuyer: boolean;
  isMaker: boolean;
  isBestMatch: boolean;
  note?: string;
}

export interface PnLEntry {
  date: string;
  pnl: number;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
}
