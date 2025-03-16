'use client';

import { useQuery, useMutation, useQueryClient } from 'react-query';
import api from '@/lib/api';
import { Settings, UpdateSettingsRequest } from '@/types/settings';

// モックデータ
const mockSettings: Settings = {
  notification: {
    emailNotifications: false,
    lowStockAlert: false,
    deliveryUpdates: false,
    orderUpdates: false,
  },
  display: {
    theme: 'light',
    language: 'ja',
    timezone: 'Asia/Tokyo',
    dateFormat: 'YYYY/MM/DD',
  },
  system: {
    lowStockThreshold: 10,
    defaultPageSize: 20,
    autoLogoutMinutes: 30,
  },
};

/**
 * 設定管理のためのカスタムフック
 */
export const useSettings = () => {
  const queryClient = useQueryClient();

  // 設定の取得
  const getSettings = () => {
    return useQuery<Settings>(
      'settings',
      async () => {
        // モックデータを使用（実際のAPIが実装されるまで）
        return mockSettings;
      }
    );
  };

  // 設定の更新
  const updateSettings = useMutation<Settings, Error, UpdateSettingsRequest>(
    async (data) => {
      // モックデータを使用（実際のAPIが実装されるまで）
      const updatedSettings = {
        ...mockSettings,
        ...data,
      };
      return updatedSettings;
    },
    {
      onSuccess: () => {
        queryClient.invalidateQueries('settings');
      },
    }
  );

  return {
    getSettings,
    updateSettings,
  };
}; 