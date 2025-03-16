'use client';

import { useQuery } from 'react-query';
import api from '@/lib/api';
import {
  ReportQueryParams,
  ReportResponse,
  SalesReportData,
  InventoryReportData,
  DeliveryReportData,
} from '@/types/report';

/**
 * レポート取得のためのカスタムフック
 */
export function useReport() {
  // レポートデータの取得
  const getReport = (params: ReportQueryParams) => {
    return useQuery<ReportResponse>(
      ['report', params],
      async () => {
        // モックデータ
        const mockData: ReportResponse = {
          salesReport: [
            { period: '1月', totalSales: 1200000 },
            { period: '2月', totalSales: 1500000 },
            { period: '3月', totalSales: 1300000 },
            { period: '4月', totalSales: 1600000 },
            { period: '5月', totalSales: 1400000 },
            { period: '6月', totalSales: 1700000 },
          ],
          inventoryReport: [
            { period: '1月', totalInventory: 5000 },
            { period: '2月', totalInventory: 4800 },
            { period: '3月', totalInventory: 5200 },
            { period: '4月', totalInventory: 5100 },
            { period: '5月', totalInventory: 5300 },
            { period: '6月', totalInventory: 5400 },
          ],
          deliveryReport: [
            { period: '1月', onTimeDeliveryRate: 95 },
            { period: '2月', onTimeDeliveryRate: 92 },
            { period: '3月', onTimeDeliveryRate: 88 },
            { period: '4月', onTimeDeliveryRate: 94 },
            { period: '5月', onTimeDeliveryRate: 96 },
            { period: '6月', onTimeDeliveryRate: 93 },
          ],
        };

        // APIが実装されたら、以下のように変更
        // const response = await api.get('/reports', { params });
        // return response.data;

        // モックデータを使用
        return mockData;
      }
    );
  };

  // 売上レポートの取得
  const getSalesReport = (params: ReportQueryParams) => {
    return useQuery<SalesReportData[]>(
      ['salesReport', params],
      async () => {
        // モックデータ
        return [
          { period: '1月', totalSales: 1200000 },
          { period: '2月', totalSales: 1500000 },
          { period: '3月', totalSales: 1300000 },
          { period: '4月', totalSales: 1600000 },
          { period: '5月', totalSales: 1400000 },
          { period: '6月', totalSales: 1700000 },
        ];
      }
    );
  };

  // 在庫レポートの取得
  const getInventoryReport = (params: ReportQueryParams) => {
    return useQuery<InventoryReportData[]>(
      ['inventoryReport', params],
      async () => {
        // モックデータ
        return [
          { period: '1月', totalInventory: 5000 },
          { period: '2月', totalInventory: 4800 },
          { period: '3月', totalInventory: 5200 },
          { period: '4月', totalInventory: 5100 },
          { period: '5月', totalInventory: 5300 },
          { period: '6月', totalInventory: 5400 },
        ];
      }
    );
  };

  // 配送レポートの取得
  const getDeliveryReport = (params: ReportQueryParams) => {
    return useQuery<DeliveryReportData[]>(
      ['deliveryReport', params],
      async () => {
        // モックデータ
        return [
          { period: '1月', onTimeDeliveryRate: 95 },
          { period: '2月', onTimeDeliveryRate: 92 },
          { period: '3月', onTimeDeliveryRate: 88 },
          { period: '4月', onTimeDeliveryRate: 94 },
          { period: '5月', onTimeDeliveryRate: 96 },
          { period: '6月', onTimeDeliveryRate: 93 },
        ];
      }
    );
  };

  return {
    getReport,
    getSalesReport,
    getInventoryReport,
    getDeliveryReport,
  };
} 