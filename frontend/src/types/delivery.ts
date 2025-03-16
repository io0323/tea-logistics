/**
 * 配送ステータス
 */
export enum DeliveryStatus {
  PENDING = 'pending',
  IN_TRANSIT = 'in_transit',
  DELIVERED = 'delivered',
  CANCELLED = 'cancelled',
}

/**
 * 配送情報
 */
export interface Delivery {
  id: number;
  orderId: number;
  customerName: string;
  customerAddress: string;
  customerPhone: string;
  status: DeliveryStatus;
  estimatedDeliveryDate?: string;
  actualDeliveryDate?: string;
  note?: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * 配送作成リクエスト
 */
export interface CreateDeliveryRequest {
  orderId: number;
  customerName: string;
  customerAddress: string;
  customerPhone: string;
  estimatedDeliveryDate?: string;
  note?: string;
}

/**
 * 配送更新リクエスト
 */
export interface UpdateDeliveryRequest extends Partial<CreateDeliveryRequest> {
  id: number;
  status: DeliveryStatus;
  actualDeliveryDate?: string;
}

/**
 * 配送一覧のクエリパラメータ
 */
export interface DeliveryQueryParams {
  page?: number;
  limit?: number;
  status?: DeliveryStatus;
  search?: string;
  startDate?: string;
  endDate?: string;
}

/**
 * 配送一覧のレスポンス
 */
export interface DeliveryListResponse {
  items: Delivery[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
} 