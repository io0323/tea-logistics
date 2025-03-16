'use client';

import { Box, Container, useToast } from '@chakra-ui/react';
import BulkImageUpload from '@/components/BulkImageUpload';

export default function TestPage() {
  const toast = useToast();

  const handleUpload = (files: File[]) => {
    console.log('アップロードされたファイル:', files);
    toast({
      title: '成功',
      description: `${files.length}枚の画像がアップロードされました`,
      status: 'success',
      duration: 3000,
      isClosable: true,
    });
  };

  return (
    <Container maxW="container.xl" py={8}>
      <Box bg="white" p={6} borderRadius="md" shadow="md">
        <BulkImageUpload
          onUpload={handleUpload}
          maxFiles={5}
          maxWidth={1200}
          maxHeight={1200}
          quality={0.8}
          maxSizeKB={1000}
        />
      </Box>
    </Container>
  );
} 