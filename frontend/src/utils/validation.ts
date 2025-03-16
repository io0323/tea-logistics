/**
 * 商品フォームのバリデーションエラー
 */
export interface ProductValidationError {
  field: string;
  message: string;
}

/**
 * 商品データのバリデーション
 */
export function validateProduct(data: {
  name: string;
  description?: string;
  category: string;
  unit: string;
  price: number;
}): ProductValidationError[] {
  const errors: ProductValidationError[] = [];

  // 商品名のバリデーション
  if (!data.name.trim()) {
    errors.push({
      field: 'name',
      message: '商品名は必須です',
    });
  } else if (data.name.length > 100) {
    errors.push({
      field: 'name',
      message: '商品名は100文字以内で入力してください',
    });
  }

  // 説明のバリデーション
  if (data.description && data.description.length > 1000) {
    errors.push({
      field: 'description',
      message: '説明は1000文字以内で入力してください',
    });
  }

  // カテゴリーのバリデーション
  if (!data.category) {
    errors.push({
      field: 'category',
      message: 'カテゴリーは必須です',
    });
  }

  // 単位のバリデーション
  if (!data.unit.trim()) {
    errors.push({
      field: 'unit',
      message: '単位は必須です',
    });
  } else if (data.unit.length > 10) {
    errors.push({
      field: 'unit',
      message: '単位は10文字以内で入力してください',
    });
  }

  // 価格のバリデーション
  if (data.price <= 0) {
    errors.push({
      field: 'price',
      message: '価格は0より大きい値を入力してください',
    });
  } else if (data.price > 1000000) {
    errors.push({
      field: 'price',
      message: '価格は1,000,000円以内で入力してください',
    });
  }

  return errors;
} 