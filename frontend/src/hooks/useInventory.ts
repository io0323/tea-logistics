'use client';

import { useState, useEffect } from 'react';
import api from '@/lib/api';

interface Inventory {
  id: string;
  code: string;
  name: string;
  category: string;
  stock: number;
  unit: string;
}

interface UseInventoryParams {
  page: number;
  pageSize: number;
  searchQuery?: string;
}

interface UseInventoryReturn {
  inventories: Inventory[];
  loading: boolean;
  error: Error | null;
  totalPages: number;
  createInventory: (data: Omit<Inventory, 'id'>) => Promise<void>;
  updateInventory: (id: string, data: Partial<Inventory>) => Promise<void>;
  deleteInventory: (id: string) => Promise<void>;
}

/**
 * 在庫管理用のカスタムフック
 */
export function useInventory({ page, pageSize, searchQuery }: UseInventoryParams): UseInventoryReturn {
  const [inventories, setInventories] = useState<Inventory[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [totalPages, setTotalPages] = useState(1);

  useEffect(() => {
    const fetchInventories = async () => {
      try {
        setLoading(true);
        setError(null);
        
        // モックデータを使用（実際のAPIが実装されるまで）
        const mockData = {
          items: [
            {
              id: '1',
              code: '001',
              name: '煎茶 - 特上',
              category: '緑茶',
              stock: 100,
              unit: 'kg'
            },
            {
              id: '2',
              code: '002',
              name: '玉露',
              category: '緑茶',
              stock: 50,
              unit: 'kg'
            }
          ],
          total: 2,
          totalPages: 1
        };

        setInventories(mockData.items);
        setTotalPages(mockData.totalPages);
      } catch (err) {
        setError(err instanceof Error ? err : new Error('在庫情報の取得に失敗しました'));
      } finally {
        setLoading(false);
      }
    };

    fetchInventories();
  }, [page, pageSize, searchQuery]);

  const createInventory = async (data: Omit<Inventory, 'id'>) => {
    try {
      setLoading(true);
      // モック実装（実際のAPIが実装されるまで）
      const mockResponse = {
        id: Math.random().toString(),
        ...data
      };
      setInventories(prev => [...prev, mockResponse]);
    } catch (err) {
      throw err instanceof Error ? err : new Error('在庫情報の作成に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const updateInventory = async (id: string, data: Partial<Inventory>) => {
    try {
      setLoading(true);
      // モック実装（実際のAPIが実装されるまで）
      setInventories(prev =>
        prev.map(item => (item.id === id ? { ...item, ...data } : item))
      );
    } catch (err) {
      throw err instanceof Error ? err : new Error('在庫情報の更新に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const deleteInventory = async (id: string) => {
    try {
      setLoading(true);
      // モック実装（実際のAPIが実装されるまで）
      setInventories(prev => prev.filter(item => item.id !== id));
    } catch (err) {
      throw err instanceof Error ? err : new Error('在庫情報の削除に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  return {
    inventories,
    loading,
    error,
    totalPages,
    createInventory,
    updateInventory,
    deleteInventory,
  };
} 