'use client';

import { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Container,
  Heading,
  Input,
  Select,
  Stack,
  useToast,
  Skeleton,
  useDisclosure,
  VStack,
  HStack,
} from '@chakra-ui/react';
import { FiPlus, FiDownload, FiUpload, FiHistory, FiEdit } from 'react-icons/fi';
import { useProduct } from '@/hooks/useProduct';
import ProductTable from '@/components/ProductTable';
import ProductFormModal from '@/components/ProductFormModal';
import Pagination from '@/components/Pagination';
import ErrorDisplay from '@/components/ErrorDisplay';
import { SortConfig, SortableField, defaultSortConfig, toggleSortDirection } from '@/utils/sorting';
import { defaultFilterConfig, FilterConfig, clearFilterConfig } from '@/utils/filtering';
import ProductFilterForm from '@/components/ProductFilterForm';
import { ProductImportModal } from '@/components/ProductImportModal';
import StockHistoryModal from '@/components/StockHistoryModal';
import BulkUpdateModal from '@/components/BulkUpdateModal';

/**
 * 商品管理ページ
 */
export default function ProductsPage() {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const {
    isOpen: isImportOpen,
    onOpen: onImportOpen,
    onClose: onImportClose,
  } = useDisclosure();
  const {
    isOpen: isHistoryOpen,
    onOpen: onHistoryOpen,
    onClose: onHistoryClose,
  } = useDisclosure();
  const {
    isOpen: isBulkUpdateOpen,
    onOpen: onBulkUpdateOpen,
    onClose: onBulkUpdateClose,
  } = useDisclosure();
  const [selectedProduct, setSelectedProduct] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize] = useState(10);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState<string>('');
  const [sortConfig, setSortConfig] = useState<SortConfig>(defaultSortConfig);
  const [filterConfig, setFilterConfig] = useState<FilterConfig>(defaultFilterConfig);
  const [selectedProducts, setSelectedProducts] = useState<string[]>([]);
  const [selectedProductId, setSelectedProductId] = useState<string>('');
  const toast = useToast();

  const {
    products,
    totalPages,
    isLoading,
    error,
    fetchProducts,
    createProduct,
    updateProduct,
    deleteProduct,
    bulkDeleteProducts,
    categories,
    exportProducts,
    importProducts,
    bulkUpdateProducts,
  } = useProduct({
    page: currentPage,
    pageSize,
    searchQuery,
    category: selectedCategory,
    sortField: sortConfig.field,
    sortDirection: sortConfig.direction,
    ...filterConfig,
  });

  useEffect(() => {
    fetchProducts();
  }, [currentPage, pageSize, sortConfig, filterConfig]);

  const handleCreate = () => {
    setSelectedProduct(null);
    onOpen();
  };

  const handleEdit = (product) => {
    setSelectedProduct(product);
    onOpen();
  };

  const handleViewHistory = (productId: string) => {
    setSelectedProductId(productId);
    onHistoryOpen();
  };

  const handleSubmit = async (formData) => {
    try {
      if (selectedProduct) {
        await updateProduct(selectedProduct.id, formData);
        toast({
          title: '商品を更新しました',
          status: 'success',
          duration: 3000,
        });
      } else {
        await createProduct(formData);
        toast({
          title: '商品を作成しました',
          status: 'success',
          duration: 3000,
        });
      }
      onClose();
      fetchProducts();
    } catch (error) {
      toast({
        title: 'エラーが発生しました',
        description: error.message,
        status: 'error',
        duration: 5000,
      });
    }
  };

  const handleDelete = async (id) => {
    if (window.confirm('この商品を削除してもよろしいですか？')) {
      try {
        await deleteProduct(id);
        toast({
          title: '商品を削除しました',
          status: 'success',
          duration: 3000,
        });
        fetchProducts();
      } catch (error) {
        toast({
          title: 'エラーが発生しました',
          description: error.message,
          status: 'error',
          duration: 5000,
        });
      }
    }
  };

  const handleBulkDelete = async (ids) => {
    try {
      await bulkDeleteProducts(ids);
      toast({
        title: '商品を一括削除しました',
        status: 'success',
        duration: 3000,
      });
      setSelectedProducts([]);
      fetchProducts();
    } catch (error) {
      toast({
        title: 'エラーが発生しました',
        description: error.message,
        status: 'error',
        duration: 5000,
      });
    }
  };

  const handleBulkUpdate = async (data) => {
    try {
      await bulkUpdateProducts(selectedProducts, data);
      toast({
        title: '商品を一括更新しました',
        status: 'success',
        duration: 3000,
      });
      setSelectedProducts([]);
      fetchProducts();
    } catch (error) {
      toast({
        title: 'エラーが発生しました',
        description: error.message,
        status: 'error',
        duration: 5000,
      });
    }
  };

  const handleSort = (field: string) => {
    setSortConfig((prev) => ({
      field,
      direction:
        prev.field === field && prev.direction === 'asc' ? 'desc' : 'asc',
    }));
    setCurrentPage(1);
  };

  const handleFilterChange = (newConfig: FilterConfig) => {
    setFilterConfig(newConfig);
    setCurrentPage(1);
  };

  const handleClearFilters = () => {
    setFilterConfig(clearFilterConfig());
    setCurrentPage(1);
  };

  const handleSelectProduct = (id: string) => {
    setSelectedProducts((prev) =>
      prev.includes(id)
        ? prev.filter((productId) => productId !== id)
        : [...prev, id]
    );
  };

  const handleSelectAll = (checked: boolean) => {
    setSelectedProducts(checked ? products.map((p) => p.id) : []);
  };

  const handleExport = async () => {
    try {
      await exportProducts();
      toast({
        title: '商品データをエクスポートしました',
        status: 'success',
        duration: 3000,
      });
    } catch (error) {
      toast({
        title: 'エラーが発生しました',
        description: error instanceof Error ? error.message : '予期せぬエラーが発生しました',
        status: 'error',
        duration: 5000,
      });
    }
  };

  const handleImport = async (products: any[]) => {
    try {
      await importProducts(products);
      toast({
        title: '商品データをインポートしました',
        status: 'success',
        duration: 3000,
      });
      onImportClose();
      fetchProducts();
    } catch (error) {
      toast({
        title: 'エラーが発生しました',
        description: error instanceof Error ? error.message : '予期せぬエラーが発生しました',
        status: 'error',
        duration: 5000,
      });
    }
  };

  return (
    <Container maxW="container.xl" py={8}>
      <VStack spacing={8} align="stretch">
        <HStack justify="space-between">
          <Heading size="lg">商品管理</Heading>
          <HStack>
            <Button
              leftIcon={<FiDownload />}
              colorScheme="blue"
              onClick={handleExport}
              isLoading={isLoading}
            >
              エクスポート
            </Button>
            <Button
              leftIcon={<FiUpload />}
              colorScheme="green"
              onClick={onImportOpen}
            >
              インポート
            </Button>
            <Button
              leftIcon={<FiEdit />}
              colorScheme="yellow"
              onClick={onBulkUpdateOpen}
              isDisabled={selectedProducts.length === 0}
            >
              一括更新
            </Button>
            <Button
              leftIcon={<FiPlus />}
              colorScheme="blue"
              onClick={handleCreate}
            >
              新規作成
            </Button>
          </HStack>
        </HStack>

        <ProductFilterForm
          filterConfig={filterConfig}
          categories={categories}
          onFilterChange={handleFilterChange}
          onClear={handleClearFilters}
        />

        <Box overflowX="auto">
          {isLoading ? (
            <Stack spacing={4}>
              <Skeleton height="40px" />
              <Skeleton height="40px" />
              <Skeleton height="40px" />
              <Skeleton height="40px" />
              <Skeleton height="40px" />
            </Stack>
          ) : error ? (
            <ErrorDisplay error={error} onRetry={fetchProducts} />
          ) : (
            <>
              <ProductTable
                products={products}
                isLoading={isLoading}
                onEdit={handleEdit}
                onDelete={handleDelete}
                sortConfig={sortConfig}
                onSort={handleSort}
                selectedProducts={selectedProducts}
                onSelectProduct={handleSelectProduct}
                onSelectAll={handleSelectAll}
                onBulkDelete={handleBulkDelete}
                onViewHistory={handleViewHistory}
              />
              <Pagination
                currentPage={currentPage}
                totalPages={totalPages}
                onPageChange={setCurrentPage}
              />
            </>
          )}
        </Box>

        <ProductFormModal
          isOpen={isOpen}
          onClose={onClose}
          onSubmit={handleSubmit}
          product={selectedProduct}
        />

        <ProductImportModal
          isOpen={isImportOpen}
          onClose={onImportClose}
          onImport={handleImport}
          isLoading={isLoading}
        />

        <StockHistoryModal
          isOpen={isHistoryOpen}
          onClose={onHistoryClose}
          productId={selectedProductId}
        />

        <BulkUpdateModal
          isOpen={isBulkUpdateOpen}
          onClose={onBulkUpdateClose}
          onSubmit={handleBulkUpdate}
          isLoading={isLoading}
          selectedCount={selectedProducts.length}
        />
      </VStack>
    </Container>
  );
} 