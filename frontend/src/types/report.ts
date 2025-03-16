/**
 * レポートの期間タイプ
 */
export enum ReportPeriodType {
  DAILY = 'daily',
  WEEKLY = 'weekly',
  MONTHLY = 'monthly',
  YEARLY = 'yearly',
}

/**
 * 売上レポートのデータ
 */
export interface SalesReportData {
  period: string;
  totalSales: number;
}

/**
 * 在庫レポートのデータ
 */
export interface InventoryReportData {
  period: string;
  stockQuantity: number;
  lowStockItems: number;
}

/**
 * 配送レポートのデータ
 */
export interface DeliveryReportData {
  period: string;
  onTimeDeliveryRate: number;
}

/**
 * レポートのクエリパラメータ
 */
export interface ReportQueryParams {
  periodType: string;
  startDate: string;
  endDate: string;
}

/**
 * レポートのレスポンス
 */
export interface ReportResponse {
  salesReport: SalesReportData[];
  inventoryReport: InventoryReportData[];
  deliveryReport: DeliveryReportData[];
}

export interface SalesReportItem {
  period: string;
  totalSales: number;
}

export interface DeliveryReportItem {
  period: string;
  onTimeDeliveryRate: number;
} 