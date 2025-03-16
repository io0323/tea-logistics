'use client';

import { Box, Container, Heading, Text, VStack, useColorModeValue } from '@chakra-ui/react';
import { FiSettings } from 'react-icons/fi';
import { useSettings } from '@/hooks/useSettings';
import { SettingsForm } from '@/components/settings/SettingsForm';
import AuthGuard from '@/components/auth/AuthGuard';
import DashboardLayout from '@/components/layout/DashboardLayout';

/**
 * 設定画面コンポーネント
 */
export default function SettingsPage() {
  const { data: settings, isLoading, error } = useSettings();
  const bgColor = useColorModeValue('white', 'gray.800');

  if (isLoading) {
    return (
      <AuthGuard>
        <DashboardLayout>
          <Container maxW="container.xl" py={8}>
            <Text>読み込み中...</Text>
          </Container>
        </DashboardLayout>
      </AuthGuard>
    );
  }

  if (error) {
    return (
      <AuthGuard>
        <DashboardLayout>
          <Container maxW="container.xl" py={8}>
            <Text color="red.500">エラーが発生しました: {error.message}</Text>
          </Container>
        </DashboardLayout>
      </AuthGuard>
    );
  }

  return (
    <AuthGuard>
      <DashboardLayout>
        <Box bg={bgColor} minH="100vh" pt={16}>
          <Container maxW="container.xl" py={8}>
            <VStack spacing={8} align="stretch">
              <Box>
                <Heading size="lg" mb={2}>
                  設定
                </Heading>
                <Text color="gray.600">
                  システムの設定を管理します
                </Text>
              </Box>

              <SettingsForm initialSettings={settings} />
            </VStack>
          </Container>
        </Box>
      </DashboardLayout>
    </AuthGuard>
  );
} 