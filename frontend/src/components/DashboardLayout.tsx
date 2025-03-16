'use client';

import { ReactNode } from 'react';
import {
  Box,
  useColorModeValue,
} from '@chakra-ui/react';
import Navbar from './Navbar';

interface DashboardLayoutProps {
  children: ReactNode;
}

/**
 * ダッシュボードレイアウトコンポーネント
 */
export default function DashboardLayout({ children }: DashboardLayoutProps) {
  return (
    <Box minH="100vh" bg={useColorModeValue('gray.50', 'gray.900')}>
      <Navbar onOpen={() => {}} />
      <Box pt="20">
        {children}
      </Box>
    </Box>
  );
} 