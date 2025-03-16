/**
 * 画像圧縮のオプション
 */
interface CompressOptions {
  quality?: number;
  maxSizeKB?: number;
  maxIterations?: number;
}

/**
 * 画像を圧縮する
 */
export async function compressImage(
  file: File,
  options: CompressOptions = { quality: 0.8, maxSizeKB: 500, maxIterations: 5 }
): Promise<File> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    const reader = new FileReader();

    reader.onload = (e) => {
      img.src = e.target?.result as string;
      img.onload = async () => {
        let quality = options.quality || 0.8;
        let iteration = 0;
        let compressedFile = file;

        // 目標サイズに達するまで圧縮を繰り返す
        while (
          compressedFile.size > (options.maxSizeKB || 500) * 1024 &&
          iteration < (options.maxIterations || 5)
        ) {
          const canvas = document.createElement('canvas');
          canvas.width = img.width;
          canvas.height = img.height;

          const ctx = canvas.getContext('2d');
          if (!ctx) {
            reject(new Error('Canvasコンテキストの取得に失敗しました'));
            return;
          }

          // 画像を描画
          ctx.drawImage(img, 0, 0);

          // 品質を下げて圧縮
          quality *= 0.8;

          // Canvasから画像データを取得
          const blob = await new Promise<Blob | null>((resolve) => {
            canvas.toBlob(resolve, file.type, quality);
          });

          if (!blob) {
            reject(new Error('画像の変換に失敗しました'));
            return;
          }

          compressedFile = new File([blob], file.name, {
            type: file.type,
            lastModified: Date.now(),
          });

          iteration++;
        }

        resolve(compressedFile);
      };
    };

    reader.onerror = () => {
      reject(new Error('画像の読み込みに失敗しました'));
    };

    reader.readAsDataURL(file);
  });
}

/**
 * ファイルサイズを人間が読みやすい形式に変換する
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
} 