'use client';

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
  Input,
  Select,
  VStack,
  useToast,
} from '@chakra-ui/react';
import { useState } from 'react';
import { Product } from '@/types/product';

interface BulkUpdateModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: Partial<Product>) => Promise<void>;
  isLoading: boolean;
  selectedCount: number;
}

/**
 * 商品一括更新モーダルコンポーネント
 */
export default function BulkUpdateModal({
  isOpen,
  onClose,
  onSubmit,
  isLoading,
  selectedCount,
}: BulkUpdateModalProps) {
  const [formData, setFormData] = useState<Partial<Product>>({
    category: '',
    price: undefined,
    stock: undefined,
  });
  const toast = useToast();

  const handleSubmit = async () => {
    try {
      await onSubmit(formData);
      toast({
        title: '商品を一括更新しました',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
      onClose();
    } catch (error) {
      toast({
        title: 'エラーが発生しました',
        description: error instanceof Error ? error.message : '予期せぬエラーが発生しました',
        status: 'error',
        duration: 5000,
        isClosable: true,
      });
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>商品一括更新</ModalHeader>
        <ModalBody>
          <VStack spacing={4}>
            <FormControl>
              <FormLabel>カテゴリー</FormLabel>
              <Select
                value={formData.category}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, category: e.target.value }))
                }
                placeholder="カテゴリーを選択"
              >
                <option value="茶葉">茶葉</option>
                <option value="茶器">茶器</option>
                <option value="菓子">菓子</option>
                <option value="その他">その他</option>
              </Select>
            </FormControl>

            <FormControl>
              <FormLabel>価格</FormLabel>
              <Input
                type="number"
                value={formData.price || ''}
                onChange={(e) =>
                  setFormData((prev) => ({
                    ...prev,
                    price: e.target.value ? Number(e.target.value) : undefined,
                  }))
                }
                placeholder="価格を入力"
              />
            </FormControl>

            <FormControl>
              <FormLabel>在庫数</FormLabel>
              <Input
                type="number"
                value={formData.stock || ''}
                onChange={(e) =>
                  setFormData((prev) => ({
                    ...prev,
                    stock: e.target.value ? Number(e.target.value) : undefined,
                  }))
                }
                placeholder="在庫数を入力"
              />
            </FormControl>
          </VStack>
        </ModalBody>
        <ModalFooter>
          <Button variant="ghost" mr={3} onClick={onClose}>
            キャンセル
          </Button>
          <Button
            colorScheme="blue"
            onClick={handleSubmit}
            isLoading={isLoading}
            isDisabled={!Object.values(formData).some((value) => value !== undefined)}
          >
            更新 ({selectedCount}件)
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
} 