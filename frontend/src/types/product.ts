/**
 * 商品のカテゴリー
 */
export enum ProductCategory {
  GREEN_TEA = 'green_tea',
  BLACK_TEA = 'black_tea',
  OOLONG_TEA = 'oolong_tea',
  PUERH_TEA = 'puerh_tea',
  HERBAL_TEA = 'herbal_tea',
  OTHER = 'other',
}

/**
 * 商品のステータス
 */
export enum ProductStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  DISCONTINUED = 'discontinued',
}

/**
 * 商品情報
 */
export interface Product {
  id: string;
  name: string;
  category: string;
  price: number;
  stock: number;
  description?: string;
  imageUrl?: string;
  imageFile?: File;
  createdAt: string;
  updatedAt: string;
}

/**
 * 商品作成リクエスト
 */
export interface CreateProductRequest {
  name: string;
  description: string;
  category: ProductCategory;
  price: number;
  stock: number;
  status: ProductStatus;
  imageUrl?: string;
}

/**
 * 商品更新リクエスト
 */
export interface UpdateProductRequest extends Partial<CreateProductRequest> {
  id: number;
}

/**
 * 商品一覧のクエリパラメータ
 */
export interface ProductQueryParams {
  page?: number;
  limit?: number;
  category?: ProductCategory;
  status?: ProductStatus;
  search?: string;
}

/**
 * 商品一覧のレスポンス
 */
export interface ProductListResponse {
  items: Product[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

/**
 * 商品の作成パラメータ
 */
export interface CreateProductParams {
  name: string;
  category: string;
  price: number;
  stock: number;
  description?: string;
  imageFile?: File;
}

/**
 * 商品の更新パラメータ
 */
export interface UpdateProductParams {
  name?: string;
  category?: string;
  price?: number;
  stock?: number;
  description?: string;
  imageFile?: File;
} 