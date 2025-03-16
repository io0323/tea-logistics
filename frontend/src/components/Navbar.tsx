'use client';

import {
  IconButton,
  Box,
  Flex,
  HStack,
  VStack,
  useColorModeValue,
  Text,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  FlexProps,
  Avatar,
} from '@chakra-ui/react';
import { HamburgerIcon } from '@chakra-ui/icons';
import { useAuth } from '@/hooks/useAuth';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import { FiHome, FiBox, FiTruck, FiDownload, FiBarChart2, FiSettings } from 'react-icons/fi';
import NextLink from 'next/link';

interface NavbarProps extends FlexProps {
  onOpen: () => void;
}

/**
 * ナビゲーションバーコンポーネント
 * 全画面共通のヘッダーUIを提供します
 */
export default function Navbar() {
  const { user, logout } = useAuth();
  const pathname = usePathname();
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');
  const selectedBg = useColorModeValue('green.500', 'green.200');
  const selectedColor = useColorModeValue('white', 'gray.800');

  const menuItems = [
    { href: '/dashboard', label: 'ダッシュボード', icon: FiHome },
    { href: '/inventory', label: '在庫管理', icon: FiBox },
    { href: '/shipping', label: '出荷管理', icon: FiTruck },
    { href: '/receiving', label: '入荷管理', icon: FiDownload },
    { href: '/reports', label: 'レポート', icon: FiBarChart2 },
    { href: '/settings', label: '設定', icon: FiSettings },
  ];

  const handleLogout = () => {
    logout();
  };

  return (
    <Box bg={bgColor} borderBottom="1px" borderColor={borderColor} position="fixed" w="100%" zIndex={1}>
      <Flex h={16} alignItems="center" justifyContent="space-between" maxW="7xl" mx="auto" px={4}>
        <HStack spacing={8} alignItems="center">
          <Text fontSize="lg" fontWeight="bold" fontFamily="'Noto Sans JP', sans-serif">
            物流管理システム
          </Text>
          {menuItems.map((item) => (
            <NextLink key={item.href} href={item.href} passHref style={{ textDecoration: 'none' }}>
              <Box
                px={3}
                py={2}
                rounded="md"
                bg={pathname === item.href ? selectedBg : 'transparent'}
                color={pathname === item.href ? selectedColor : undefined}
                _hover={{
                  bg: pathname === item.href ? selectedBg : useColorModeValue('gray.100', 'gray.700'),
                  color: pathname === item.href ? selectedColor : undefined,
                }}
                cursor="pointer"
                display="flex"
                alignItems="center"
              >
                <Box as={item.icon} mr={2} />
                <Text>{item.label}</Text>
              </Box>
            </NextLink>
          ))}
        </HStack>

        <Menu>
          <MenuButton
            as={Box}
            rounded="full"
            cursor="pointer"
            display="flex"
            alignItems="center"
          >
            <HStack spacing={2}>
              <Avatar
                size="sm"
                name={user?.name}
                bg={selectedBg}
                color={selectedColor}
              />
              <VStack
                spacing="1px"
                align="start"
                display={{ base: 'none', md: 'flex' }}
              >
                <Text fontSize="sm" fontWeight="medium">
                  {user?.name || 'ゲスト'}
                </Text>
                <Text fontSize="xs" color="gray.500">
                  {user?.email}
                </Text>
              </VStack>
            </HStack>
          </MenuButton>
          <MenuList>
            <MenuItem onClick={handleLogout}>ログアウト</MenuItem>
          </MenuList>
        </Menu>
      </Flex>
    </Box>
  );
} 