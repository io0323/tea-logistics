/**
 * 画像リサイズのオプション
 */
interface ResizeOptions {
  maxWidth: number;
  maxHeight: number;
  quality?: number;
}

/**
 * 画像をリサイズする
 */
export async function resizeImage(
  file: File,
  options: ResizeOptions = { maxWidth: 800, maxHeight: 800, quality: 0.8 }
): Promise<File> {
  return new Promise((resolve, reject) => {
    const img = new Image();
    const reader = new FileReader();

    reader.onload = (e) => {
      img.src = e.target?.result as string;
      img.onload = () => {
        const canvas = document.createElement('canvas');
        let { width, height } = img;

        // アスペクト比を保持しながらリサイズ
        if (width > options.maxWidth) {
          height = (height * options.maxWidth) / width;
          width = options.maxWidth;
        }
        if (height > options.maxHeight) {
          width = (width * options.maxHeight) / height;
          height = options.maxHeight;
        }

        canvas.width = width;
        canvas.height = height;

        const ctx = canvas.getContext('2d');
        if (!ctx) {
          reject(new Error('Canvasコンテキストの取得に失敗しました'));
          return;
        }

        // 画像を描画
        ctx.drawImage(img, 0, 0, width, height);

        // Canvasから画像データを取得
        canvas.toBlob(
          (blob) => {
            if (!blob) {
              reject(new Error('画像の変換に失敗しました'));
              return;
            }
            // 新しいFileオブジェクトを作成
            const resizedFile = new File([blob], file.name, {
              type: file.type,
              lastModified: Date.now(),
            });
            resolve(resizedFile);
          },
          file.type,
          options.quality
        );
      };
    };

    reader.onerror = () => {
      reject(new Error('画像の読み込みに失敗しました'));
    };

    reader.readAsDataURL(file);
  });
} 