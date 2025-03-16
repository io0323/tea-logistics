'use client';

import { useState } from 'react';
import {
  Box,
  Container,
  Heading,
  VStack,
  HStack,
  Select,
  Input,
  Button,
  Card,
  CardBody,
  Text,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  StatArrow,
  useColorModeValue,
  SimpleGrid,
} from '@chakra-ui/react';
import { Line, Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { useReport } from '@/hooks/useReport';
import DashboardLayout from '@/components/DashboardLayout';
import AuthGuard from '@/components/AuthGuard';

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

export default function ReportsPage() {
  const [periodType, setPeriodType] = useState('monthly');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');

  const { getReport } = useReport();
  const { data, isLoading } = getReport({
    periodType,
    startDate,
    endDate,
  });

  const cardBg = useColorModeValue('white', 'gray.700');

  if (isLoading) {
    return (
      <AuthGuard>
        <DashboardLayout>
          <Box minH="100vh" bg={useColorModeValue('gray.50', 'gray.900')}>
            <Container maxW="container.xl" py={8}>
              <Text>読み込み中...</Text>
            </Container>
          </Box>
        </DashboardLayout>
      </AuthGuard>
    );
  }

  return (
    <AuthGuard>
      <DashboardLayout>
        <Box minH="100vh" bg={useColorModeValue('gray.50', 'gray.900')}>
          <Container maxW="container.xl" py={8}>
            <SimpleGrid columns={{ base: 1, md: 4 }} spacing={6} mb={8}>
              <Box p={6} bg={cardBg} rounded="lg" shadow="md">
                <Stat>
                  <StatLabel>総売上</StatLabel>
                  <StatNumber>¥{data?.salesReport[5]?.totalSales.toLocaleString() || 0}</StatNumber>
                  <StatHelpText>
                    <StatArrow type="increase" />
                    前月比 +12.3%
                  </StatHelpText>
                </Stat>
              </Box>
              <Box p={6} bg={cardBg} rounded="lg" shadow="md">
                <Stat>
                  <StatLabel>総出荷数</StatLabel>
                  <StatNumber>1,234</StatNumber>
                  <StatHelpText>
                    <StatArrow type="increase" />
                    前月比 +8.5%
                  </StatHelpText>
                </Stat>
              </Box>
              <Box p={6} bg={cardBg} rounded="lg" shadow="md">
                <Stat>
                  <StatLabel>総入荷数</StatLabel>
                  <StatNumber>987</StatNumber>
                  <StatHelpText>
                    <StatArrow type="decrease" />
                    前月比 -3.2%
                  </StatHelpText>
                </Stat>
              </Box>
              <Box p={6} bg={cardBg} rounded="lg" shadow="md">
                <Stat>
                  <StatLabel>配送完了率</StatLabel>
                  <StatNumber>98.5%</StatNumber>
                  <StatHelpText>
                    <StatArrow type="increase" />
                    前月比 +2.1%
                  </StatHelpText>
                </Stat>
              </Box>
            </SimpleGrid>

            <Card mb={6} bg={cardBg}>
              <CardBody>
                <HStack spacing={4}>
                  <Select
                    value={periodType}
                    onChange={(e) => setPeriodType(e.target.value)}
                    w="200px"
                  >
                    <option value="daily">日次</option>
                    <option value="weekly">週次</option>
                    <option value="monthly">月次</option>
                  </Select>
                  <Input
                    type="date"
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                    w="200px"
                  />
                  <Input
                    type="date"
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                    w="200px"
                  />
                  <Button colorScheme="blue">更新</Button>
                </HStack>
              </CardBody>
            </Card>

            <SimpleGrid columns={{ base: 1, lg: 2 }} spacing={6} mb={8}>
              <Card>
                <CardBody>
                  <VStack spacing={4} align="stretch">
                    <Heading size="md">売上推移</Heading>
                    <Box h="300px">
                      <Line
                        data={{
                          labels: data?.salesReport.map((item) => item.period) || [],
                          datasets: [
                            {
                              label: '売上',
                              data: data?.salesReport.map((item) => item.totalSales) || [],
                              borderColor: 'rgb(75, 192, 192)',
                              tension: 0.1,
                            },
                          ],
                        }}
                        options={{
                          responsive: true,
                          plugins: {
                            legend: {
                              position: 'top' as const,
                            },
                          },
                        }}
                      />
                    </Box>
                  </VStack>
                </CardBody>
              </Card>
              <Card>
                <CardBody>
                  <VStack spacing={4} align="stretch">
                    <Heading size="md">配送完了率</Heading>
                    <Box h="300px">
                      <Line
                        data={{
                          labels: data?.deliveryReport.map((item) => item.period) || [],
                          datasets: [
                            {
                              label: '配送完了率',
                              data: data?.deliveryReport.map((item) => item.onTimeDeliveryRate) || [],
                              borderColor: 'rgb(53, 162, 235)',
                              tension: 0.1,
                            },
                          ],
                        }}
                        options={{
                          responsive: true,
                          plugins: {
                            legend: {
                              position: 'top' as const,
                            },
                          },
                        }}
                      />
                    </Box>
                  </VStack>
                </CardBody>
              </Card>
            </SimpleGrid>

            <Card bg={cardBg}>
              <CardBody>
                <VStack spacing={4} align="stretch">
                  <Heading size="md">月間出荷・入荷数量</Heading>
                  <Box h="300px">
                    <Bar
                      data={{
                        labels: ['1月', '2月', '3月', '4月', '5月', '6月'],
                        datasets: [
                          {
                            label: '出荷数',
                            data: [65, 59, 80, 81, 56, 55],
                            backgroundColor: 'rgba(75, 192, 192, 0.5)',
                          },
                          {
                            label: '入荷数',
                            data: [28, 48, 40, 19, 86, 27],
                            backgroundColor: 'rgba(53, 162, 235, 0.5)',
                          },
                        ],
                      }}
                      options={{
                        responsive: true,
                        plugins: {
                          legend: {
                            position: 'top' as const,
                          },
                        },
                      }}
                    />
                  </Box>
                </VStack>
              </CardBody>
            </Card>
          </Container>
        </Box>
      </DashboardLayout>
    </AuthGuard>
  );
} 