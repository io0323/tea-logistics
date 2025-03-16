import { useQuery, useMutation, useQueryClient } from 'react-query';
import api from '@/lib/api';
import {
  Delivery,
  CreateDeliveryRequest,
  UpdateDeliveryRequest,
  DeliveryQueryParams,
  DeliveryListResponse,
} from '@/types/delivery';

/**
 * 配送管理のためのカスタムフック
 */
export const useDelivery = () => {
  const queryClient = useQueryClient();

  // 配送一覧の取得
  const getDeliveries = (params: DeliveryQueryParams = {}) => {
    return useQuery<DeliveryListResponse>(
      ['deliveries', params],
      async () => {
        const response = await api.get('/deliveries', { params });
        return response.data;
      }
    );
  };

  // 配送詳細の取得
  const getDelivery = (id: number) => {
    return useQuery<Delivery>(
      ['delivery', id],
      async () => {
        const response = await api.get(`/deliveries/${id}`);
        return response.data;
      },
      {
        enabled: !!id,
      }
    );
  };

  // 配送の作成
  const createDelivery = useMutation<Delivery, Error, CreateDeliveryRequest>(
    async (data) => {
      const response = await api.post('/deliveries', data);
      return response.data;
    },
    {
      onSuccess: () => {
        queryClient.invalidateQueries('deliveries');
      },
    }
  );

  // 配送の更新
  const updateDelivery = useMutation<Delivery, Error, UpdateDeliveryRequest>(
    async (data) => {
      const response = await api.put(`/deliveries/${data.id}`, data);
      return response.data;
    },
    {
      onSuccess: (data) => {
        queryClient.invalidateQueries('deliveries');
        queryClient.invalidateQueries(['delivery', data.id]);
      },
    }
  );

  // 配送の削除
  const deleteDelivery = useMutation<void, Error, number>(
    async (id) => {
      await api.delete(`/deliveries/${id}`);
    },
    {
      onSuccess: (_, id) => {
        queryClient.invalidateQueries('deliveries');
        queryClient.invalidateQueries(['delivery', id]);
      },
    }
  );

  return {
    getDeliveries,
    getDelivery,
    createDelivery,
    updateDelivery,
    deleteDelivery,
  };
}; 