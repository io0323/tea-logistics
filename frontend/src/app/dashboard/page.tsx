'use client';

import {
  Box,
  Container,
  Grid,
  Heading,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  StatArrow,
  SimpleGrid,
  Card,
  CardHeader,
  CardBody,
  Text,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Badge,
  useColorModeValue,
} from '@chakra-ui/react';
import DashboardLayout from '@/components/DashboardLayout';
import AuthGuard from '@/components/AuthGuard';
import { useAuth } from '@/hooks/useAuth';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  BarElement,
} from 'chart.js';
import { Line, Bar } from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend
);

/**
 * ダッシュボードページ
 */
export default function DashboardPage() {
  const { user } = useAuth();
  const bgColor = useColorModeValue('white', 'gray.700');
  const borderColor = useColorModeValue('gray.200', 'gray.600');

  const salesData = {
    labels: ['1月', '2月', '3月', '4月', '5月', '6月'],
    datasets: [
      {
        label: '売上高',
        data: [650, 590, 800, 810, 760, 850],
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.1,
      },
    ],
  };

  const inventoryData = {
    labels: ['煎茶', '玉露', 'ほうじ茶', '抹茶', '玄米茶'],
    datasets: [
      {
        label: '在庫数',
        data: [100, 50, 80, 30, 60],
        backgroundColor: [
          'rgba(255, 99, 132, 0.5)',
          'rgba(54, 162, 235, 0.5)',
          'rgba(255, 206, 86, 0.5)',
          'rgba(75, 192, 192, 0.5)',
          'rgba(153, 102, 255, 0.5)',
        ],
      },
    ],
  };

  return (
    <AuthGuard>
      <DashboardLayout>
        <Container maxW="container.xl" py={8}>
          <Box mb={6}>
            <Heading size="lg" mb={2}>ダッシュボード</Heading>
            <Text>ようこそ、{user?.name || 'ゲスト'}さん</Text>
          </Box>

          {/* 統計カード */}
          <SimpleGrid columns={{ base: 1, md: 4 }} spacing={6} mb={8}>
            <Box p={6} bg={bgColor} rounded="lg" shadow="md">
              <Stat>
                <StatLabel>在庫総数</StatLabel>
                <StatNumber>1,234</StatNumber>
                <StatHelpText>
                  <StatArrow type="increase" />
                  12%
                </StatHelpText>
              </Stat>
            </Box>
            <Box p={6} bg={bgColor} rounded="lg" shadow="md">
              <Stat>
                <StatLabel>今月の出荷数</StatLabel>
                <StatNumber>456</StatNumber>
                <StatHelpText>
                  <StatArrow type="increase" />
                  8%
                </StatHelpText>
              </Stat>
            </Box>
            <Box p={6} bg={bgColor} rounded="lg" shadow="md">
              <Stat>
                <StatLabel>今月の入荷数</StatLabel>
                <StatNumber>789</StatNumber>
                <StatHelpText>
                  <StatArrow type="increase" />
                  15%
                </StatHelpText>
              </Stat>
            </Box>
            <Box p={6} bg={bgColor} rounded="lg" shadow="md">
              <Stat>
                <StatLabel>売上高</StatLabel>
                <StatNumber>¥2.4M</StatNumber>
                <StatHelpText>
                  <StatArrow type="increase" />
                  23%
                </StatHelpText>
              </Stat>
            </Box>
          </SimpleGrid>

          {/* グラフ */}
          <SimpleGrid columns={{ base: 1, lg: 2 }} spacing={6} mb={8}>
            <Card>
              <CardHeader>
                <Heading size="md">売上推移</Heading>
              </CardHeader>
              <CardBody>
                <Line data={salesData} />
              </CardBody>
            </Card>
            <Card>
              <CardHeader>
                <Heading size="md">商品別在庫状況</Heading>
              </CardHeader>
              <CardBody>
                <Bar data={inventoryData} />
              </CardBody>
            </Card>
          </SimpleGrid>

          {/* 最近の取引 */}
          <Card mb={8}>
            <CardHeader>
              <Heading size="md">最近の取引</Heading>
            </CardHeader>
            <CardBody>
              <Table variant="simple">
                <Thead>
                  <Tr>
                    <Th>日時</Th>
                    <Th>取引種別</Th>
                    <Th>商品名</Th>
                    <Th>数量</Th>
                    <Th>状態</Th>
                  </Tr>
                </Thead>
                <Tbody>
                  <Tr>
                    <Td>2024/03/15 15:30</Td>
                    <Td>出荷</Td>
                    <Td>煎茶 - 特上</Td>
                    <Td>50kg</Td>
                    <Td><Badge colorScheme="green">完了</Badge></Td>
                  </Tr>
                  <Tr>
                    <Td>2024/03/15 14:15</Td>
                    <Td>入荷</Td>
                    <Td>玉露</Td>
                    <Td>30kg</Td>
                    <Td><Badge colorScheme="yellow">処理中</Badge></Td>
                  </Tr>
                  <Tr>
                    <Td>2024/03/15 13:45</Td>
                    <Td>出荷</Td>
                    <Td>ほうじ茶</Td>
                    <Td>100kg</Td>
                    <Td><Badge colorScheme="green">完了</Badge></Td>
                  </Tr>
                </Tbody>
              </Table>
            </CardBody>
          </Card>

          {/* 在庫アラート */}
          <Card>
            <CardHeader>
              <Heading size="md">在庫アラート</Heading>
            </CardHeader>
            <CardBody>
              <SimpleGrid columns={{ base: 1, md: 2 }} spacing={4}>
                <Box p={4} bg="red.50" rounded="md" borderWidth={1} borderColor="red.200">
                  <Text color="red.600" fontWeight="bold">在庫不足</Text>
                  <Text>玉露 - 特上（残り5kg）</Text>
                </Box>
                <Box p={4} bg="yellow.50" rounded="md" borderWidth={1} borderColor="yellow.200">
                  <Text color="yellow.600" fontWeight="bold">在庫要確認</Text>
                  <Text>煎茶 - 上級（残り20kg）</Text>
                </Box>
              </SimpleGrid>
            </CardBody>
          </Card>
        </Container>
      </DashboardLayout>
    </AuthGuard>
  );
} 