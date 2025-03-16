'use client';

import {
  Box,
  Flex,
  HStack,
  IconButton,
  Button,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  useDisclosure,
  useColorModeValue,
  Stack,
  Text,
} from '@chakra-ui/react';
import { HamburgerIcon, CloseIcon } from '@chakra-ui/icons';
import { useAuth } from '@/hooks/useAuth';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

interface NavItem {
  label: string;
  href: string;
  roles?: string[];
}

const NavItems: NavItem[] = [
  { label: 'ダッシュボード', href: '/dashboard' },
  { label: '在庫管理', href: '/inventory', roles: ['admin', 'user'] },
  { label: '出荷管理', href: '/shipping', roles: ['admin', 'user'] },
  { label: '入荷管理', href: '/receiving', roles: ['admin', 'user'] },
  { label: 'レポート', href: '/reports', roles: ['admin'] },
  { label: 'ユーザー管理', href: '/users', roles: ['admin'] },
];

export default function Navigation() {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const { user, logout } = useAuth();
  const pathname = usePathname();

  const filteredNavItems = NavItems.filter(item => {
    if (!item.roles) return true;
    return user && item.roles.includes(user.role);
  });

  return (
    <Box bg={useColorModeValue('white', 'gray.900')} px={4} shadow="md">
      <Flex h={16} alignItems="center" justifyContent="space-between">
        <IconButton
          size="md"
          icon={isOpen ? <CloseIcon /> : <HamburgerIcon />}
          aria-label="メニューを開く"
          display={{ md: 'none' }}
          onClick={isOpen ? onClose : onOpen}
        />

        <HStack spacing={8} alignItems="center">
          <Box>
            <Text fontSize="lg" fontWeight="bold">
              茶物流管理システム
            </Text>
          </Box>
          <HStack as="nav" spacing={4} display={{ base: 'none', md: 'flex' }}>
            {filteredNavItems.map((item) => (
              <Link key={item.href} href={item.href}>
                <Button
                  variant={pathname === item.href ? 'solid' : 'ghost'}
                  colorScheme={pathname === item.href ? 'green' : 'gray'}
                >
                  {item.label}
                </Button>
              </Link>
            ))}
          </HStack>
        </HStack>

        <Flex alignItems="center">
          <Menu>
            <MenuButton
              as={Button}
              rounded="full"
              variant="link"
              cursor="pointer"
              minW={0}
            >
              {user?.name || 'ユーザー'}
            </MenuButton>
            <MenuList>
              <MenuItem onClick={logout}>ログアウト</MenuItem>
            </MenuList>
          </Menu>
        </Flex>
      </Flex>

      {isOpen && (
        <Box pb={4} display={{ md: 'none' }}>
          <Stack as="nav" spacing={4}>
            {filteredNavItems.map((item) => (
              <Link key={item.href} href={item.href}>
                <Button
                  w="full"
                  variant={pathname === item.href ? 'solid' : 'ghost'}
                  colorScheme={pathname === item.href ? 'green' : 'gray'}
                >
                  {item.label}
                </Button>
              </Link>
            ))}
          </Stack>
        </Box>
      )}
    </Box>
  );
} 