'use client';

import {
  Box,
  Container,
  Heading,
  Text,
  VStack,
  HStack,
  Badge,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Button,
  useColorModeValue,
} from '@chakra-ui/react';
import { useRouter } from 'next/navigation';
import { FiArrowLeft } from 'react-icons/fi';
import AuthGuard from '@/components/auth/AuthGuard';
import DashboardLayout from '@/components/layout/DashboardLayout';

interface ReceivingDetailPageProps {
  params: {
    id: string;
  };
}

/**
 * 入荷詳細ページコンポーネント
 */
export default function ReceivingDetailPage({ params }: ReceivingDetailPageProps) {
  const router = useRouter();
  const bgColor = useColorModeValue('white', 'gray.800');

  // モックデータ（実際のアプリケーションではAPIから取得）
  const receivingData = {
    id: params.id,
    orderNumber: 'PO-2024-001',
    status: '入荷準備中',
    supplierName: '山田茶園',
    receivingDate: '2024-03-20',
    arrivalDate: '2024-03-22',
    warehouse: '東京倉庫',
    items: [
      { id: 1, name: '煎茶（一番茶）', quantity: 50, unit: 'kg' },
      { id: 2, name: '玉露（新芽）', quantity: 30, unit: 'kg' },
    ],
  };

  return (
    <AuthGuard>
      <DashboardLayout>
        <Box bg={bgColor} minH="100vh" pt={16}>
          <Container maxW="container.xl" py={8}>
            <VStack spacing={8} align="stretch">
              <HStack justify="space-between">
                <Button
                  leftIcon={<FiArrowLeft />}
                  variant="ghost"
                  onClick={() => router.back()}
                >
                  戻る
                </Button>
                <Badge
                  colorScheme={receivingData.status === '入荷完了' ? 'green' : 'orange'}
                  fontSize="md"
                  px={3}
                  py={1}
                  borderRadius="full"
                >
                  {receivingData.status}
                </Badge>
              </HStack>

              <Box>
                <Heading size="lg" mb={2}>
                  入荷詳細 #{receivingData.orderNumber}
                </Heading>
                <Text color="gray.600">
                  {receivingData.supplierName}からの入荷情報
                </Text>
              </Box>

              <Box bg={useColorModeValue('white', 'gray.700')} p={6} rounded="lg" shadow="base">
                <VStack spacing={4} align="stretch">
                  <HStack>
                    <Text fontWeight="bold" minW="150px">発注番号:</Text>
                    <Text>{receivingData.orderNumber}</Text>
                  </HStack>
                  <HStack>
                    <Text fontWeight="bold" minW="150px">仕入先:</Text>
                    <Text>{receivingData.supplierName}</Text>
                  </HStack>
                  <HStack>
                    <Text fontWeight="bold" minW="150px">入荷予定日:</Text>
                    <Text>{receivingData.receivingDate}</Text>
                  </HStack>
                  <HStack>
                    <Text fontWeight="bold" minW="150px">到着予定日:</Text>
                    <Text>{receivingData.arrivalDate}</Text>
                  </HStack>
                  <HStack>
                    <Text fontWeight="bold" minW="150px">入荷倉庫:</Text>
                    <Text>{receivingData.warehouse}</Text>
                  </HStack>
                </VStack>
              </Box>

              <Box bg={useColorModeValue('white', 'gray.700')} p={6} rounded="lg" shadow="base">
                <Heading size="md" mb={4}>入荷商品一覧</Heading>
                <Table variant="simple">
                  <Thead>
                    <Tr>
                      <Th>商品名</Th>
                      <Th isNumeric>数量</Th>
                      <Th>単位</Th>
                    </Tr>
                  </Thead>
                  <Tbody>
                    {receivingData.items.map((item) => (
                      <Tr key={item.id}>
                        <Td>{item.name}</Td>
                        <Td isNumeric>{item.quantity}</Td>
                        <Td>{item.unit}</Td>
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