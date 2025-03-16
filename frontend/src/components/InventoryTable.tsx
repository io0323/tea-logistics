'use client';

import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  IconButton,
  useToast,
} from '@chakra-ui/react';
import { FiEdit2, FiTrash2 } from 'react-icons/fi';

interface Inventory {
  id: string;
  productId: string;
  quantity: number;
  type: 'in' | 'out';
  note?: string;
  createdAt: string;
  updatedAt: string;
}

interface InventoryTableProps {
  inventories: Inventory[];
  onEdit: (inventory: Inventory) => void;
  onDelete: (id: string) => void;
}

/**
 * 在庫一覧テーブルコンポーネント
 */
export default function InventoryTable({
  inventories,
  onEdit,
  onDelete,
}: InventoryTableProps) {
  const toast = useToast();

  const handleDelete = async (id: string) => {
    if (window.confirm('この在庫情報を削除してもよろしいですか？')) {
      try {
        await onDelete(id);
        toast({
          title: '削除成功',
          description: '在庫情報を削除しました',
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
      } catch (error) {
        toast({
          title: 'エラー',
          description: '在庫情報の削除に失敗しました',
          status: 'error',
          duration: 3000,
          isClosable: true,
        });
      }
    }
  };

  return (
    <Table variant="simple">
      <Thead>
        <Tr>
          <Th>商品ID</Th>
          <Th>数量</Th>
          <Th>タイプ</Th>
          <Th>メモ</Th>
          <Th>作成日時</Th>
          <Th>操作</Th>
        </Tr>
      </Thead>
      <Tbody>
        {inventories.map((inventory) => (
          <Tr key={inventory.id}>
            <Td>{inventory.productId}</Td>
            <Td>{inventory.quantity}</Td>
            <Td>{inventory.type === 'in' ? '入庫' : '出庫'}</Td>
            <Td>{inventory.note || '-'}</Td>
            <Td>{new Date(inventory.createdAt).toLocaleString()}</Td>
            <Td>
              <IconButton
                aria-label="編集"
                icon={<FiEdit2 />}
                size="sm"
                mr={2}
                onClick={() => onEdit(inventory)}
                colorScheme="blue"
              />
              <IconButton
                aria-label="削除"
                icon={<FiTrash2 />}
                size="sm"
                onClick={() => handleDelete(inventory.id)}
                colorScheme="red"
              />
            </Td>
          </Tr>
        ))}
      </Tbody>
    </Table>
  );
} 