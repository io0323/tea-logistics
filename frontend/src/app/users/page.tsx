'use client';

import {
  Box,
  Container,
  Heading,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Button,
  Badge,
  useColorModeValue,
} from '@chakra-ui/react';
import Navigation from '@/components/Navigation';

export default function UsersPage() {
  return (
    <Box minH="100vh" bg={useColorModeValue('gray.50', 'gray.900')}>
      <Navigation />
      <Container maxW="container.xl" py={8}>
        <Heading mb={6}>ユーザー管理</Heading>
        <Box bg="white" rounded="lg" shadow="md" overflow="hidden">
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>ID</Th>
                <Th>名前</Th>
                <Th>メールアドレス</Th>
                <Th>役割</Th>
                <Th>状態</Th>
                <Th>操作</Th>
              </Tr>
            </Thead>
            <Tbody>
              <Tr>
                <Td>1</Td>
                <Td>管理者</Td>
                <Td>admin@example.com</Td>
                <Td>
                  <Badge colorScheme="red">管理者</Badge>
                </Td>
                <Td>
                  <Badge colorScheme="green">有効</Badge>
                </Td>
                <Td>
                  <Button size="sm" colorScheme="blue" mr={2}>
                    編集
                  </Button>
                </Td>
              </Tr>
              <Tr>
                <Td>2</Td>
                <Td>一般ユーザー</Td>
                <Td>user@example.com</Td>
                <Td>
                  <Badge colorScheme="blue">一般</Badge>
                </Td>
                <Td>
                  <Badge colorScheme="green">有効</Badge>
                </Td>
                <Td>
                  <Button size="sm" colorScheme="blue" mr={2}>
                    編集
                  </Button>
                </Td>
              </Tr>
            </Tbody>
          </Table>
        </Box>
      </Container>
    </Box>
  );
} 