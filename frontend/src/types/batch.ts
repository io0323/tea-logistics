/**
 * バッチ処理のタイプ
 */
export enum BatchType {
  STOCK_CHECK = 'stock_check',          // 在庫確認
  DELIVERY_STATUS_UPDATE = 'delivery_status_update',  // 配送ステータス更新
  DATA_CLEANUP = 'data_cleanup',        // データクリーンアップ
  REPORT_GENERATION = 'report_generation', // レポート生成
}

/**
 * バッチ処理のステータス
 */
export enum BatchStatus {
  PENDING = 'pending',     // 待機中
  RUNNING = 'running',     // 実行中
  COMPLETED = 'completed', // 完了
  FAILED = 'failed',       // 失敗
  CANCELLED = 'cancelled', // キャンセル
}

/**
 * バッチ処理の設定
 */
export interface BatchConfig {
  type: BatchType;
  schedule?: string;      // cronスケジュール
  retryCount?: number;    // リトライ回数
  timeout?: number;       // タイムアウト（秒）
  params?: Record<string, any>; // 追加パラメータ
}

/**
 * バッチ処理の実行結果
 */
export interface BatchResult {
  id: number;
  type: BatchType;
  status: BatchStatus;
  startTime: string;
  endTime?: string;
  duration?: number;      // 実行時間（秒）
  processedItems: number; // 処理件数
  successCount: number;   // 成功件数
  errorCount: number;     // エラー件数
  errors?: Array<{
    message: string;
    details?: any;
  }>;
  logs?: string[];       // 実行ログ
}

/**
 * バッチ処理の検索条件
 */
export interface BatchQueryParams {
  type?: BatchType;
  status?: BatchStatus;
  startDate?: string;
  endDate?: string;
  page?: number;
  limit?: number;
}

/**
 * バッチ処理一覧のレスポンス
 */
export interface BatchListResponse {
  items: BatchResult[];
  total: number;
  totalPages: number;
  currentPage: number;
} 