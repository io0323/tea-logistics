import { useMutation } from 'react-query';
import { api } from '@/lib/api';
import {
  DataFormat,
  DataType,
  ExportOptions,
  ImportOptions,
  ExportResult,
  ImportResult,
} from '@/types/export';

/**
 * データのエクスポート/インポート機能を提供するカスタムフック
 */
export function useDataExportImport() {
  /**
   * データのエクスポート
   */
  const exportData = useMutation(
    async (options: ExportOptions): Promise<ExportResult> => {
      const response = await api.post('/export', options);
      return response.data;
    }
  );

  /**
   * データのインポート
   */
  const importData = useMutation(
    async ({
      file,
      options,
    }: {
      file: File;
      options: ImportOptions;
    }): Promise<ImportResult> => {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('options', JSON.stringify(options));

      const response = await api.post('/import', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return response.data;
    }
  );

  /**
   * エクスポートしたファイルのダウンロード
   */
  const downloadExportedFile = async (url: string, filename: string) => {
    const response = await fetch(url);
    const blob = await response.blob();
    const downloadUrl = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = downloadUrl;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(downloadUrl);
  };

  return {
    exportData,
    importData,
    downloadExportedFile,
  };
} 