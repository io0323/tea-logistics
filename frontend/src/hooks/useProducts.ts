import { useQuery, useMutation, useQueryClient } from 'react-query';
import api from '@/lib/api';
import {
  Product,
  CreateProductRequest,
  UpdateProductRequest,
  ProductQueryParams,
  ProductListResponse,
} from '@/types/product';

/**
 * 商品管理のためのカスタムフック
 */
export const useProducts = () => {
  const queryClient = useQueryClient();

  // 商品一覧の取得
  const getProducts = (params: ProductQueryParams = {}) => {
    return useQuery<ProductListResponse>(
      ['products', params],
      async () => {
        const response = await api.get('/products', { params });
        return response.data;
      }
    );
  };

  // 商品詳細の取得
  const getProduct = (id: number) => {
    return useQuery<Product>(
      ['product', id],
      async () => {
        const response = await api.get(`/products/${id}`);
        return response.data;
      },
      {
        enabled: !!id,
      }
    );
  };

  // 商品の作成
  const createProduct = useMutation<Product, Error, CreateProductRequest>(
    async (data) => {
      const response = await api.post('/products', data);
      return response.data;
    },
    {
      onSuccess: () => {
        queryClient.invalidateQueries('products');
      },
    }
  );

  // 商品の更新
  const updateProduct = useMutation<Product, Error, UpdateProductRequest>(
    async (data) => {
      const response = await api.put(`/products/${data.id}`, data);
      return response.data;
    },
    {
      onSuccess: (data) => {
        queryClient.invalidateQueries('products');
        queryClient.invalidateQueries(['product', data.id]);
      },
    }
  );

  // 商品の削除
  const deleteProduct = useMutation<void, Error, number>(
    async (id) => {
      await api.delete(`/products/${id}`);
    },
    {
      onSuccess: (_, id) => {
        queryClient.invalidateQueries('products');
        queryClient.invalidateQueries(['product', id]);
      },
    }
  );

  return {
    getProducts,
    getProduct,
    createProduct,
    updateProduct,
    deleteProduct,
  };
}; 