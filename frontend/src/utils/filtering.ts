/**
 * 価格範囲
 */
export interface PriceRange {
  min: number;
  max: number;
}

/**
 * フィルター設定
 */
export interface FilterConfig {
  category: string;
  priceRange: PriceRange;
  searchQuery: string;
}

/**
 * フィルター設定の初期値
 */
export const defaultFilterConfig: FilterConfig = {
  category: '',
  priceRange: {
    min: 0,
    max: 1000000,
  },
  searchQuery: '',
};

/**
 * フィルター設定をクリアする
 */
export function clearFilterConfig(): FilterConfig {
  return defaultFilterConfig;
}

/**
 * フィルター設定が有効かどうかを判定する
 */
export function hasActiveFilters(config: FilterConfig): boolean {
  return (
    config.category !== '' ||
    config.searchQuery !== '' ||
    config.priceRange.min > 0 ||
    config.priceRange.max < 1000000
  );
} 