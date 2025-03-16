/**
 * 通知設定の型定義
 */
export interface NotificationSettings {
  emailNotifications: boolean;
  lowStockAlert: boolean;
  deliveryUpdates: boolean;
  orderUpdates: boolean;
}

/**
 * 表示設定の型定義
 */
export interface DisplaySettings {
  theme: 'light' | 'dark';
  language: 'ja' | 'en';
  timezone: string;
  dateFormat: string;
}

/**
 * システム設定の型定義
 */
export interface SystemSettings {
  lowStockThreshold: number;
  defaultPageSize: number;
  autoLogoutMinutes: number;
}

/**
 * 設定全体の型定義
 */
export interface Settings {
  notification: NotificationSettings;
  display: DisplaySettings;
  system: SystemSettings;
}

/**
 * 設定更新リクエストの型定義
 */
export type UpdateSettingsRequest = Partial<Settings>; 