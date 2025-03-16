'use client';

import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  Select,
  HStack,
  VStack,
  useColorModeValue,
} from '@chakra-ui/react';
import { FiFilter, FiX } from 'react-icons/fi';
import { FilterConfig, PriceRange } from '@/utils/filtering';

interface ProductFilterFormProps {
  filterConfig: FilterConfig;
  categories: string[];
  onFilterChange: (config: FilterConfig) => void;
  onClear: () => void;
}

/**
 * 商品フィルターフォームコンポーネント
 */
export default function ProductFilterForm({
  filterConfig,
  categories,
  onFilterChange,
  onClear,
}: ProductFilterFormProps) {
  const bgColor = useColorModeValue('white', 'gray.700');
  const borderColor = useColorModeValue('gray.200', 'gray.600');

  const handlePriceRangeChange = (
    field: keyof PriceRange,
    value: string
  ) => {
    const numValue = value ? parseInt(value) : 0;
    onFilterChange({
      ...filterConfig,
      priceRange: {
        ...filterConfig.priceRange,
        [field]: numValue,
      },
    });
  };

  return (
    <Box
      p={4}
      bg={bgColor}
      borderWidth={1}
      borderColor={borderColor}
      borderRadius="md"
    >
      <VStack spacing={4} align="stretch">
        <HStack justify="space-between">
          <FormLabel mb={0}>フィルター</FormLabel>
          <Button
            size="sm"
            variant="ghost"
            leftIcon={<FiX />}
            onClick={onClear}
          >
            クリア
          </Button>
        </HStack>

        <FormControl>
          <FormLabel>カテゴリー</FormLabel>
          <Select
            value={filterConfig.category}
            onChange={(e) =>
              onFilterChange({
                ...filterConfig,
                category: e.target.value,
              })
            }
            placeholder="カテゴリーを選択"
          >
            {categories.map((category) => (
              <option key={category} value={category}>
                {category}
              </option>
            ))}
          </Select>
        </FormControl>

        <FormControl>
          <FormLabel>価格範囲</FormLabel>
          <HStack spacing={4}>
            <Input
              type="number"
              value={filterConfig.priceRange.min}
              onChange={(e) =>
                handlePriceRangeChange('min', e.target.value)
              }
              placeholder="最小価格"
            />
            <Input
              type="number"
              value={filterConfig.priceRange.max}
              onChange={(e) =>
                handlePriceRangeChange('max', e.target.value)
              }
              placeholder="最大価格"
            />
          </HStack>
        </FormControl>

        <FormControl>
          <FormLabel>キーワード検索</FormLabel>
          <Input
            value={filterConfig.searchQuery}
            onChange={(e) =>
              onFilterChange({
                ...filterConfig,
                searchQuery: e.target.value,
              })
            }
            placeholder="商品名で検索"
          />
        </FormControl>
      </VStack>
    </Box>
  );
} 