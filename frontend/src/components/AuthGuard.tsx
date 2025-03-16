'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/hooks/useAuth';
import { UserRole } from '@/types/auth';
import { Spinner, Center } from '@chakra-ui/react';

interface AuthGuardProps {
  children: React.ReactNode;
  requiredRoles?: UserRole[];
}

/**
 * 認証ガードコンポーネント
 * 認証が必要なページを保護し、必要な権限を持っているかチェックする
 */
export default function AuthGuard({
  children,
  requiredRoles,
}: AuthGuardProps) {
  const { isAuthenticated, loading, user } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push('/login');
    }

    if (
      !loading &&
      isAuthenticated &&
      requiredRoles &&
      user &&
      !requiredRoles.includes(user.role as UserRole)
    ) {
      router.replace('/unauthorized');
    }
  }, [isAuthenticated, loading, router, requiredRoles, user]);

  if (loading) {
    return (
      <Center h="100vh">
        <Spinner size="xl" />
      </Center>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  if (requiredRoles && user && !requiredRoles.includes(user.role as UserRole)) {
    return null;
  }

  return <>{children}</>;
} 