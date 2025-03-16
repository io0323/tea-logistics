'use client';

import { Box, Button, Container, Heading, Text, VStack } from '@chakra-ui/react';

/**
 * グローバルエラーページコンポーネント
 */
export default function Error({
  error,
  reset,
}: {
  error: Error;
  reset: () => void;
}) {
  return (
    <Container maxW="container.md" py={8}>
      <Box p={8} borderWidth={1} borderRadius="lg" boxShadow="md">
        <VStack spacing={4} align="stretch">
          <Heading size="lg" color="red.500">エラーが発生しました</Heading>
          <Text>{error.message}</Text>
          <Button onClick={reset} colorScheme="blue">
            再試行
          </Button>
        </VStack>
      </Box>
    </Container>
  );
} 