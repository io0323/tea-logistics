'use client';

import {
  Box,
  Button,
  Heading,
  Text,
  VStack,
  Icon,
  useColorModeValue,
} from '@chakra-ui/react';
import { FiAlertCircle, FiRefreshCw } from 'react-icons/fi';

interface ErrorDisplayProps {
  error: Error;
  onRetry?: () => void;
}

/**
 * エラー表示コンポーネント
 */
export default function ErrorDisplay({ error, onRetry }: ErrorDisplayProps) {
  const bgColor = useColorModeValue('red.50', 'red.900');
  const borderColor = useColorModeValue('red.200', 'red.700');
  const textColor = useColorModeValue('red.700', 'red.200');

  return (
    <Box
      p={6}
      bg={bgColor}
      borderWidth={1}
      borderColor={borderColor}
      borderRadius="md"
      textAlign="center"
    >
      <VStack spacing={4}>
        <Icon as={FiAlertCircle} w={8} h={8} color="red.500" />
        <Heading size="md" color={textColor}>
          エラーが発生しました
        </Heading>
        <Text color={textColor}>{error.message}</Text>
        {onRetry && (
          <Button
            leftIcon={<FiRefreshCw />}
            colorScheme="red"
            variant="outline"
            onClick={onRetry}
          >
            再試行
          </Button>
        )}
      </VStack>
    </Box>
  );
} 