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
  useToast,
} from '@chakra-ui/react';
import { useEffect, useState } from 'react';
import { useStockHistory } from '@/hooks/useStockHistory';
import StockHistoryTable from './StockHistoryTable';
import Pagination from './Pagination';

interface StockHistoryModalProps {
  isOpen: boolean;
  onClose: () => void;
  productId: string;
}

/**
 * 在庫履歴モーダルコンポーネント
 */
export default function StockHistoryModal({
  isOpen,
  onClose,
  productId,
}: StockHistoryModalProps) {
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(10);
  const toast = useToast();

  const {
    histories,
    totalPages,
    isLoading,
    error,
    fetchHistories,
  } = useStockHistory({
    productId,
    page: currentPage,
    pageSize,
  });

  useEffect(() => {
    if (isOpen) {
      fetchHistories();
    }
  }, [isOpen, currentPage, fetchHistories]);

  if (error) {
    toast({
      title: 'エラーが発生しました',
      description: error.message,
      status: 'error',
      duration: 5000,
      isClosable: true,
    });
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="xl">
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>在庫履歴</ModalHeader>
        <ModalBody>
          <VStack spacing={4} align="stretch">
            <StockHistoryTable
              histories={histories}
              isLoading={isLoading}
            />
            <Pagination
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={setCurrentPage}
            />
          </VStack>
        </ModalBody>
        <ModalFooter>
          <Button variant="ghost" onClick={onClose}>
            閉じる
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
} 