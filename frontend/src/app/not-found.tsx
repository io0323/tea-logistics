'use client';

import {
  Box,
  Button,
  Container,
  Heading,
  Text,
  VStack,
} from '@chakra-ui/react';
import { useRouter } from 'next/navigation';

/**
 * 404ページコンポーネント
 */
export default function NotFound() {
  const router = useRouter();

  return (
    <Container maxW="container.md" py={8}>
      <Box
        p={8}
        borderWidth={1}
        borderRadius="lg"
        boxShadow="md"
      >
        <VStack spacing={6} align="center">
          <Heading size="lg">ページが見つかりません</Heading>
          <Text>お探しのページは存在しないか、移動または削除された可能性があります。</Text>
          <Button colorScheme="blue" onClick={() => router.push('/')}>
            ホームに戻る
          </Button>
        </VStack>
      </Box>
    </Container>
  );
} 