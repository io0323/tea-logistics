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

interface ShippingDetailPageProps {
  params: {
    id: string;
  };
}

/**
 * 出荷詳細ページコンポーネント
 */
export default function ShippingDetailPage({ params }: ShippingDetailPageProps) {
  const router = useRouter();
  const bgColor = useColorModeValue('white', 'gray.800');

  // モックデータ（実際のアプリケーションではAPIから取得）
  const shippingData = {
    id: params.id,
    orderNumber: 'ORD-2024-001',
    status: '出荷準備中',
    customerName: '株式会社茶商',
    shippingDate: '2024-03-20',
    deliveryDate: '2024-03-22',
    address: '東京都中央区日本橋1-1-1',
    items: [
      { id: 1, name: '煎茶', quantity: 10, unit: 'kg' },
      { id: 2, name: '玉露', quantity: 5, unit: 'kg' },
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
                  colorScheme={shippingData.status === '出荷完了' ? 'green' : 'orange'}
                  fontSize="md"
                  px={3}
                  py={1}
                  borderRadius="full"
                >
                  {shippingData.status}
                </Badge>
              </HStack>

              <Box>
                <Heading size="lg" mb={2}>
                  出荷詳細 #{shippingData.orderNumber}
                </Heading>
                <Text color="gray.600">
                  {shippingData.customerName}様の出荷情報
                </Text>
              </Box>

              <Box bg={useColorModeValue('white', 'gray.700')} p={6} rounded="lg" shadow="base">
                <VStack spacing={4} align="stretch">
                  <HStack>
                    <Text fontWeight="bold" minW="150px">注文番号:</Text>
                    <Text>{shippingData.orderNumber}</Text>
                  </HStack>
                  <HStack>
                    <Text fontWeight="bold" minW="150px">顧客名:</Text>
                    <Text>{shippingData.customerName}</Text>
                  </HStack>
                  <HStack>
                    <Text fontWeight="bold" minW="150px">出荷予定日:</Text>
                    <Text>{shippingData.shippingDate}</Text>
                  </HStack>
                  <HStack>
                    <Text fontWeight="bold" minW="150px">配送予定日:</Text>
                    <Text>{shippingData.deliveryDate}</Text>
                  </HStack>
                  <HStack>
                    <Text fontWeight="bold" minW="150px">配送先住所:</Text>
                    <Text>{shippingData.address}</Text>
                  </HStack>
                </VStack>
              </Box>

              <Box bg={useColorModeValue('white', 'gray.700')} p={6} rounded="lg" shadow="base">
                <Heading size="md" mb={4}>出荷商品一覧</Heading>
                <Table variant="simple">
                  <Thead>
                    <Tr>
                      <Th>商品名</Th>
                      <Th isNumeric>数量</Th>
                      <Th>単位</Th>
                    </Tr>
                  </Thead>
                  <Tbody>
                    {shippingData.items.map((item) => (
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