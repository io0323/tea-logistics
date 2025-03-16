import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  FormControl,
  FormLabel,
  Select,
  Input,
  VStack,
  useToast,
  Text,
  Switch,
  FormHelperText,
  Box,
  Progress,
} from '@chakra-ui/react';
import { useState, useRef } from 'react';
import { useDataExportImport } from '@/hooks/useDataExportImport';
import {
  DataFormat,
  DataType,
  ExportOptions,
  ImportOptions,
} from '@/types/export';

interface ExportImportModalProps {
  isOpen: boolean;
  onClose: () => void;
  type: DataType;
  mode: 'export' | 'import';
}

/**
 * エクスポート/インポートモーダル
 */
export default function ExportImportModal({
  isOpen,
  onClose,
  type,
  mode,
}: ExportImportModalProps) {
  const toast = useToast();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { exportData, importData, downloadExportedFile } = useDataExportImport();
  const [format, setFormat] = useState<DataFormat>(DataFormat.CSV);
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [includeHeaders, setIncludeHeaders] = useState(true);
  const [validateData, setValidateData] = useState(true);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);

  const handleExport = async () => {
    try {
      const options: ExportOptions = {
        format,
        type,
        startDate: startDate || undefined,
        endDate: endDate || undefined,
        includeHeaders,
      };

      const result = await exportData.mutateAsync(options);
      await downloadExportedFile(result.url, result.filename);

      toast({
        title: 'エクスポートが完了しました',
        description: `${result.totalRecords}件のデータをエクスポートしました`,
        status: 'success',
        duration: 5000,
        isClosable: true,
      });

      onClose();
    } catch (error) {
      toast({
        title: 'エクスポートに失敗しました',
        description: 'もう一度お試しください',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  const handleImport = async () => {
    if (!selectedFile) return;

    try {
      const options: ImportOptions = {
        format,
        type,
        skipHeaders: !includeHeaders,
        validateData,
      };

      const result = await importData.mutateAsync({
        file: selectedFile,
        options,
      });

      toast({
        title: 'インポートが完了しました',
        description: `${result.successCount}件のデータをインポートしました（エラー: ${result.errorCount}件）`,
        status: result.errorCount > 0 ? 'warning' : 'success',
        duration: 5000,
        isClosable: true,
      });

      if (result.errors && result.errors.length > 0) {
        console.error('インポートエラー:', result.errors);
      }

      onClose();
    } catch (error) {
      toast({
        title: 'インポートに失敗しました',
        description: 'もう一度お試しください',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setSelectedFile(file);
    }
  };

  const isLoading = exportData.isLoading || importData.isLoading;

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="xl">
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>
          {mode === 'export' ? 'データのエクスポート' : 'データのインポート'}
        </ModalHeader>
        <ModalBody>
          <VStack spacing={4}>
            <FormControl>
              <FormLabel>ファイル形式</FormLabel>
              <Select
                value={format}
                onChange={(e) => setFormat(e.target.value as DataFormat)}
              >
                <option value={DataFormat.CSV}>CSV</option>
                <option value={DataFormat.JSON}>JSON</option>
                <option value={DataFormat.EXCEL}>Excel</option>
              </Select>
            </FormControl>

            {mode === 'export' && (
              <>
                <FormControl>
                  <FormLabel>期間</FormLabel>
                  <VStack spacing={2}>
                    <Input
                      type="date"
                      value={startDate}
                      onChange={(e) => setStartDate(e.target.value)}
                      placeholder="開始日"
                    />
                    <Input
                      type="date"
                      value={endDate}
                      onChange={(e) => setEndDate(e.target.value)}
                      placeholder="終了日"
                    />
                  </VStack>
                  <FormHelperText>
                    期間を指定しない場合、全期間のデータがエクスポートされます
                  </FormHelperText>
                </FormControl>
              </>
            )}

            {mode === 'import' && (
              <FormControl>
                <FormLabel>インポートファイル</FormLabel>
                <Input
                  type="file"
                  accept={`.${format}`}
                  onChange={handleFileChange}
                  ref={fileInputRef}
                />
                <FormHelperText>
                  {format === DataFormat.CSV && '.csvファイル'}
                  {format === DataFormat.JSON && '.jsonファイル'}
                  {format === DataFormat.EXCEL &&
                    '.xlsx, .xlsファイル'}
                  をアップロードしてください
                </FormHelperText>
              </FormControl>
            )}

            <FormControl>
              <FormLabel>ヘッダー行</FormLabel>
              <Switch
                isChecked={includeHeaders}
                onChange={(e) => setIncludeHeaders(e.target.checked)}
              />
              <FormHelperText>
                {mode === 'export'
                  ? 'ヘッダー行を含める'
                  : 'ファイルにヘッダー行が含まれている'}
              </FormHelperText>
            </FormControl>

            {mode === 'import' && (
              <FormControl>
                <FormLabel>データの検証</FormLabel>
                <Switch
                  isChecked={validateData}
                  onChange={(e) => setValidateData(e.target.checked)}
                />
                <FormHelperText>
                  インポート前にデータの形式を検証する
                </FormHelperText>
              </FormControl>
            )}
          </VStack>

          {isLoading && (
            <Box mt={4}>
              <Text mb={2}>処理中...</Text>
              <Progress size="xs" isIndeterminate />
            </Box>
          )}
        </ModalBody>

        <ModalFooter>
          <Button variant="ghost" mr={3} onClick={onClose} isDisabled={isLoading}>
            キャンセル
          </Button>
          <Button
            colorScheme="blue"
            onClick={mode === 'export' ? handleExport : handleImport}
            isLoading={isLoading}
            isDisabled={mode === 'import' && !selectedFile}
          >
            {mode === 'export' ? 'エクスポート' : 'インポート'}
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
} 