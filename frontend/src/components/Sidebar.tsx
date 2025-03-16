'use client';

import {
  Box,
  CloseButton,
  Flex,
  useColorModeValue,
  Text,
  BoxProps,
  Link,
} from '@chakra-ui/react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/hooks/useAuth';

interface SidebarProps extends BoxProps {
  onClose: () => void;
}

/**
 * サイドバーコンポーネント
 */
export default function Sidebar({ onClose, ...rest }: SidebarProps) {
  const router = useRouter();
  const { user } = useAuth();

  const LinkItems = [
    { name: 'ダッシュボード', href: '/dashboard' },
    { name: '在庫管理', href: '/inventory' },
    { name: '商品管理', href: '/products' },
    { name: '発注管理', href: '/orders' },
    { name: '設定', href: '/settings' },
  ];

  return (
    <Box
      transition="3s ease"
      bg={useColorModeValue('white', 'gray.900')}
      borderRight="1px"
      borderRightColor={useColorModeValue('gray.200', 'gray.700')}
      w={{ base: 'full', md: 60 }}
      pos="fixed"
      h="full"
      {...rest}
    >
      <Flex h="20" alignItems="center" mx="8" justifyContent="space-between">
        <Text fontSize="2xl" fontFamily="monospace" fontWeight="bold">
          茶葉物流管理
        </Text>
        <CloseButton display={{ base: 'flex', md: 'none' }} onClick={onClose} />
      </Flex>
      {LinkItems.map((link) => (
        <Link
          key={link.name}
          onClick={() => router.push(link.href)}
          style={{ textDecoration: 'none' }}
          _focus={{ boxShadow: 'none' }}
        >
          <Box
            p={4}
            mx={4}
            borderRadius="lg"
            role="group"
            cursor="pointer"
            _hover={{
              bg: 'blue.400',
              color: 'white',
            }}
          >
            {link.name}
          </Box>
        </Link>
      ))}
    </Box>
  );
} 