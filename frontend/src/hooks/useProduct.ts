'use client';

import { useState, useCallback } from 'react';
import api from '@/lib/api';
import { convertToCSV, downloadFile } from '@/utils/productExport';
import { Product, CreateProductParams, UpdateProductParams } from '@/types/product';

interface UseProductParams {
  page: number;
  pageSize: number;
  sortField?: string;
  sortDirection?: 'asc' | 'desc';
  category?: string;
  searchQuery?: string;
  priceRange?: {
    min: number;
    max: number;
  };
}

interface UseProductReturn {
  products: Product[];
  product: Product | null;
  totalPages: number;
  isLoading: boolean;
  error: Error | null;
  categories: string[];
  fetchProducts: () => Promise<void>;
  fetchProduct: (id: string) => Promise<void>;
  createProduct: (params: CreateProductParams) => Promise<void>;
  updateProduct: (id: string, params: UpdateProductParams) => Promise<void>;
  deleteProduct: (id: string) => Promise<void>;
  bulkDeleteProducts: (ids: string[]) => Promise<void>;
  exportProducts: () => Promise<void>;
  importProducts: (products: Partial<Product>[]) => Promise<void>;
  bulkUpdateProducts: (ids: string[], params: Partial<UpdateProductParams>) => Promise<void>;
  uploadProductImage: (file: File) => Promise<string>;
}

/**
 * 商品管理用カスタムフック
 */
export function useProduct(params: UseProductParams = {}): UseProductReturn {
  const [products, setProducts] = useState<Product[]>([]);
  const [product, setProduct] = useState<Product | null>(null);
  const [totalPages, setTotalPages] = useState(1);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const [categories, setCategories] = useState<string[]>([]);

  const fetchProducts = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await api.get('/products', { params });
      setProducts(response.data.products);
      setTotalPages(response.data.totalPages);
      setCategories(response.data.categories);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('商品の取得に失敗しました'));
    } finally {
      setIsLoading(false);
    }
  }, [params]);

  const fetchProduct = useCallback(async (id: string) => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await api.get(`/products/${id}`);
      setProduct(response.data);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('商品の取得に失敗しました'));
    } finally {
      setIsLoading(false);
    }
  }, []);

  const createProduct = useCallback(async (params: CreateProductParams) => {
    try {
      setIsLoading(true);
      setError(null);
      const formData = new FormData();
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          formData.append(key, value);
        }
      });
      const response = await api.post<Product>('/products', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      setProducts((prev) => [...prev, response.data]);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('商品の作成に失敗しました'));
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const updateProduct = useCallback(async (id: string, params: UpdateProductParams) => {
    try {
      setIsLoading(true);
      setError(null);
      const formData = new FormData();
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined) {
          formData.append(key, value);
        }
      });
      const response = await api.put<Product>(`/products/${id}`, formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      setProducts((prev) =>
        prev.map((product) => (product.id === id ? response.data : product))
      );
    } catch (err) {
      setError(err instanceof Error ? err : new Error('商品の更新に失敗しました'));
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const deleteProduct = useCallback(async (id: string) => {
    try {
      setIsLoading(true);
      setError(null);
      await api.delete(`/products/${id}`);
      setProducts((prev) => prev.filter((product) => product.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err : new Error('商品の削除に失敗しました'));
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const bulkDeleteProducts = useCallback(async (ids: string[]) => {
    try {
      setIsLoading(true);
      setError(null);
      await api.post('/products/bulk-delete', { ids });
      setProducts((prev) => prev.filter((product) => !ids.includes(product.id)));
    } catch (err) {
      setError(err instanceof Error ? err : new Error('商品の一括削除に失敗しました'));
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const exportProducts = useCallback(async () => {
    try {
      setIsLoading(true);
      setError(null);
      const response = await api.get('/products/export');
      const csvContent = convertToCSV(response.data);
      downloadFile(csvContent, `products_${new Date().toISOString().split('T')[0]}.csv`);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('商品のエクスポートに失敗しました'));
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const importProducts = useCallback(async (products: Partial<Product>[]) => {
    try {
      setIsLoading(true);
      setError(null);
      await api.post('/products/import', { products });
    } catch (err) {
      setError(err instanceof Error ? err : new Error('商品のインポートに失敗しました'));
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  const bulkUpdateProducts = useCallback(
    async (ids: string[], params: Partial<UpdateProductParams>) => {
      try {
        setIsLoading(true);
        setError(null);
        const formData = new FormData();
        formData.append('ids', JSON.stringify(ids));
        Object.entries(params).forEach(([key, value]) => {
          if (value !== undefined) {
            formData.append(key, value);
          }
        });
        const response = await api.put<Product[]>('/products/bulk-update', formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        });
        setProducts((prev) =>
          prev.map((product) =>
            ids.includes(product.id)
              ? response.data.find((p) => p.id === product.id) || product
              : product
          )
        );
      } catch (err) {
        setError(err instanceof Error ? err : new Error('商品の一括更新に失敗しました'));
        throw err;
      } finally {
        setIsLoading(false);
      }
    },
    []
  );

  const uploadProductImage = useCallback(async (file: File): Promise<string> => {
    try {
      setIsLoading(true);
      setError(null);
      const formData = new FormData();
      formData.append('image', file);
      const response = await api.post<{ url: string }>('/products/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return response.data.url;
    } catch (err) {
      setError(err instanceof Error ? err : new Error('画像のアップロードに失敗しました'));
      throw err;
    } finally {
      setIsLoading(false);
    }
  }, []);

  return {
    products,
    product,
    totalPages,
    isLoading,
    error,
    categories,
    fetchProducts,
    fetchProduct,
    createProduct,
    updateProduct,
    deleteProduct,
    bulkDeleteProducts,
    exportProducts,
    importProducts,
    bulkUpdateProducts,
    uploadProductImage,
  };
} 