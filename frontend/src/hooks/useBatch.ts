import { useQuery, useMutation, useQueryClient } from 'react-query';
import { api } from '@/lib/api';
import {
  BatchType,
  BatchStatus,
  BatchConfig,
  BatchResult,
  BatchQueryParams,
  BatchListResponse,
} from '@/types/batch';

/**
 * バッチ処理のエラー型
 */
export class BatchError extends Error {
  constructor(
    message: string,
    public code: string,
    public details?: any
  ) {
    super(message);
    this.name = 'BatchError';
  }
}

/**
 * バッチ処理を管理するカスタムフック
 */
export function useBatch() {
  const queryClient = useQueryClient();

  /**
   * エラーハンドリング関数
   */
  const handleError = (error: any): BatchError => {
    if (error.response) {
      return new BatchError(
        error.response.data.message || 'バッチ処理でエラーが発生しました',
        error.response.data.code || 'BATCH_ERROR',
        error.response.data.details
      );
    }
    return new BatchError(
      'ネットワークエラーが発生しました',
      'NETWORK_ERROR'
    );
  };

  /**
   * バッチ処理一覧の取得
   */
  const getBatches = (params: BatchQueryParams = {}) => {
    return useQuery<BatchListResponse>(
      ['batches', params],
      async () => {
        const response = await api.get('/batches', { params });
        return response.data;
      },
      {
        staleTime: 30000, // 30秒間はキャッシュを新鮮とみなす
        cacheTime: 300000, // 5分間キャッシュを保持
        refetchOnWindowFocus: false, // ウィンドウフォーカス時の自動リフェッチを無効化
      }
    );
  };

  /**
   * バッチ処理の詳細取得
   */
  const getBatch = (id: number) => {
    return useQuery<BatchResult>(
      ['batch', id],
      async () => {
        const response = await api.get(`/batches/${id}`);
        return response.data;
      },
      {
        enabled: !!id,
        staleTime: 30000,
        cacheTime: 300000,
        refetchInterval: (data) => 
          data?.status === BatchStatus.RUNNING ? 5000 : false, // 実行中のバッチは5秒ごとに更新
      }
    );
  };

  /**
   * バッチ処理の作成
   */
  const createBatch = useMutation(
    async (config: BatchConfig): Promise<BatchResult> => {
      try {
        const response = await api.post('/batches', config);
        return response.data;
      } catch (error) {
        throw handleError(error);
      }
    },
    {
      onSuccess: () => {
        queryClient.invalidateQueries('batches');
      },
    }
  );

  /**
   * バッチ処理の実行
   */
  const executeBatch = useMutation(
    async (id: number): Promise<BatchResult> => {
      try {
        const response = await api.post(`/batches/${id}/execute`);
        return response.data;
      } catch (error) {
        throw handleError(error);
      }
    },
    {
      onSuccess: (_, id) => {
        queryClient.invalidateQueries(['batch', id]);
        queryClient.invalidateQueries('batches');
      },
    }
  );

  /**
   * バッチ処理のキャンセル
   */
  const cancelBatch = useMutation(
    async (id: number): Promise<BatchResult> => {
      try {
        const response = await api.post(`/batches/${id}/cancel`);
        return response.data;
      } catch (error) {
        throw handleError(error);
      }
    },
    {
      onSuccess: (_, id) => {
        queryClient.invalidateQueries(['batch', id]);
        queryClient.invalidateQueries('batches');
      },
    }
  );

  /**
   * バッチ処理の削除
   */
  const deleteBatch = useMutation(
    async (id: number): Promise<void> => {
      try {
        await api.delete(`/batches/${id}`);
      } catch (error) {
        throw handleError(error);
      }
    },
    {
      onSuccess: () => {
        queryClient.invalidateQueries('batches');
      },
    }
  );

  /**
   * バッチ処理のログ取得
   */
  const getBatchLogs = (id: number) => {
    return useQuery<string[]>(
      ['batch-logs', id],
      async () => {
        const response = await api.get(`/batches/${id}/logs`);
        return response.data;
      },
      {
        enabled: !!id,
        staleTime: 5000,
        cacheTime: 300000,
        refetchInterval: (data, query) => {
          const batch = queryClient.getQueryData<BatchResult>(['batch', id]);
          return batch?.status === BatchStatus.RUNNING ? 5000 : false;
        },
      }
    );
  };

  return {
    getBatches,
    getBatch,
    createBatch,
    executeBatch,
    cancelBatch,
    deleteBatch,
    getBatchLogs,
  };
} 