/**
 * 在庫操作の種類
 */
export enum InventoryOperationType {
  IN = 'in',
  OUT = 'out',
  ADJUSTMENT = 'adjustment',
}

/**
 * 在庫操作のステータス
 */
export enum InventoryOperationStatus {
  PENDING = 'pending',
  COMPLETED = 'completed',
  CANCELLED = 'cancelled',
}

/**
 * 在庫操作情報
 */
export interface InventoryOperation {
  id: number;
  productId: number;
  type: InventoryOperationType;
  quantity: number;
  status: InventoryOperationStatus;
  note?: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * 在庫操作作成リクエスト
 */
export interface CreateInventoryOperationRequest {
  productId: number;
  type: InventoryOperationType;
  quantity: number;
  note?: string;
}

/**
 * 在庫操作更新リクエスト
 */
export interface UpdateInventoryOperationRequest extends Partial<CreateInventoryOperationRequest> {
  id: number;
  status: InventoryOperationStatus;
}

/**
 * 在庫一覧のクエリパラメータ
 */
export interface InventoryQueryParams {
  page?: number;
  limit?: number;
  type?: InventoryOperationType;
  status?: InventoryOperationStatus;
  search?: string;
}

/**
 * 在庫一覧のレスポンス
 */
export interface InventoryListResponse {
  items: InventoryOperation[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
} 