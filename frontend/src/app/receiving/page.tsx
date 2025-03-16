'use client';

import {
  Box,
  Container,
  Heading,
  Text,
  VStack,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Button,
  HStack,
  Badge,
  useColorModeValue,
} from '@chakra-ui/react';
import { useRouter } from 'next/navigation';
import { FiEye } from 'react-icons/fi';
import AuthGuard from '@/components/auth/AuthGuard';
import DashboardLayout from '@/components/layout/DashboardLayout';

/**
 * 入荷管理ページコンポーネント
 */
export default function ReceivingPage() {
  const router = useRouter();
  const bgColor = useColorModeValue('white', 'gray.800');

  // モックデータ（実際のアプリケーションではAPIから取得）
  const receivingList = [
    {
      id: '1',
      orderNumber: 'PO-2024-001',
      supplierName: '山田茶園',
      status: '入荷準備中',
      receivingDate: '2024-03-20',
      arrivalDate: '2024-03-22',
    },
    {
      id: '2',
      orderNumber: 'PO-2024-002',
      supplierName: '緑茶農園',
      status: '入荷完了',
      receivingDate: '2024-03-19',
      arrivalDate: '2024-03-21',
    },
  ];

  return (
    <AuthGuard>
      <DashboardLayout>
        <Box bg={bgColor} minH="100vh" pt={16}>
          <Container maxW="container.xl" py={8}>
            <VStack spacing={8} align="stretch">
              <Box>
                <Heading size="lg" mb={2}>
                  入荷管理
                </Heading>
                <Text color="gray.600">
                  入荷情報の一覧を表示します
                </Text>
              </Box>

              <Box bg={useColorModeValue('white', 'gray.700')} p={6} rounded="lg" shadow="base">
                <Table variant="simple">
                  <Thead>
                    <Tr>
                      <Th>発注番号</Th>
                      <Th>仕入先</Th>
                      <Th>ステータス</Th>
                      <Th>入荷予定日</Th>
                      <Th>到着予定日</Th>
                      <Th>アクション</Th>
                    </Tr>
                  </Thead>
                  <Tbody>
                    {receivingList.map((receiving) => (
                      <Tr key={receiving.id}>
                        <Td>{receiving.orderNumber}</Td>
                        <Td>{receiving.supplierName}</Td>
                        <Td>
                          <Badge
                            colorScheme={receiving.status === '入荷完了' ? 'green' : 'orange'}
                          >
                            {receiving.status}
                          </Badge>
                        </Td>
                        <Td>{receiving.receivingDate}</Td>
                        <Td>{receiving.arrivalDate}</Td>
                        <Td>
                          <Button
                            size="sm"
                            colorScheme="blue"
                            leftIcon={<FiEye />}
                            onClick={() => router.push(`/receiving/${receiving.id}`)}
                          >
                            詳細
                          </Button>
                        </Td>
                      </Tr>
                    ))}
                  </Tbody>
                </Table>
              </Box>
            </VStack>
          </Container>
        </Box>
      </DashboardLayout>
    </AuthGuard>
  );
} 