'use client';

import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Badge,
  Skeleton,
  Box,
} from '@chakra-ui/react';
import { StockHistory } from '@/types/stockHistory';
import { format } from 'date-fns';
import { ja } from 'date-fns/locale';

interface StockHistoryTableProps {
  histories: StockHistory[];
  isLoading: boolean;
}

/**
 * 在庫履歴テーブルコンポーネント
 */
export default function StockHistoryTable({
  histories,
  isLoading,
}: StockHistoryTableProps) {
  const getTypeColor = (type: StockHistory['type']) => {
    switch (type) {
      case 'in':
        return 'green';
      case 'out':
        return 'red';
      case 'adjustment':
        return 'yellow';
      default:
        return 'gray';
    }
  };

  const getTypeLabel = (type: StockHistory['type']) => {
    switch (type) {
      case 'in':
        return '入庫';
      case 'out':
        return '出庫';
      case 'adjustment':
        return '調整';
      default:
        return type;
    }
  };

  if (isLoading) {
    return (
      <Box>
        <Skeleton height="40px" mb={2} />
        <Skeleton height="40px" mb={2} />
        <Skeleton height="40px" mb={2} />
        <Skeleton height="40px" mb={2} />
        <Skeleton height="40px" />
      </Box>
    );
  }

  return (
    <Table variant="simple">
      <Thead>
        <Tr>
          <Th>日時</Th>
          <Th>種類</Th>
          <Th isNumeric>変更前</Th>
          <Th isNumeric>変更後</Th>
          <Th isNumeric>変更量</Th>
          <Th>理由</Th>
          <Th>操作者</Th>
        </Tr>
      </Thead>
      <Tbody>
        {histories.map((history) => (
          <Tr key={history.id}>
            <Td>
              {format(new Date(history.createdAt), 'yyyy/MM/dd HH:mm', {
                locale: ja,
              })}
            </Td>
            <Td>
              <Badge colorScheme={getTypeColor(history.type)}>
                {getTypeLabel(history.type)}
              </Badge>
            </Td>
            <Td isNumeric>{history.previousStock}</Td>
            <Td isNumeric>{history.newStock}</Td>
            <Td isNumeric>{history.changeAmount}</Td>
            <Td>{history.reason}</Td>
            <Td>{history.createdBy}</Td>
          </Tr>
        ))}
      </Tbody>
    </Table>
  );
} 