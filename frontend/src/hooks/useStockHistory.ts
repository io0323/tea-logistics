'use client';

import { useState, useCallback } from 'react';
import api from '@/lib/api';
import { StockHistory, CreateStockHistoryParams } from '@/types/stockHistory';

interface UseStockHistoryParams {
  productId: string;
  page: number;
  pageSize: number;
}

interface UseStockHistoryReturn {
  histories: StockHistory[];
  totalPages: number;
  isLoading: boolean;
  error: Error | null;
  fetchHistories: () => Promise<void>;
  createHistory: (data: CreateStockHistoryParams) => Promise<void>;
}

/**
 * 在庫履歴管理用カスタムフック
 */
export function useStockHistory(params: UseStockHistoryParams): UseStockHistoryReturn {
  const [histories, setHistories] = useState<StockHistory[]>([]);
  const [totalPages, setTotalPages] = useState(1);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchHistories = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await api.get(`/products/${params.productId}/stock-history`, {
        params: {
          page: params.page,
          pageSize: params.pageSize,
        },
      });
      setHistories(response.data.histories);
      setTotalPages(response.data.totalPages);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('在庫履歴の取得に失敗しました'));
    } finally {
      setIsLoading(false);
    }
  }, [params]);

  const createHistory = useCallback(async (data: CreateStockHistoryParams) => {
    try {
      setIsLoading(true);
      setError(null);
      await api.post(`/products/${params.productId}/stock-history`, data);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('在庫履歴の作成に失敗しました'));
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, [params.productId]);

  return {
    histories,
    totalPages,
    isLoading,
    error,
    fetchHistories,
    createHistory,
  };
} 