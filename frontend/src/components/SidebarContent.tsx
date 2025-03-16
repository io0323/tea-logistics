import {
  Box,
  BoxProps,
  CloseButton,
  Flex,
  Text,
  useColorModeValue,
} from '@chakra-ui/react';
import { NavItem } from './NavItem';
import {
  FiHome,
  FiBox,
  FiPackage,
  FiTruck,
  FiUsers,
  FiSettings,
} from 'react-icons/fi';
import { IconType } from 'react-icons';
import { useAuth } from '@/hooks/useAuth';
import { UserRole } from '@/types/auth';

interface LinkItemProps {
  name: string;
  icon: IconType;
  href: string;
  roles?: UserRole[];
}

/**
 * サイドバーのナビゲーションリンク
 */
const LinkItems: Array<LinkItemProps> = [
  { name: 'ダッシュボード', icon: FiHome, href: '/dashboard' },
  { name: '商品管理', icon: FiBox, href: '/products' },
  { name: '在庫管理', icon: FiPackage, href: '/inventory' },
  { name: '配送管理', icon: FiTruck, href: '/deliveries' },
  {
    name: 'ユーザー管理',
    icon: FiUsers,
    href: '/users',
    roles: [UserRole.ADMIN, UserRole.MANAGER],
  },
  {
    name: '設定',
    icon: FiSettings,
    href: '/settings',
    roles: [UserRole.ADMIN],
  },
];

interface SidebarContentProps extends BoxProps {
  onClose: () => void;
}

/**
 * サイドバーコンテンツ
 */
export const SidebarContent = ({ onClose, ...rest }: SidebarContentProps) => {
  const { user } = useAuth();

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
        <Text fontSize="2xl" fontWeight="bold">
          茶葉物流管理
        </Text>
        <CloseButton display={{ base: 'flex', md: 'none' }} onClick={onClose} />
      </Flex>
      {LinkItems.map((link) => {
        // 権限チェック
        if (link.roles && (!user || !link.roles.includes(user.role))) {
          return null;
        }

        return (
          <NavItem
            key={link.name}
            icon={link.icon}
            href={link.href}
            onClick={onClose}
          >
            {link.name}
          </NavItem>
        );
      })}
    </Box>
  );
}; 