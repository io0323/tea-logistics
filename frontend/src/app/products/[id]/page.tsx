'use client';

import { useEffect, useState } from 'react';
import {
  Box,
  Container,
  Heading,
  Text,
  VStack,
  HStack,
  Badge,
  Button,
  useToast,
  Skeleton,
  Divider,
  useColorModeValue,
} from '@chakra-ui/react';
import { FiArrowLeft, FiEdit2 } from 'react-icons/fi';
import { useRouter } from 'next/navigation';
import { useProduct } from '@/hooks/useProduct';
import ProductFormModal from '@/components/ProductFormModal';
import ErrorDisplay from '@/components/ErrorDisplay';

interface ProductDetailPageProps {
  params: {
    id: string;
  };
}

/**
 * 商品詳細ページ
 */
export default function ProductDetailPage({ params }: ProductDetailPageProps) {
  const router = useRouter();
  const toast = useToast();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [selectedProduct, setSelectedProduct] = useState(null);
  const bgColor = useColorModeValue('white', 'gray.700');
  const borderColor = useColorModeValue('gray.200', 'gray.600');

  const {
    product,
    isLoading,
    error,
    fetchProduct,
    updateProduct,
  } = useProduct();

  useEffect(() => {
    fetchProduct(params.id);
  }, [params.id]);

  const handleEdit = () => {
    setSelectedProduct(product);
    onOpen();
  };

  const handleSubmit = async (formData) => {
    try {
      await updateProduct(params.id, formData);
      toast({
        title: '商品を更新しました',
        status: 'success',
        duration: 3000,
      });
      onClose();
      fetchProduct(params.id);
    } catch (error) {
      toast({
        title: 'エラーが発生しました',
        description: error.message,
        status: 'error',
        duration: 5000,
      });
    }
  };

  if (error) {
    return <ErrorDisplay error={error} onRetry={() => fetchProduct(params.id)} />;
  }

  return (
    <Container maxW="container.xl" py={8}>
      <VStack spacing={8} align="stretch">
        <HStack justify="space-between">
          <HStack spacing={4}>
            <Button
              leftIcon={<FiArrowLeft />}
              variant="ghost"
              onClick={() => router.push('/products')}
            >
              戻る
            </Button>
            <Heading size="lg">商品詳細</Heading>
          </HStack>
          <Button
            leftIcon={<FiEdit2 />}
            colorScheme="blue"
            onClick={handleEdit}
          >
            編集
          </Button>
        </HStack>

        <Box
          p={6}
          bg={bgColor}
          borderWidth={1}
          borderColor={borderColor}
          borderRadius="md"
        >
          {isLoading ? (
            <VStack spacing={4} align="stretch">
              <Skeleton height="40px" />
              <Skeleton height="20px" />
              <Divider />
              <Skeleton height="20px" />
              <Skeleton height="20px" />
            </VStack>
          ) : product ? (
            <VStack spacing={6} align="stretch">
              <HStack justify="space-between">
                <Heading size="md">{product.name}</Heading>
                <Badge colorScheme="blue">{product.category}</Badge>
              </HStack>

              <Text color="gray.600">{product.description}</Text>

              <Divider />

              <VStack spacing={4} align="stretch">
                <HStack justify="space-between">
                  <Text fontWeight="bold">単価</Text>
                  <Text>{product.price.toLocaleString()}円</Text>
                </HStack>
                <HStack justify="space-between">
                  <Text fontWeight="bold">単位</Text>
                  <Text>{product.unit}</Text>
                </HStack>
                <HStack justify="space-between">
                  <Text fontWeight="bold">在庫数</Text>
                  <Text>{product.stock.toLocaleString()}</Text>
                </HStack>
                <HStack justify="space-between">
                  <Text fontWeight="bold">作成日</Text>
                  <Text>
                    {new Date(product.createdAt).toLocaleDateString()}
                  </Text>
                </HStack>
                <HStack justify="space-between">
                  <Text fontWeight="bold">更新日</Text>
                  <Text>
                    {new Date(product.updatedAt).toLocaleDateString()}
                  </Text>
                </HStack>
              </VStack>
            </VStack>
          ) : null}
        </Box>

        <ProductFormModal
          isOpen={isOpen}
          onClose={onClose}
          onSubmit={handleSubmit}
          product={selectedProduct}
        />
      </VStack>
    </Container>
  );
} 