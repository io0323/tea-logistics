/**
 * エクスポート/インポートのデータ形式
 */
export enum DataFormat {
  CSV = 'csv',
  JSON = 'json',
  EXCEL = 'excel',
}

/**
 * エクスポート/インポートのデータタイプ
 */
export enum DataType {
  INVENTORY = 'inventory',
  DELIVERY = 'delivery',
  PRODUCT = 'product',
}

/**
 * エクスポート設定
 */
export interface ExportOptions {
  format: DataFormat;
  type: DataType;
  startDate?: string;
  endDate?: string;
  includeHeaders?: boolean;
}

/**
 * インポート設定
 */
export interface ImportOptions {
  format: DataFormat;
  type: DataType;
  skipHeaders?: boolean;
  validateData?: boolean;
}

/**
 * エクスポート結果
 */
export interface ExportResult {
  url: string;
  filename: string;
  format: DataFormat;
  type: DataType;
  totalRecords: number;
  createdAt: string;
}

/**
 * インポート結果
 */
export interface ImportResult {
  totalRecords: number;
  successCount: number;
  errorCount: number;
  errors?: Array<{
    row: number;
    message: string;
  }>;
  createdAt: string;
} 