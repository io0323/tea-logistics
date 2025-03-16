'use client';

import {
  HStack,
  Button,
  Text,
  IconButton,
  useColorModeValue,
} from '@chakra-ui/react';
import { FiChevronLeft, FiChevronRight } from 'react-icons/fi';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

/**
 * ページネーションコンポーネント
 */
export default function Pagination({
  currentPage,
  totalPages,
  onPageChange,
}: PaginationProps) {
  const buttonBg = useColorModeValue('white', 'gray.700');
  const buttonHoverBg = useColorModeValue('gray.100', 'gray.600');

  const handlePrevious = () => {
    if (currentPage > 1) {
      onPageChange(currentPage - 1);
    }
  };

  const handleNext = () => {
    if (currentPage < totalPages) {
      onPageChange(currentPage + 1);
    }
  };

  const renderPageNumbers = () => {
    const pages = [];
    const maxVisiblePages = 5;
    let startPage = Math.max(1, currentPage - Math.floor(maxVisiblePages / 2));
    let endPage = Math.min(totalPages, startPage + maxVisiblePages - 1);

    if (endPage - startPage + 1 < maxVisiblePages) {
      startPage = Math.max(1, endPage - maxVisiblePages + 1);
    }

    for (let i = startPage; i <= endPage; i++) {
      pages.push(
        <Button
          key={i}
          size="sm"
          variant={currentPage === i ? 'solid' : 'outline'}
          colorScheme={currentPage === i ? 'blue' : 'gray'}
          onClick={() => onPageChange(i)}
          bg={currentPage === i ? 'blue.500' : buttonBg}
          _hover={{
            bg: currentPage === i ? 'blue.600' : buttonHoverBg,
          }}
        >
          {i}
        </Button>
      );
    }

    return pages;
  };

  return (
    <HStack spacing={2} justify="center" mt={4}>
      <IconButton
        aria-label="前のページ"
        icon={<FiChevronLeft />}
        onClick={handlePrevious}
        isDisabled={currentPage === 1}
        size="sm"
      />
      {renderPageNumbers()}
      <IconButton
        aria-label="次のページ"
        icon={<FiChevronRight />}
        onClick={handleNext}
        isDisabled={currentPage === totalPages}
        size="sm"
      />
      <Text fontSize="sm" color="gray.500">
        {currentPage} / {totalPages}
      </Text>
    </HStack>
  );
} 