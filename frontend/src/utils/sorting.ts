/**
 * ソート方向
 */
export type SortDirection = 'asc' | 'desc';

/**
 * ソート可能なフィールド
 */
export type SortableField = 'name' | 'category' | 'price' | 'createdAt';

/**
 * ソート設定
 */
export interface SortConfig {
  field: SortableField;
  direction: SortDirection;
}

/**
 * ソート設定の初期値
 */
export const defaultSortConfig: SortConfig = {
  field: 'createdAt',
  direction: 'desc',
};

/**
 * ソート方向を切り替える
 */
export function toggleSortDirection(direction: SortDirection): SortDirection {
  return direction === 'asc' ? 'desc' : 'asc';
} 