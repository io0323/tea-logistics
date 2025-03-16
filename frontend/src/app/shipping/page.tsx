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
 * 出荷管理ページコンポーネント
 */
export default function ShippingPage() {
  const router = useRouter();
  const bgColor = useColorModeValue('white', 'gray.800');

  // モックデータ（実際のアプリケーションではAPIから取得）
  const shippingList = [
    {
      id: '1',
      orderNumber: 'ORD-2024-001',
      customerName: '株式会社茶商',
      status: '出荷準備中',
      shippingDate: '2024-03-20',
      deliveryDate: '2024-03-22',
    },
    {
      id: '2',
      orderNumber: 'ORD-2024-002',
      customerName: '茶葉工業株式会社',
      status: '出荷完了',
      shippingDate: '2024-03-19',
      deliveryDate: '2024-03-21',
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
                  出荷管理
                </Heading>
                <Text color="gray.600">
                  出荷情報の一覧を表示します
                </Text>
              </Box>

              <Box bg={useColorModeValue('white', 'gray.700')} p={6} rounded="lg" shadow="base">
                <Table variant="simple">
                  <Thead>
                    <Tr>
                      <Th>注文番号</Th>
                      <Th>顧客名</Th>
                      <Th>ステータス</Th>
                      <Th>出荷予定日</Th>
                      <Th>配送予定日</Th>
                      <Th>アクション</Th>
                    </Tr>
                  </Thead>
                  <Tbody>
                    {shippingList.map((shipping) => (
                      <Tr key={shipping.id}>
                        <Td>{shipping.orderNumber}</Td>
                        <Td>{shipping.customerName}</Td>
                        <Td>
                          <Badge
                            colorScheme={shipping.status === '出荷完了' ? 'green' : 'orange'}
                          >
                            {shipping.status}
                          </Badge>
                        </Td>
                        <Td>{shipping.shippingDate}</Td>
                        <Td>{shipping.deliveryDate}</Td>
                        <Td>
                          <Button
                            size="sm"
                            colorScheme="blue"
                            leftIcon={<FiEye />}
                            onClick={() => router.push(`/shipping/${shipping.id}`)}
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