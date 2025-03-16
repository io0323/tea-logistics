/**
 * 在庫履歴の型定義
 */
export interface StockHistory {
  id: string;
  productId: string;
  previousStock: number;
  newStock: number;
  changeAmount: number;
  type: 'in' | 'out' | 'adjustment';
  reason: string;
  createdAt: string;
  createdBy: string;
}

/**
 * 在庫履歴の作成パラメータ
 */
export interface CreateStockHistoryParams {
  productId: string;
  previousStock: number;
  newStock: number;
  type: 'in' | 'out' | 'adjustment';
  reason: string;
} 