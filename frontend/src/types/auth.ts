/**
 * ユーザーの役割を定義する列挙型
 */
export enum UserRole {
  ADMIN = 'ADMIN',
  MANAGER = 'MANAGER',
  USER = 'USER',
}

/**
 * ユーザーのステータス
 */
export enum UserStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  SUSPENDED = 'suspended',
}

/**
 * ユーザー情報の型定義
 */
export interface User {
  id: number;
  email: string;
  name: string;
  role: UserRole;
  createdAt: string;
  updatedAt: string;
}

/**
 * ログインリクエスト
 */
export interface LoginRequest {
  email: string;
  password: string;
}

/**
 * ログインレスポンス
 */
export interface LoginResponse {
  token: string;
  user: User;
}

/**
 * 登録リクエスト
 */
export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  name: string;
}

/**
 * プロフィール更新リクエスト
 */
export interface UpdateProfileRequest {
  name: string;
  email: string;
}

export interface UpdatePasswordRequest {
  currentPassword: string;
  newPassword: string;
} 