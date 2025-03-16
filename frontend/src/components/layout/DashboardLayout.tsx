'use client';

import { Box } from '@chakra-ui/react';
import Navbar from '@/components/Navbar';

interface DashboardLayoutProps {
  children: React.ReactNode;
}

/**
 * ダッシュボードレイアウトコンポーネント
 * 認証済みページの共通レイアウトを提供します
 */
export default function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <Box minH="100vh">
      <Navbar />
      <Box>{children}</Box>
    </Box>
  );
} 