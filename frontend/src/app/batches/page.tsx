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
  IconButton,
  useDisclosure,
  Badge,
  Tooltip,
  Collapse,
  VStack,
  Code,
  Spinner,
  Input,
} from '@chakra-ui/react';
import {
  AddIcon,
  DeleteIcon,
  RepeatIcon,
  TimeIcon,
  ViewIcon,
  WarningIcon,
} from '@chakra-ui/icons';
import { useBatch } from '@/hooks/useBatch';
import {
  BatchType,
  BatchStatus,
  BatchConfig,
  BatchResult,
} from '@/types/batch';
import DashboardLayout from '@/components/DashboardLayout';
import AuthGuard from '@/components/AuthGuard';
import BatchConfigModal from '@/components/BatchConfigModal';
import Pagination from '@/components/Pagination';
import SearchBar from '@/components/SearchBar';

/**
 * バッチ処理一覧ページ
 */
export default function BatchesPage() {
  const toast = useToast();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [type, setType] = useState<BatchType | ''>('');
  const [status, setStatus] = useState<BatchStatus | ''>('');
  const [search, setSearch] = useState('');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [selectedBatch, setSelectedBatch] = useState<BatchResult | null>(null);
  const [expandedLogs, setExpandedLogs] = useState<number[]>([]);

  const {
    getBatches,
    createBatch,
    executeBatch,
    cancelBatch,
    deleteBatch,
    getBatchLogs,
  } = useBatch();

  const { data, isLoading } = getBatches({
    type: type as BatchType,
    status: status as BatchStatus,
    startDate,
    endDate,
    page,
    limit: pageSize,
  });

  const { data: logs } = getBatchLogs(selectedBatch?.id || 0);

  const handleCreate = async (config: BatchConfig) => {
    try {
      await createBatch.mutateAsync(config);
      toast({
        title: 'バッチ処理を作成しました',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
      onClose();
    } catch (error) {
      toast({
        title: '作成に失敗しました',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
  };

  const handleExecute = async (id: number) => {
    try {
      await executeBatch.mutateAsync(id);
      toast({
        title: 'バッチ処理を実行しました',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    } catch (error) {
      toast({
        title: '実行に失敗しました',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('このバッチ処理を削除してもよろしいですか？')) {
      try {
        await deleteBatch.mutateAsync(id);
        toast({
          title: 'バッチ処理を削除しました',
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

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };

  const handlePageSizeChange = (newPageSize: number) => {
    setPageSize(newPageSize);
    setPage(1);
  };

  const toggleLogs = (id: number) => {
    setExpandedLogs((prev) =>
      prev.includes(id)
        ? prev.filter((logId) => logId !== id)
        : [...prev, id]
    );
  };

  const getTypeLabel = (type: BatchType): string => {
    const labels: Record<BatchType, string> = {
      [BatchType.STOCK_CHECK]: '在庫確認',
      [BatchType.DELIVERY_STATUS_UPDATE]: '配送ステータス更新',
      [BatchType.DATA_CLEANUP]: 'データクリーンアップ',
      [BatchType.REPORT_GENERATION]: 'レポート生成',
    };
    return labels[type];
  };

  const getStatusColor = (status: BatchStatus): string => {
    const colors: Record<BatchStatus, string> = {
      [BatchStatus.PENDING]: 'yellow',
      [BatchStatus.RUNNING]: 'blue',
      [BatchStatus.COMPLETED]: 'green',
      [BatchStatus.FAILED]: 'red',
      [BatchStatus.CANCELLED]: 'gray',
    };
    return colors[status];
  };

  return (
    <AuthGuard>
      <DashboardLayout>
        <Box p={6}>
          <HStack justify="space-between" mb={6}>
            <Text fontSize="2xl" fontWeight="bold">
              バッチ処理管理
            </Text>
            <Button
              leftIcon={<AddIcon />}
              colorScheme="blue"
              onClick={onOpen}
            >
              新規バッチ
            </Button>
          </HStack>

          <HStack spacing={4} mb={6} flexWrap="wrap">
            <Select
              placeholder="処理タイプ"
              value={type}
              onChange={(e) => setType(e.target.value as BatchType)}
              width="200px"
            >
              {Object.values(BatchType).map((t) => (
                <option key={t} value={t}>
                  {getTypeLabel(t)}
                </option>
              ))}
            </Select>

            <Select
              placeholder="ステータス"
              value={status}
              onChange={(e) => setStatus(e.target.value as BatchStatus)}
              width="200px"
            >
              {Object.values(BatchStatus).map((s) => (
                <option key={s} value={s}>
                  {s}
                </option>
              ))}
            </Select>

            <Input
              type="date"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
              width="200px"
              placeholder="開始日"
            />

            <Input
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
              width="200px"
              placeholder="終了日"
            />

            <SearchBar
              placeholder="バッチ処理を検索"
              value={search}
              onChange={setSearch}
              width="300px"
            />
          </HStack>

          <Box overflowX="auto" position="relative" minHeight="400px">
            {isLoading ? (
              <Box
                position="absolute"
                top="50%"
                left="50%"
                transform="translate(-50%, -50%)"
              >
                <Spinner size="xl" />
              </Box>
            ) : data?.items.length === 0 ? (
              <Box
                textAlign="center"
                py={10}
                color="gray.500"
              >
                バッチ処理が見つかりませんでした
              </Box>
            ) : (
              <Table variant="simple">
                <Thead>
                  <Tr>
                    <Th>ID</Th>
                    <Th>タイプ</Th>
                    <Th>ステータス</Th>
                    <Th>開始時間</Th>
                    <Th>終了時間</Th>
                    <Th>処理件数</Th>
                    <Th>成功</Th>
                    <Th>エラー</Th>
                    <Th>操作</Th>
                  </Tr>
                </Thead>
                <Tbody>
                  {data?.items.map((batch) => (
                    <>
                      <Tr key={batch.id}>
                        <Td>{batch.id}</Td>
                        <Td>{getTypeLabel(batch.type)}</Td>
                        <Td>
                          <Badge colorScheme={getStatusColor(batch.status)}>
                            {batch.status}
                          </Badge>
                        </Td>
                        <Td>{batch.startTime}</Td>
                        <Td>{batch.endTime || '-'}</Td>
                        <Td>{batch.processedItems}</Td>
                        <Td>{batch.successCount}</Td>
                        <Td>
                          {batch.errorCount > 0 ? (
                            <HStack>
                              <Text color="red.500">{batch.errorCount}</Text>
                              <Tooltip
                                label="エラーの詳細を確認"
                                hasArrow
                              >
                                <WarningIcon color="red.500" />
                              </Tooltip>
                            </HStack>
                          ) : (
                            0
                          )}
                        </Td>
                        <Td>
                          <HStack spacing={2}>
                            <Tooltip label="ログを表示" hasArrow>
                              <IconButton
                                aria-label="View logs"
                                icon={<ViewIcon />}
                                size="sm"
                                variant="ghost"
                                onClick={() => toggleLogs(batch.id)}
                              />
                            </Tooltip>
                            {batch.status === BatchStatus.PENDING && (
                              <Tooltip label="実行" hasArrow>
                                <IconButton
                                  aria-label="Execute batch"
                                  icon={<RepeatIcon />}
                                  size="sm"
                                  variant="ghost"
                                  colorScheme="blue"
                                  onClick={() => handleExecute(batch.id)}
                                />
                              </Tooltip>
                            )}
                            {batch.status === BatchStatus.RUNNING && (
                              <Tooltip label="キャンセル" hasArrow>
                                <IconButton
                                  aria-label="Cancel batch"
                                  icon={<TimeIcon />}
                                  size="sm"
                                  variant="ghost"
                                  colorScheme="yellow"
                                  onClick={() => cancelBatch.mutate(batch.id)}
                                />
                              </Tooltip>
                            )}
                            <Tooltip label="削除" hasArrow>
                              <IconButton
                                aria-label="Delete batch"
                                icon={<DeleteIcon />}
                                size="sm"
                                variant="ghost"
                                colorScheme="red"
                                onClick={() => handleDelete(batch.id)}
                              />
                            </Tooltip>
                          </HStack>
                        </Td>
                      </Tr>
                      <Tr>
                        <Td colSpan={9} p={0}>
                          <Collapse in={expandedLogs.includes(batch.id)}>
                            <Box p={4} bg="gray.50">
                              <VStack align="stretch" spacing={2}>
                                <Text fontWeight="bold">実行ログ</Text>
                                <Code p={2} borderRadius="md">
                                  {batch.logs?.join('\n') || 'ログはありません'}
                                </Code>
                              </VStack>
                            </Box>
                          </Collapse>
                        </Td>
                      </Tr>
                    </>
                  ))}
                </Tbody>
              </Table>
            )}
          </Box>

          {data && data.items.length > 0 && (
            <Pagination
              currentPage={page}
              totalPages={data.totalPages}
              pageSize={pageSize}
              totalItems={data.total}
              onPageChange={handlePageChange}
              onPageSizeChange={handlePageSizeChange}
            />
          )}

          <BatchConfigModal
            isOpen={isOpen}
            onClose={onClose}
            onSubmit={handleCreate}
            isLoading={createBatch.isLoading}
          />
        </Box>
      </DashboardLayout>
    </AuthGuard>
  );
} 