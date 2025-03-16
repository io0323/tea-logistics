'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import {
  Box,
  Button,
  FormControl,
  FormLabel,
  Input,
  VStack,
  Text,
  useToast,
  Avatar,
  HStack,
  IconButton,
  Divider,
} from '@chakra-ui/react';
import { EditIcon } from '@chakra-ui/icons';
import { useAuth } from '@/hooks/useAuth';
import DashboardLayout from '@/components/DashboardLayout';
import AuthGuard from '@/components/AuthGuard';

/**
 * プロフィール更新のバリデーションスキーマ
 */
const profileSchema = yup.object().shape({
  name: yup.string().required('名前は必須です'),
  email: yup
    .string()
    .email('有効なメールアドレスを入力してください')
    .required('メールアドレスは必須です'),
});

/**
 * パスワード変更のバリデーションスキーマ
 */
const passwordSchema = yup.object().shape({
  currentPassword: yup.string().required('現在のパスワードは必須です'),
  newPassword: yup
    .string()
    .min(8, '新しいパスワードは8文字以上である必要があります')
    .required('新しいパスワードは必須です'),
  confirmPassword: yup
    .string()
    .oneOf([yup.ref('newPassword')], 'パスワードが一致しません')
    .required('パスワードの確認は必須です'),
});

/**
 * プロフィール編集ページ
 */
export default function ProfilePage() {
  const { user, updateProfile, updatePassword } = useAuth();
  const toast = useToast();
  const [isLoading, setIsLoading] = useState(false);
  const [isPasswordLoading, setIsPasswordLoading] = useState(false);

  const {
    register: profileRegister,
    handleSubmit: handleProfileSubmit,
    formState: { errors: profileErrors },
  } = useForm({
    resolver: yupResolver(profileSchema),
    defaultValues: {
      name: user?.name || '',
      email: user?.email || '',
    },
  });

  const {
    register: passwordRegister,
    handleSubmit: handlePasswordSubmit,
    formState: { errors: passwordErrors },
    reset: resetPassword,
  } = useForm({
    resolver: yupResolver(passwordSchema),
  });

  const onProfileSubmit = async (data: any) => {
    try {
      setIsLoading(true);
      await updateProfile.mutateAsync(data);
      toast({
        title: 'プロフィールを更新しました',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
    } catch (error) {
      toast({
        title: '更新に失敗しました',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  const onPasswordSubmit = async (data: any) => {
    try {
      setIsPasswordLoading(true);
      await updatePassword.mutateAsync({
        currentPassword: data.currentPassword,
        newPassword: data.newPassword,
      });
      toast({
        title: 'パスワードを変更しました',
        status: 'success',
        duration: 3000,
        isClosable: true,
      });
      resetPassword();
    } catch (error) {
      toast({
        title: 'パスワードの変更に失敗しました',
        status: 'error',
        duration: 3000,
        isClosable: true,
      });
    } finally {
      setIsPasswordLoading(false);
    }
  };

  return (
    <AuthGuard>
      <DashboardLayout>
        <Box maxW="container.md" mx="auto" py={8}>
          <VStack spacing={8} align="stretch">
            {/* プロフィール画像セクション */}
            <Box>
              <Text fontSize="2xl" fontWeight="bold" mb={4}>
                プロフィール画像
              </Text>
              <HStack spacing={4}>
                <Avatar
                  size="2xl"
                  name={user?.name}
                  src={user?.avatarUrl}
                />
                <IconButton
                  aria-label="Edit profile image"
                  icon={<EditIcon />}
                  size="sm"
                  variant="ghost"
                />
              </HStack>
            </Box>

            <Divider />

            {/* プロフィール情報セクション */}
            <Box>
              <Text fontSize="2xl" fontWeight="bold" mb={4}>
                プロフィール情報
              </Text>
              <form onSubmit={handleProfileSubmit(onProfileSubmit)}>
                <VStack spacing={4}>
                  <FormControl isInvalid={!!profileErrors.name}>
                    <FormLabel>名前</FormLabel>
                    <Input {...profileRegister('name')} />
                    {profileErrors.name && (
                      <Text color="red.500" fontSize="sm">
                        {profileErrors.name.message}
                      </Text>
                    )}
                  </FormControl>

                  <FormControl isInvalid={!!profileErrors.email}>
                    <FormLabel>メールアドレス</FormLabel>
                    <Input {...profileRegister('email')} />
                    {profileErrors.email && (
                      <Text color="red.500" fontSize="sm">
                        {profileErrors.email.message}
                      </Text>
                    )}
                  </FormControl>

                  <Button
                    type="submit"
                    colorScheme="blue"
                    isLoading={isLoading}
                    alignSelf="flex-start"
                  >
                    保存
                  </Button>
                </VStack>
              </form>
            </Box>

            <Divider />

            {/* パスワード変更セクション */}
            <Box>
              <Text fontSize="2xl" fontWeight="bold" mb={4}>
                パスワード変更
              </Text>
              <form onSubmit={handlePasswordSubmit(onPasswordSubmit)}>
                <VStack spacing={4}>
                  <FormControl isInvalid={!!passwordErrors.currentPassword}>
                    <FormLabel>現在のパスワード</FormLabel>
                    <Input
                      type="password"
                      {...passwordRegister('currentPassword')}
                    />
                    {passwordErrors.currentPassword && (
                      <Text color="red.500" fontSize="sm">
                        {passwordErrors.currentPassword.message}
                      </Text>
                    )}
                  </FormControl>

                  <FormControl isInvalid={!!passwordErrors.newPassword}>
                    <FormLabel>新しいパスワード</FormLabel>
                    <Input
                      type="password"
                      {...passwordRegister('newPassword')}
                    />
                    {passwordErrors.newPassword && (
                      <Text color="red.500" fontSize="sm">
                        {passwordErrors.newPassword.message}
                      </Text>
                    )}
                  </FormControl>

                  <FormControl isInvalid={!!passwordErrors.confirmPassword}>
                    <FormLabel>新しいパスワード（確認）</FormLabel>
                    <Input
                      type="password"
                      {...passwordRegister('confirmPassword')}
                    />
                    {passwordErrors.confirmPassword && (
                      <Text color="red.500" fontSize="sm">
                        {passwordErrors.confirmPassword.message}
                      </Text>
                    )}
                  </FormControl>

                  <Button
                    type="submit"
                    colorScheme="blue"
                    isLoading={isPasswordLoading}
                    alignSelf="flex-start"
                  >
                    パスワードを変更
                  </Button>
                </VStack>
              </form>
            </Box>
          </VStack>
        </Box>
      </DashboardLayout>
    </AuthGuard>
  );
} 