'use client';

import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Image,
  Input,
  VStack,
  Text,
  useToast,
  Spinner,
  HStack,
  Badge,
  useDisclosure,
  IconButton,
} from '@chakra-ui/react';
import { useCallback, useState } from 'react';
import { FiUpload, FiX, FiMaximize2 } from 'react-icons/fi';
import { resizeImage } from '@/utils/imageResizer';
import { compressImage, formatFileSize } from '@/utils/imageCompressor';
import ImagePreviewModal from './ImagePreviewModal';

interface ImageUploadProps {
  value?: string;
  onChange: (file: File | undefined) => void;
  label?: string;
  isRequired?: boolean;
  maxWidth?: number;
  maxHeight?: number;
  quality?: number;
  maxSizeKB?: number;
}

/**
 * 画像アップロードコンポーネント
 */
export default function ImageUpload({
  value,
  onChange,
  label = '画像',
  isRequired = false,
  maxWidth = 800,
  maxHeight = 800,
  quality = 0.8,
  maxSizeKB = 500,
}: ImageUploadProps) {
  const [preview, setPreview] = useState<string | undefined>(value);
  const [isProcessing, setIsProcessing] = useState(false);
  const [originalSize, setOriginalSize] = useState<number>(0);
  const [compressedSize, setCompressedSize] = useState<number>(0);
  const [fileName, setFileName] = useState<string>('');
  const toast = useToast();
  const { isOpen, onOpen, onClose } = useDisclosure();

  const handleFileChange = useCallback(
    async (event: React.ChangeEvent<HTMLInputElement>) => {
      const file = event.target.files?.[0];
      if (!file) return;

      try {
        setIsProcessing(true);
        setOriginalSize(file.size);
        setFileName(file.name);

        // ファイル形式のチェック
        if (!file.type.startsWith('image/')) {
          toast({
            title: 'エラー',
            description: '画像ファイルを選択してください',
            status: 'error',
            duration: 3000,
            isClosable: true,
          });
          return;
        }

        // 画像のリサイズ
        const resizedFile = await resizeImage(file, { maxWidth, maxHeight, quality });

        // 画像の圧縮
        const compressedFile = await compressImage(resizedFile, {
          quality,
          maxSizeKB,
          maxIterations: 5,
        });

        setCompressedSize(compressedFile.size);

        // プレビューの設定
        const reader = new FileReader();
        reader.onloadend = () => {
          setPreview(reader.result as string);
          onChange(compressedFile);
          
          const compressionRatio = ((1 - compressedFile.size / file.size) * 100).toFixed(1);
          toast({
            title: '成功',
            description: `画像を最適化しました（${compressionRatio}%削減）`,
            status: 'success',
            duration: 3000,
            isClosable: true,
          });
        };
        reader.readAsDataURL(compressedFile);
      } catch (error) {
        toast({
          title: 'エラー',
          description: '画像の処理中にエラーが発生しました',
          status: 'error',
          duration: 3000,
          isClosable: true,
        });
      } finally {
        setIsProcessing(false);
      }
    },
    [onChange, toast, maxWidth, maxHeight, quality, maxSizeKB]
  );

  const handleRemove = useCallback(() => {
    setPreview(undefined);
    onChange(undefined);
    setOriginalSize(0);
    setCompressedSize(0);
    setFileName('');
  }, [onChange]);

  return (
    <>
      <FormControl isRequired={isRequired}>
        <FormLabel>{label}</FormLabel>
        <VStack spacing={4} align="stretch">
          {preview ? (
            <Box>
              <Box position="relative" width="200px" height="200px">
                <Image
                  src={preview}
                  alt={fileName || 'プレビュー'}
                  objectFit="cover"
                  width="100%"
                  height="100%"
                  borderRadius="md"
                  cursor="pointer"
                  onClick={onOpen}
                  _hover={{ opacity: 0.8 }}
                />
                <HStack position="absolute" top={2} right={2} spacing={2}>
                  <IconButton
                    aria-label="拡大表示"
                    icon={<FiMaximize2 />}
                    size="sm"
                    colorScheme="blue"
                    onClick={onOpen}
                  />
                  <IconButton
                    aria-label="削除"
                    icon={<FiX />}
                    size="sm"
                    colorScheme="red"
                    onClick={handleRemove}
                  />
                </HStack>
              </Box>
              {originalSize > 0 && compressedSize > 0 && (
                <HStack mt={2} spacing={2}>
                  <Badge colorScheme="gray">元のサイズ: {formatFileSize(originalSize)}</Badge>
                  <Badge colorScheme="green">
                    圧縮後: {formatFileSize(compressedSize)}
                    （{((1 - compressedSize / originalSize) * 100).toFixed(1)}%削減）
                  </Badge>
                </HStack>
              )}
            </Box>
          ) : (
            <Box
              border="2px dashed"
              borderColor="gray.300"
              borderRadius="md"
              p={8}
              textAlign="center"
              cursor="pointer"
              _hover={{ borderColor: 'blue.500' }}
              position="relative"
            >
              <Input
                type="file"
                accept="image/*"
                onChange={handleFileChange}
                display="none"
                id="image-upload"
                disabled={isProcessing}
              />
              <label htmlFor="image-upload">
                <VStack spacing={2}>
                  {isProcessing ? (
                    <Spinner size="lg" />
                  ) : (
                    <>
                      <FiUpload size={24} />
                      <Text>画像をアップロード</Text>
                      <Text fontSize="sm" color="gray.500">
                        またはドラッグ＆ドロップ
                      </Text>
                      <Text fontSize="xs" color="gray.400">
                        推奨サイズ: {maxWidth}x{maxHeight}px
                      </Text>
                      <Text fontSize="xs" color="gray.400">
                        最大ファイルサイズ: {formatFileSize(maxSizeKB * 1024)}
                      </Text>
                    </>
                  )}
                </VStack>
              </label>
            </Box>
          )}
        </VStack>
      </FormControl>

      {preview && (
        <ImagePreviewModal
          isOpen={isOpen}
          onClose={onClose}
          imageUrl={preview}
          imageName={fileName}
          originalSize={originalSize}
          compressedSize={compressedSize}
        />
      )}
    </>
  );
} 