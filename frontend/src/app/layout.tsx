import { Providers } from '@/providers/providers';
import { metadata } from './metadata';
import './globals.css';

export { metadata };

/**
 * ルートレイアウト
 */
export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
