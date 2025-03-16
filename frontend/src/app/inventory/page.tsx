'use client';

import { useState } from 'react';
import {
  Box,
  Button,
  Container,
  Heading,
  Input,
  InputGroup,
  InputLeftElement,
  Stack,
  useToast,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  useColorModeValue,
} from '@chakra-ui/react';
import { FiSearch, FiPlus } from 'react-icons/fi';
import DashboardLayout from '@/components/DashboardLayout';
import InventoryTable from '@/components/InventoryTable';
import InventoryFormModal from '@/components/InventoryFormModal';
import Pagination from '@/components/Pagination';
import { useInventory } from '@/hooks/useInventory';
import AuthGuard from '@/components/AuthGuard';

/**
 * 在庫管理ページ
 */
export default function InventoryPage() {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [searchQuery, setSearchQuery] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingInventory, setEditingInventory] = useState<any>(null);
  const toast = useToast();

  const {
    inventories,
    loading,
    error,
    totalPages,
    createInventory,
    updateInventory,
    deleteInventory,
  } = useInventory({
    page,
    pageSize,
    searchQuery,
  });

  const handleCreate = async (data: any) => {
    try {
      await createInventory(data);
      setIsModalOpen(false);
      toast({
        title: '作成成功',
        description: '在庫情報を作成しました',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    } catch (error) {
      toast({
        title: 'エラー',
        description: '在庫情報の作成に失敗しました',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
  };

  const handleEdit = async (data: any) => {
    if (editingInventory) {
      try {
        await updateInventory(editingInventory.id, data);
        setIsModalOpen(false);
        setEditingInventory(null);
        toast({
          title: '更新成功',
          description: '在庫情報を更新しました',
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
      } catch (error) {
        toast({
          title: 'エラー',
          description: '在庫情報の更新に失敗しました',
          status: 'error',
          duration: 3000,
          isClosable: true,
        });
      }
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteInventory(id);
    } catch (error) {
      toast({
        title: 'エラー',
        description: '在庫情報の削除に失敗しました',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    }
  };

  return (
    <AuthGuard>
      <DashboardLayout>
        <Container maxW="container.xl" py={8}>
          <Stack spacing={6}>
            <Box>
              <Heading size="lg">在庫管理</Heading>
              <Button
                leftIcon={<FiPlus />}
                colorScheme="blue"
                mt={4}
                onClick={() => {
                  setEditingInventory(null);
                  setIsModalOpen(true);
                }}
              >
                新規作成
              </Button>
            </Box>

            <InputGroup>
              <InputLeftElement pointerEvents="none">
                <FiSearch color="gray.300" />
              </InputLeftElement>
              <Input
                placeholder="商品名で検索"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </InputGroup>

            {error ? (
              <Box color="red.500">エラーが発生しました: {error.message}</Box>
            ) : (
              <>
                <Box bg="white" rounded="lg" shadow="md" overflow="hidden">
                  <Table variant="simple">
                    <Thead>
                      <Tr>
                        <Th>商品コード</Th>
                        <Th>商品名</Th>
                        <Th>カテゴリー</Th>
                        <Th>在庫数</Th>
                        <Th>単位</Th>
                        <Th>操作</Th>
                      </Tr>
                    </Thead>
                    <Tbody>
                      {inventories.map((inventory) => (
                        <Tr key={inventory.id}>
                          <Td>{inventory.code}</Td>
                          <Td>{inventory.name}</Td>
                          <Td>{inventory.category}</Td>
                          <Td>{inventory.stock}</Td>
                          <Td>{inventory.unit}</Td>
                          <Td>
                            <Button size="sm" colorScheme="blue" mr={2} onClick={() => {
                              setEditingInventory(inventory);
                              setIsModalOpen(true);
                            }}>
                              編集
                            </Button>
                            <Button size="sm" colorScheme="red" onClick={() => handleDelete(inventory.id)}>
                              削除
                            </Button>
                          </Td>
                        </Tr>
                      ))}
                    </Tbody>
                  </Table>
                </Box>

                <Pagination
                  currentPage={page}
                  totalPages={totalPages}
                  pageSize={pageSize}
                  onPageChange={setPage}
                  onPageSizeChange={setPageSize}
                />
              </>
            )}
          </Stack>

          <InventoryFormModal
            isOpen={isModalOpen}
            onClose={() => {
              setIsModalOpen(false);
              setEditingInventory(null);
            }}
            onSubmit={editingInventory ? handleEdit : handleCreate}
            initialData={editingInventory}
          />
        </Container>
      </DashboardLayout>
    </AuthGuard>
  );
} 