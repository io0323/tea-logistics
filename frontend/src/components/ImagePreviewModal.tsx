'use client';

import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalCloseButton,
  Image,
  Text,
  VStack,
  Badge,
  HStack,
} from '@chakra-ui/react';
import { formatFileSize } from '@/utils/imageCompressor';

interface ImagePreviewModalProps {
  isOpen: boolean;
  onClose: () => void;
  imageUrl: string;
  imageName?: string;
  originalSize?: number;
  compressedSize?: number;
}

/**
 * 画像プレビューモーダルコンポーネント
 */
export default function ImagePreviewModal({
  isOpen,
  onClose,
  imageUrl,
  imageName,
  originalSize,
  compressedSize,
}: ImagePreviewModalProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} size="xl" isCentered>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>
          {imageName || '画像プレビュー'}
        </ModalHeader>
        <ModalCloseButton />
        <ModalBody pb={6}>
          <VStack spacing={4} align="stretch">
            <Image
              src={imageUrl}
              alt={imageName || '画像プレビュー'}
              width="100%"
              height="auto"
              maxH="70vh"
              objectFit="contain"
              borderRadius="md"
            />
            {originalSize && compressedSize && (
              <HStack spacing={2} justify="center">
                <Badge colorScheme="gray">元のサイズ: {formatFileSize(originalSize)}</Badge>
                <Badge colorScheme="green">
                  圧縮後: {formatFileSize(compressedSize)}
                  （{((1 - compressedSize / originalSize) * 100).toFixed(1)}%削減）
                </Badge>
              </HStack>
            )}
            <Text fontSize="sm" color="gray.500" textAlign="center">
              クリックして閉じる
            </Text>
          </VStack>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
} 