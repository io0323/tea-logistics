'use client';

import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
  VStack,
  Text,
  useToast,
  Input,
  FormControl,
  FormLabel,
} from '@chakra-ui/react';
import { FiUpload } from 'react-icons/fi';
import { useState } from 'react';
import { readFile, convertFromCSV } from '@/utils/productExport';

interface ProductImportModalProps {
  isOpen: boolean;
  onClose: () => void;
  onImport: (products: any[]) => Promise<void>;
}

/**
 * 商品データインポートモーダルコンポーネント
 */
export default function ProductImportModal({
  isOpen,
  onClose,
  onImport,
}: ProductImportModalProps) {
  const [file, setFile] = useState<File | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const toast = useToast();

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile && selectedFile.type === 'text/csv') {
      setFile(selectedFile);
    } else {
      toast({
        title: 'エラー',
        description: 'CSVファイルを選択してください',
        status: 'error',
        duration: 3000,
      });
    }
  };

  const handleImport = async () => {
    if (!file) return;

    try {
      setIsLoading(true);
      const content = await readFile(file);
      const products = convertFromCSV(content);
      await onImport(products);
      toast({
        title: 'インポート完了',
        description: `${products.length}件の商品をインポートしました`,
        status: 'success',
        duration: 3000,
      });
      onClose();
    } catch (error) {
      toast({
        title: 'エラー',
        description: 'インポートに失敗しました',
        status: 'error',
        duration: 3000,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>商品データのインポート</ModalHeader>
        <ModalBody>
          <VStack spacing={4}>
            <FormControl>
              <FormLabel>CSVファイル</FormLabel>
              <Input
                type="file"
                accept=".csv"
                onChange={handleFileChange}
                p={1}
              />
            </FormControl>
            <Text fontSize="sm" color="gray.500">
              CSVファイルには以下の列が必要です：
              商品名、カテゴリー、価格、在庫数、説明
            </Text>
          </VStack>
        </ModalBody>
        <ModalFooter>
          <Button variant="ghost" mr={3} onClick={onClose}>
            キャンセル
          </Button>
          <Button
            colorScheme="blue"
            leftIcon={<FiUpload />}
            onClick={handleImport}
            isLoading={isLoading}
            isDisabled={!file}
          >
            インポート
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
} 