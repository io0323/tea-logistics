import { Product } from '@/types/product';

/**
 * 商品データをCSV形式に変換する
 */
export function convertToCSV(products: Product[]): string {
  const headers = [
    '商品名',
    'カテゴリー',
    '価格',
    '在庫数',
    '説明',
  ];

  const rows = products.map((product) => [
    product.name,
    product.category,
    product.price.toString(),
    product.stock.toString(),
    product.description || '',
  ]);

  const csvContent = [
    headers.join(','),
    ...rows.map((row) => row.join(',')),
  ].join('\n');

  return csvContent;
}

/**
 * CSVデータを商品データに変換する
 */
export function convertFromCSV(csvContent: string): Partial<Product>[] {
  const lines = csvContent.split('\n');
  const headers = lines[0].split(',');
  const products: Partial<Product>[] = [];

  for (let i = 1; i < lines.length; i++) {
    const values = lines[i].split(',');
    const product: Partial<Product> = {};

    headers.forEach((header, index) => {
      switch (header) {
        case '商品名':
          product.name = values[index];
          break;
        case 'カテゴリー':
          product.category = values[index];
          break;
        case '価格':
          product.price = parseInt(values[index], 10);
          break;
        case '在庫数':
          product.stock = parseInt(values[index], 10);
          break;
        case '説明':
          product.description = values[index];
          break;
      }
    });

    products.push(product);
  }

  return products;
}

/**
 * ファイルをダウンロードする
 */
export function downloadFile(content: string, filename: string) {
  const blob = new Blob([content], { type: 'text/csv;charset=utf-8;' });
  const link = document.createElement('a');
  const url = URL.createObjectURL(blob);
  link.setAttribute('href', url);
  link.setAttribute('download', filename);
  link.style.visibility = 'hidden';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

/**
 * ファイルを読み込む
 */
export function readFile(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      if (e.target?.result) {
        resolve(e.target.result as string);
      } else {
        reject(new Error('ファイルの読み込みに失敗しました'));
      }
    };
    reader.onerror = () => reject(new Error('ファイルの読み込みに失敗しました'));
    reader.readAsText(file);
  });
} 