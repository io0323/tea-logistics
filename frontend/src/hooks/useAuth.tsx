'use client';

import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { useRouter } from 'next/navigation';

// テスト用のユーザー情報
const TEST_USERS = [
  {
    id: '1',
    email: 'admin@example.com',
    password: 'admin123',
    name: '管理者',
    role: 'ADMIN'
  },
  {
    id: '2',
    email: 'user@example.com',
    password: 'user123',
    name: '一般ユーザー',
    role: 'USER'
  }
];

interface User {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'user';
}

interface AuthContextType {
  user: User | null;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

/**
 * 認証用のカスタムフック
 */
export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return {
    ...context,
    isAuthenticated: !!context.user,
  };
}

/**
 * 認証プロバイダーコンポーネント
 */
export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    // ローカルストレージからユーザー情報を取得
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    }
    setIsLoading(false);
  }, []);

  const login = async (email: string, password: string) => {
    try {
      // モックの認証処理
      let mockUser: User;
      if (email === 'admin@example.com' && password === 'admin123') {
        mockUser = {
          id: '1',
          name: '管理者',
          email: 'admin@example.com',
          role: 'admin',
        };
      } else if (email === 'user@example.com' && password === 'user123') {
        mockUser = {
          id: '2',
          name: '一般ユーザー',
          email: 'user@example.com',
          role: 'user',
        };
      } else {
        throw new Error('メールアドレスまたはパスワードが正しくありません');
      }

      // ユーザー情報をローカルストレージに保存
      localStorage.setItem('user', JSON.stringify(mockUser));
      setUser(mockUser);
      router.push('/dashboard');
    } catch (error) {
      throw error;
    }
  };

  const logout = () => {
    localStorage.removeItem('user');
    setUser(null);
    router.push('/login');
  };

  return (
    <AuthContext.Provider value={{ user, login, logout, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
} 