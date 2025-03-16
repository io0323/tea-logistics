'use client';

import {
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  IconButton,
  Checkbox,
  Button,
  HStack,
  Text,
  Link,
  useToast,
  Box,
  useColorModeValue,
  Skeleton,
  Image,
  useDisclosure,
} from '@chakra-ui/react';
import { FiEdit2, FiTrash2, FiEye } from 'react-icons/fi';
import { useRouter } from 'next/navigation';
import { SortConfig } from '@/utils/sorting';
import { Product } from '@/types/product';
import StockHistoryModal from './StockHistoryModal';

interface ProductTableProps {
  products: Product[];
  isLoading: boolean;
  onEdit: (product: Product) => void;
  onDelete: (id: string) => void;
  onBulkDelete: (ids: string[]) => void;
  sortConfig: SortConfig;
  onSort: (field: string) => void;
  selectedProducts: string[];
  onSelectProduct: (id: string) => void;
  onSelectAll: (checked: boolean) => void;
  onViewHistory: (id: string) => void;
}

/**
 * 商品テーブルコンポーネント
 */
export default function ProductTable({
  products,
  isLoading,
  onEdit,
  onDelete,
  onBulkDelete,
  sortConfig,
  onSort,
  selectedProducts,
  onSelectProduct,
  onSelectAll,
  onViewHistory,
}: ProductTableProps) {
  const router = useRouter();
  const toast = useToast();
  const { isOpen, onOpen, onClose } = useDisclosure();
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const hoverBgColor = useColorModeValue('gray.50', 'gray.700');

  const renderSortIcon = (field: string) => {
    if (sortConfig.field !== field) return null;
    return sortConfig.direction === 'asc' ? '↑' : '↓';
  };

  const handleViewDetails = useCallback(
    (product: Product) => {
      router.push(`/products/${product.id}`);
    },
    [router]
  );

  const handleViewHistory = useCallback(
    (product: Product) => {
      setSelectedProduct(product);
      onOpen();
    },
    [onOpen]
  );

  const handleDelete = useCallback(
    async (id: string) => {
      if (window.confirm('この商品を削除してもよろしいですか？')) {
        try {
          await onDelete(id);
          toast({
            title: '成功',
            description: '商品を削除しました',
            status: 'success',
            duration: 3000,
            isClosable: true,
          });
        } catch (error) {
          toast({
            title: 'エラー',
            description: error instanceof Error ? error.message : '予期せぬエラーが発生しました',
            status: 'error',
            duration: 3000,
            isClosable: true,
          });
        }
      }
    },
    [onDelete, toast]
  );

  const handleBulkDelete = useCallback(
    async () => {
      if (window.confirm(`${selectedProducts.length}件の商品を削除してもよろしいですか？`)) {
        try {
          await onBulkDelete(selectedProducts);
          toast({
            title: '成功',
            description: '選択した商品を削除しました',
            status: 'success',
            duration: 3000,
            isClosable: true,
          });
        } catch (error) {
          toast({
            title: 'エラー',
            description: error instanceof Error ? error.message : '予期せぬエラーが発生しました',
            status: 'error',
            duration: 3000,
            isClosable: true,
          });
        }
      }
    },
    [selectedProducts, onBulkDelete, toast]
  );

  if (isLoading) {
    return (
      <Box>
        <Skeleton height="40px" mb={2} />
        <Skeleton height="40px" mb={2} />
        <Skeleton height="40px" mb={2} />
        <Skeleton height="40px" mb={2} />
        <Skeleton height="40px" />
      </Box>
    );
  }

  return (
    <>
      <Box>
        {selectedProducts.length > 0 && (
          <HStack spacing={4} mb={4}>
            <Text>{selectedProducts.length}件選択中</Text>
            <Button
              size="sm"
              colorScheme="red"
              leftIcon={<FiTrash2 />}
              onClick={handleBulkDelete}
            >
              一括削除
            </Button>
          </HStack>
        )}
        <Table variant="simple">
          <Thead>
            <Tr>
              <Th width="50px">
                <Checkbox
                  isChecked={selectedProducts.length === products.length}
                  onChange={(e) => onSelectAll(e.target.checked)}
                />
              </Th>
              <Th width="100px">画像</Th>
              <Th
                cursor="pointer"
                onClick={() => onSort('name')}
                _hover={{ bg: hoverBgColor }}
              >
                商品名 {renderSortIcon('name')}
              </Th>
              <Th
                cursor="pointer"
                onClick={() => onSort('category')}
                _hover={{ bg: hoverBgColor }}
              >
                カテゴリー {renderSortIcon('category')}
              </Th>
              <Th
                cursor="pointer"
                onClick={() => onSort('price')}
                _hover={{ bg: hoverBgColor }}
              >
                価格 {renderSortIcon('price')}
              </Th>
              <Th
                cursor="pointer"
                onClick={() => onSort('stock')}
                _hover={{ bg: hoverBgColor }}
              >
                在庫 {renderSortIcon('stock')}
              </Th>
              <Th
                cursor="pointer"
                onClick={() => onSort('createdAt')}
                _hover={{ bg: hoverBgColor }}
              >
                作成日 {renderSortIcon('createdAt')}
              </Th>
              <Th>操作</Th>
            </Tr>
          </Thead>
          <Tbody>
            {products.map((product) => (
              <Tr key={product.id}>
                <Td>
                  <Checkbox
                    isChecked={selectedProducts.includes(product.id)}
                    onChange={() => onSelectProduct(product.id)}
                  />
                </Td>
                <Td>
                  {product.imageUrl && (
                    <Image
                      src={product.imageUrl}
                      alt={product.name}
                      width="60px"
                      height="60px"
                      objectFit="cover"
                      borderRadius="md"
                    />
                  )}
                </Td>
                <Td>
                  <Link
                    color="blue.500"
                    onClick={() => handleViewDetails(product)}
                    _hover={{ textDecoration: 'none' }}
                  >
                    {product.name}
                  </Link>
                </Td>
                <Td>{product.category}</Td>
                <Td>{product.price.toLocaleString()}円</Td>
                <Td>{product.stock.toLocaleString()}</Td>
                <Td>
                  {new Date(product.createdAt).toLocaleDateString()}
                </Td>
                <Td>
                  <HStack spacing={2}>
                    <IconButton
                      aria-label="詳細表示"
                      icon={<FiEye />}
                      size="sm"
                      onClick={() => handleViewDetails(product)}
                    />
                    <IconButton
                      aria-label="編集"
                      icon={<FiEdit2 />}
                      size="sm"
                      onClick={() => onEdit(product)}
                    />
                    <IconButton
                      aria-label="削除"
                      icon={<FiTrash2 />}
                      size="sm"
                      colorScheme="red"
                      onClick={() => handleDelete(product.id)}
                    />
                  </HStack>
                </Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </Box>

      {selectedProduct && (
        <StockHistoryModal
          isOpen={isOpen}
          onClose={onClose}
          productId={selectedProduct.id}
        />
      )}
    </>
  );
} 