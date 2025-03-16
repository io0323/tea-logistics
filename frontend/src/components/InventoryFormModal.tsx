'use client';

import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  Button,
  FormControl,
  FormLabel,
  Input,
  Select,
  VStack,
  useToast,
  NumberInput,
  NumberInputField,
  Textarea,
} from '@chakra-ui/react';
import { useState, useEffect } from 'react';
import { InventoryOperation, InventoryOperationType, InventoryOperationStatus } from '@/types/inventory';

interface Inventory {
  id: string;
  productId: string;
  quantity: number;
  type: 'in' | 'out';
  note?: string;
  createdAt: string;
  updatedAt: string;
}

interface InventoryFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: Partial<Inventory>) => Promise<void>;
  initialData?: Inventory | null;
}

/**
 * 在庫情報入力モーダルコンポーネント
 */
export default function InventoryFormModal({
  isOpen,
  onClose,
  onSubmit,
  initialData,
}: InventoryFormModalProps) {
  const toast = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const [formData, setFormData] = useState<Partial<Inventory>>({
    productId: '',
    quantity: 0,
    type: 'in',
    note: '',
  });

  useEffect(() => {
    if (initialData) {
      setFormData({
        productId: initialData.productId,
        quantity: initialData.quantity,
        type: initialData.type,
        note: initialData.note || '',
      });
    } else {
      setFormData({
        productId: '',
        quantity: 0,
        type: 'in',
        note: '',
      });
    }
  }, [initialData]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      await onSubmit(formData);
      toast({
        title: '保存しました',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
      onClose();
    } catch (error) {
      toast({
        title: 'エラーが発生しました',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <form onSubmit={handleSubmit}>
          <ModalHeader>
            {initialData ? '在庫情報の編集' : '新規在庫情報'}
          </ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <VStack spacing={4}>
              <FormControl isRequired>
                <FormLabel>商品ID</FormLabel>
                <Input
                  value={formData.productId}
                  onChange={(e) =>
                    setFormData({ ...formData, productId: e.target.value })
                  }
                  placeholder="商品IDを入力"
                />
              </FormControl>

              <FormControl isRequired>
                <FormLabel>数量</FormLabel>
                <NumberInput
                  value={formData.quantity}
                  onChange={(value) =>
                    setFormData({ ...formData, quantity: Number(value) })
                  }
                  min={0}
                >
                  <NumberInputField />
                </NumberInput>
              </FormControl>

              <FormControl isRequired>
                <FormLabel>タイプ</FormLabel>
                <Select
                  value={formData.type}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      type: e.target.value as 'in' | 'out',
                    })
                  }
                >
                  <option value="in">入庫</option>
                  <option value="out">出庫</option>
                </Select>
              </FormControl>

              <FormControl>
                <FormLabel>メモ</FormLabel>
                <Textarea
                  value={formData.note}
                  onChange={(e) =>
                    setFormData({ ...formData, note: e.target.value })
                  }
                  placeholder="メモを入力"
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
              type="submit"
              isLoading={isLoading}
            >
              {initialData ? '更新' : '作成'}
            </Button>
          </ModalFooter>
        </form>
      </ModalContent>
    </Modal>
  );
} 