import {
  Flex,
  FlexProps,
  Avatar,
  Button,
  Menu,
  MenuButton,
  MenuList,
  MenuItem,
  MenuDivider,
  Text,
  useColorMode,
  useColorModeValue,
} from '@chakra-ui/react';
import { FiMoon, FiSun, FiUser, FiLogOut } from 'react-icons/fi';
import { useAuth } from '@/hooks/useAuth';
import { User } from '@/types/auth';

interface NavbarContentProps extends FlexProps {
  user?: User;
}

/**
 * ナビバーコンテンツ
 */
export const NavbarContent = ({ user, ...rest }: NavbarContentProps) => {
  const { colorMode, toggleColorMode } = useColorMode();
  const { logout } = useAuth();

  return (
    <Flex alignItems="center" {...rest}>
      <Button
        variant="ghost"
        onClick={toggleColorMode}
        aria-label="ダークモード切り替え"
        mr={4}
      >
        {colorMode === 'light' ? <FiMoon /> : <FiSun />}
      </Button>

      <Menu>
        <MenuButton
          as={Button}
          rounded="full"
          variant="link"
          cursor="pointer"
          minW={0}
        >
          <Avatar
            size="sm"
            name={user?.name}
            bg={useColorModeValue('blue.500', 'blue.200')}
            color="white"
          />
        </MenuButton>
        <MenuList>
          <Text px="3" py="2" fontSize="sm" fontWeight="bold">
            {user?.name}
          </Text>
          <Text px="3" pb="2" fontSize="xs" color="gray.500">
            {user?.email}
          </Text>
          <MenuDivider />
          <MenuItem icon={<FiUser />}>プロフィール</MenuItem>
          <MenuItem icon={<FiLogOut />} onClick={logout}>
            ログアウト
          </MenuItem>
        </MenuList>
      </Menu>
    </Flex>
  );
}; 