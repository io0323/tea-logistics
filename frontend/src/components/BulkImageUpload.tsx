'use client';

import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Grid,
  VStack,
  Text,
  useToast,
  Spinner,
  Progress,
  HStack,
  IconButton,
  Badge,
} from '@chakra-ui/react';
import { useCallback, useState } from 'react';
import { FiUpload, FiX, FiMaximize2 } from 'react-icons/fi';
import { resizeImage } from '@/utils/imageResizer';
import { compressImage, formatFileSize } from '@/utils/imageCompressor';
import ImagePreviewModal from './ImagePreviewModal';

interface ProcessedImage {
  id: string;
  file: File;
  preview: string;
  originalSize: number;
  compressedSize: number;
}

interface BulkImageUploadProps {
  onUpload: (files: File[]) => void;
  maxFiles?: number;
  maxWidth?: number;
  maxHeight?: number;
  quality?: number;
  maxSizeKB?: number;
}

/**
 * 一括画像アップロードコンポーネント
 */
export default function BulkImageUpload({
  onUpload,
  maxFiles = 10,
  maxWidth = 800,
  maxHeight = 800,
  quality = 0.8,
  maxSizeKB = 500,
}: BulkImageUploadProps) {
  const [images, setImages] = useState<ProcessedImage[]>([]);
  const [isProcessing, setIsProcessing] = useState(false);
  const [progress, setProgress] = useState(0);
  const [previewImage, setPreviewImage] = useState<ProcessedImage | null>(null);
  const toast = useToast();

  const processImage = async (file: File): Promise<ProcessedImage> => {
    // 画像のリサイズ
    const resizedFile = await resizeImage(file, { maxWidth, maxHeight, quality });

    // 画像の圧縮
    const compressedFile = await compressImage(resizedFile, {
      quality,
      maxSizeKB,
      maxIterations: 5,
    });

    // プレビューの生成
    const preview = await new Promise<string>((resolve) => {
      const reader = new FileReader();
      reader.onloadend = () => resolve(reader.result as string);
      reader.readAsDataURL(compressedFile);
    });

    return {
      id: Math.random().toString(36).substring(2),
      file: compressedFile,
      preview,
      originalSize: file.size,
      compressedSize: compressedFile.size,
    };
  };

  const handleFileChange = useCallback(
    async (event: React.ChangeEvent<HTMLInputElement>) => {
      const files = Array.from(event.target.files || []);
      if (files.length === 0) return;

      try {
        setIsProcessing(true);
        setProgress(0);

        // ファイル数の制限チェック
        if (images.length + files.length > maxFiles) {
          toast({
            title: 'エラー',
            description: `アップロードできる画像は最大${maxFiles}枚です`,
            status: 'error',
            duration: 3000,
            isClosable: true,
          });
          return;
        }

        const processedImages: ProcessedImage[] = [];
        let processedCount = 0;

        for (const file of files) {
          // ファイル形式のチェック
          if (!file.type.startsWith('image/')) {
            toast({
              title: '警告',
              description: `${file.name}は画像ファイルではありません`,
              status: 'warning',
              duration: 3000,
              isClosable: true,
            });
            continue;
          }

          try {
            const processedImage = await processImage(file);
            processedImages.push(processedImage);
          } catch (error) {
            toast({
              title: '警告',
              description: `${file.name}の処理中にエラーが発生しました`,
              status: 'warning',
              duration: 3000,
              isClosable: true,
            });
          }

          processedCount++;
          setProgress((processedCount / files.length) * 100);
        }

        setImages((prev) => [...prev, ...processedImages]);
        onUpload(processedImages.map((img) => img.file));

        toast({
          title: '成功',
          description: `${processedImages.length}枚の画像をアップロードしました`,
          status: 'success',
          duration: 3000,
          isClosable: true,
        });
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
        setProgress(0);
      }
    },
    [images, maxFiles, onUpload, toast, maxWidth, maxHeight, quality, maxSizeKB]
  );

  const handleRemove = useCallback((id: string) => {
    setImages((prev) => prev.filter((img) => img.id !== id));
  }, []);

  const handlePreview = useCallback((image: ProcessedImage) => {
    setPreviewImage(image);
  }, []);

  return (
    <>
      <FormControl>
        <FormLabel>画像一括アップロード</FormLabel>
        <VStack spacing={4} align="stretch">
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
            <input
              type="file"
              accept="image/*"
              onChange={handleFileChange}
              multiple
              style={{ display: 'none' }}
              id="bulk-image-upload"
              disabled={isProcessing}
            />
            <label htmlFor="bulk-image-upload">
              <VStack spacing={2}>
                {isProcessing ? (
                  <>
                    <Spinner size="lg" />
                    <Progress
                      value={progress}
                      width="100%"
                      size="sm"
                      colorScheme="blue"
                      hasStripe
                      isAnimated
                    />
                  </>
                ) : (
                  <>
                    <FiUpload size={24} />
                    <Text>画像をアップロード</Text>
                    <Text fontSize="sm" color="gray.500">
                      またはドラッグ＆ドロップ
                    </Text>
                    <Text fontSize="xs" color="gray.400">
                      最大{maxFiles}枚まで
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

          {images.length > 0 && (
            <Grid templateColumns="repeat(auto-fill, minmax(200px, 1fr))" gap={4}>
              {images.map((image) => (
                <Box key={image.id} position="relative">
                  <Box
                    width="100%"
                    paddingTop="100%"
                    position="relative"
                    overflow="hidden"
                    borderRadius="md"
                  >
                    <Box
                      as="img"
                      src={image.preview}
                      alt="プレビュー"
                      position="absolute"
                      top="0"
                      left="0"
                      width="100%"
                      height="100%"
                      objectFit="cover"
                      cursor="pointer"
                      onClick={() => handlePreview(image)}
                      _hover={{ opacity: 0.8 }}
                    />
                    <HStack position="absolute" top={2} right={2} spacing={2}>
                      <IconButton
                        aria-label="拡大表示"
                        icon={<FiMaximize2 />}
                        size="sm"
                        colorScheme="blue"
                        onClick={() => handlePreview(image)}
                      />
                      <IconButton
                        aria-label="削除"
                        icon={<FiX />}
                        size="sm"
                        colorScheme="red"
                        onClick={() => handleRemove(image.id)}
                      />
                    </HStack>
                  </Box>
                  <VStack mt={2} spacing={1}>
                    <Badge colorScheme="gray" fontSize="xs">
                      元のサイズ: {formatFileSize(image.originalSize)}
                    </Badge>
                    <Badge colorScheme="green" fontSize="xs">
                      圧縮後: {formatFileSize(image.compressedSize)}
                      （{((1 - image.compressedSize / image.originalSize) * 100).toFixed(1)}%削減）
                    </Badge>
                  </VStack>
                </Box>
              ))}
            </Grid>
          )}
        </VStack>
      </FormControl>

      {previewImage && (
        <ImagePreviewModal
          isOpen={!!previewImage}
          onClose={() => setPreviewImage(null)}
          imageUrl={previewImage.preview}
          originalSize={previewImage.originalSize}
          compressedSize={previewImage.compressedSize}
        />
      )}
    </>
  );
} 