'use client';

import {
  Button,
  FormControl,
  FormLabel,
  Input,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  NumberIncrementStepper,
  NumberDecrementStepper,
  Select,
  Textarea,
  VStack,
  useToast,
} from '@chakra-ui/react';
import { useCallback, useState } from 'react';
import { Product, CreateProductParams, UpdateProductParams } from '@/types/product';
import ImageUpload from './ImageUpload';

interface ProductFormProps {
  product?: Product;
  categories: string[];
  onSubmit: (params: CreateProductParams | UpdateProductParams) => Promise<void>;
  onCancel: () => void;
  isLoading?: boolean;
}

/**
 * 商品フォームコンポーネント
 */
export default function ProductForm({
  product,
  categories,
  onSubmit,
  onCancel,
  isLoading = false,
}: ProductFormProps) {
  const [formData, setFormData] = useState<CreateProductParams | UpdateProductParams>({
    name: product?.name || '',
    category: product?.category || '',
    price: product?.price || 0,
    stock: product?.stock || 0,
    unit: product?.unit || '',
    description: product?.description || '',
    imageUrl: product?.imageUrl,
  });
  const [imageFile, setImageFile] = useState<File | undefined>();
  const toast = useToast();

  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      try {
        const params = { ...formData };
        if (imageFile) {
          params.imageFile = imageFile;
        }
        await onSubmit(params);
        toast({
          title: '成功',
          description: product ? '商品を更新しました' : '商品を作成しました',
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
    },
    [formData, imageFile, onSubmit, product, toast]
  );

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
      const { name, value } = e.target;
      setFormData((prev) => ({ ...prev, [name]: value }));
    },
    []
  );

  const handleNumberChange = useCallback(
    (name: string, value: string) => {
      setFormData((prev) => ({ ...prev, [name]: Number(value) }));
    },
    []
  );

  return (
    <form onSubmit={handleSubmit}>
      <VStack spacing={4} align="stretch">
        <FormControl isRequired>
          <FormLabel>商品名</FormLabel>
          <Input
            name="name"
            value={formData.name}
            onChange={handleChange}
            placeholder="商品名を入力"
          />
        </FormControl>

        <FormControl isRequired>
          <FormLabel>カテゴリー</FormLabel>
          <Select name="category" value={formData.category} onChange={handleChange}>
            <option value="">カテゴリーを選択</option>
            {categories.map((category) => (
              <option key={category} value={category}>
                {category}
              </option>
            ))}
          </Select>
        </FormControl>

        <FormControl isRequired>
          <FormLabel>価格</FormLabel>
          <NumberInput
            name="price"
            value={formData.price}
            onChange={(value) => handleNumberChange('price', value)}
            min={0}
            precision={0}
          >
            <NumberInputField />
            <NumberInputStepper>
              <NumberIncrementStepper />
              <NumberDecrementStepper />
            </NumberInputStepper>
          </NumberInput>
        </FormControl>

        <FormControl isRequired>
          <FormLabel>在庫数</FormLabel>
          <NumberInput
            name="stock"
            value={formData.stock}
            onChange={(value) => handleNumberChange('stock', value)}
            min={0}
            precision={0}
          >
            <NumberInputField />
            <NumberInputStepper>
              <NumberIncrementStepper />
              <NumberDecrementStepper />
            </NumberInputStepper>
          </NumberInput>
        </FormControl>

        <FormControl isRequired>
          <FormLabel>単位</FormLabel>
          <Input
            name="unit"
            value={formData.unit}
            onChange={handleChange}
            placeholder="個、kg、g など"
          />
        </FormControl>

        <FormControl>
          <FormLabel>説明</FormLabel>
          <Textarea
            name="description"
            value={formData.description}
            onChange={handleChange}
            placeholder="商品の説明を入力"
            rows={4}
          />
        </FormControl>

        <ImageUpload
          value={formData.imageUrl}
          onChange={setImageFile}
          label="商品画像"
        />

        <Button
          type="submit"
          colorScheme="blue"
          isLoading={isLoading}
          loadingText={product ? '更新中...' : '作成中...'}
        >
          {product ? '更新' : '作成'}
        </Button>

        <Button onClick={onCancel}>キャンセル</Button>
      </VStack>
    </form>
  );
} 