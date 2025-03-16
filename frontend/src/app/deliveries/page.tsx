'use client';

import { useState } from 'react';
import {
  Box,
  Button,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Text,
  useToast,
  HStack,
  Select,
  Input,
  IconButton,
  useDisclosure,
  Badge,
} from '@chakra-ui/react';
import { AddIcon, EditIcon, DeleteIcon, DownloadIcon, UploadIcon } from '@chakra-ui/icons';
import { useDelivery } from '@/hooks/useDelivery';
import { Delivery, DeliveryStatus } from '@/types/delivery';
import DashboardLayout from '@/components/DashboardLayout';
import AuthGuard from '@/components/AuthGuard';
import { useAuth } from '@/hooks/useAuth';
import DeliveryFormModal from '@/components/DeliveryFormModal';
import Pagination from '@/components/Pagination';
import { useSettings } from '@/hooks/useSettings';
import SearchBar from '@/components/SearchBar';
import ExportImportModal from '@/components/ExportImportModal';
import { DataType } from '@/types/export';

/**
 * 配送一覧ページ
 */
export default function DeliveriesPage() {
  const { user } = useAuth();
  const toast = useToast();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [status, setStatus] = useState<DeliveryStatus | ''>('');
  const [search, setSearch] = useState('');
  const [selectedDelivery, setSelectedDelivery] = useState<Delivery | undefined>();
  const [exportImportMode, setExportImportMode] = useState<'export' | 'import' | null>(null);

  const { getDeliveries, createDelivery, updateDelivery, deleteDelivery } = useDelivery();
  const { getSettings } = useSettings();
  const { data: settings } = getSettings();

  const { data, isLoading } = getDeliveries({
    page,
    limit: pageSize,
    status: status as DeliveryStatus,
    search,
  });

  const handleCreate = async (data: any) => {
    await createDelivery.mutateAsync(data);
  };

  const handleUpdate = async (data: any) => {
    if (selectedDelivery) {
      await updateDelivery.mutateAsync({
        id: selectedDelivery.id,
        ...data,
      });
    }
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('この配送を削除してもよろしいですか？')) {
      try {
        await deleteDelivery.mutateAsync(id);
        toast({
          title: '配送を削除しました',
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
      } catch (error) {
        toast({
          title: '削除に失敗しました',
          status: 'error',
          duration: 3000,
          isClosable: true,
        });
      }
    }
  };

  const handleOpenModal = (delivery?: Delivery) => {
    setSelectedDelivery(delivery);
    onOpen();
  };

  const handleCloseModal = () => {
    setSelectedDelivery(undefined);
    onClose();
  };

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };

  const handlePageSizeChange = (newPageSize: number) => {
    setPageSize(newPageSize);
    setPage(1); // ページサイズが変更されたら1ページ目に戻る
  };

  const handleOpenExportModal = () => {
    setExportImportMode('export');
  };

  const handleOpenImportModal = () => {
    setExportImportMode('import');
  };

  const handleCloseExportImportModal = () => {
    setExportImportMode(null);
  };

  const getStatusLabel = (status: DeliveryStatus) => {
    const labels: Record<DeliveryStatus, string> = {
      [DeliveryStatus.PENDING]: '保留中',
      [DeliveryStatus.IN_TRANSIT]: '配送中',
      [DeliveryStatus.DELIVERED]: '配送完了',
      [DeliveryStatus.CANCELLED]: 'キャンセル',
    };
    return labels[status];
  };

  const getStatusColor = (status: DeliveryStatus) => {
    const colors: Record<DeliveryStatus, string> = {
      [DeliveryStatus.PENDING]: 'yellow',
      [DeliveryStatus.IN_TRANSIT]: 'blue',
      [DeliveryStatus.DELIVERED]: 'green',
      [DeliveryStatus.CANCELLED]: 'red',
    };
    return colors[status];
  };

  return (
    <AuthGuard>
      <DashboardLayout>
        <Box p={6}>
          <HStack justify="space-between" mb={6}>
            <Text fontSize="2xl" fontWeight="bold">
              配送管理
            </Text>
            <HStack spacing={2}>
              <Button
                leftIcon={<DownloadIcon />}
                variant="outline"
                onClick={handleOpenExportModal}
              >
                エクスポート
              </Button>
              <Button
                leftIcon={<UploadIcon />}
                variant="outline"
                onClick={handleOpenImportModal}
              >
                インポート
              </Button>
              <Button
                leftIcon={<AddIcon />}
                colorScheme="blue"
                onClick={() => handleOpenModal()}
              >
                新規配送
              </Button>
            </HStack>
          </HStack>

          <HStack spacing={4} mb={6}>
            <Select
              placeholder="ステータス"
              value={status}
              onChange={(e) => setStatus(e.target.value as DeliveryStatus)}
              width="200px"
            >
              {Object.values(DeliveryStatus).map((s) => (
                <option key={s} value={s}>
                  {getStatusLabel(s)}
                </option>
              ))}
            </Select>

            <SearchBar
              placeholder="顧客名で検索"
              value={search}
              onChange={setSearch}
              width="300px"
            />
          </HStack>

          <Box overflowX="auto">
            <Table variant="simple">
              <Thead>
                <Tr>
                  <Th>ID</Th>
                  <Th>注文ID</Th>
                  <Th>顧客名</Th>
                  <Th>配送先住所</Th>
                  <Th>電話番号</Th>
                  <Th>ステータス</Th>
                  <Th>予定配送日</Th>
                  <Th>操作</Th>
                </Tr>
              </Thead>
              <Tbody>
                {data?.items.map((delivery) => (
                  <Tr key={delivery.id}>
                    <Td>{delivery.id}</Td>
                    <Td>{delivery.orderId}</Td>
                    <Td>{delivery.customerName}</Td>
                    <Td>{delivery.customerAddress}</Td>
                    <Td>{delivery.customerPhone}</Td>
                    <Td>
                      <Badge colorScheme={getStatusColor(delivery.status)}>
                        {getStatusLabel(delivery.status)}
                      </Badge>
                    </Td>
                    <Td>{delivery.estimatedDeliveryDate}</Td>
                    <Td>
                      <HStack spacing={2}>
                        <IconButton
                          aria-label="Edit delivery"
                          icon={<EditIcon />}
                          size="sm"
                          variant="ghost"
                          onClick={() => handleOpenModal(delivery)}
                        />
                        <IconButton
                          aria-label="Delete delivery"
                          icon={<DeleteIcon />}
                          size="sm"
                          variant="ghost"
                          colorScheme="red"
                          onClick={() => handleDelete(delivery.id)}
                        />
                      </HStack>
                    </Td>
                  </Tr>
                ))}
              </Tbody>
            </Table>
          </Box>

          {data && (
            <Pagination
              currentPage={page}
              totalPages={data.totalPages}
              pageSize={pageSize}
              totalItems={data.total}
              onPageChange={handlePageChange}
              onPageSizeChange={handlePageSizeChange}
            />
          )}

          <DeliveryFormModal
            isOpen={isOpen}
            onClose={handleCloseModal}
            onSubmit={selectedDelivery ? handleUpdate : handleCreate}
            initialData={selectedDelivery}
            title={selectedDelivery ? '配送の編集' : '新規配送'}
          />

          <ExportImportModal
            isOpen={exportImportMode !== null}
            onClose={handleCloseExportImportModal}
            type={DataType.DELIVERY}
            mode={exportImportMode || 'export'}
          />
        </Box>
      </DashboardLayout>
    </AuthGuard>
  );
} 