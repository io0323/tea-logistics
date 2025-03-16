# 茶物流管理システム

茶葉の在庫管理、出荷管理、入荷管理を一元化して行うための物流管理システムです。

## 機能

- **在庫管理**
  - 商品の在庫状況の確認
  - 在庫の入出荷記録
  - 在庫アラート設定

- **出荷管理**
  - 出荷指示の作成・編集
  - 出荷状況の追跡
  - 配送スケジュール管理
  - 顧客情報管理

- **入荷管理**
  - 発注管理
  - 入荷予定の管理
  - 仕入先情報管理
  - 入荷実績の記録

- **レポート機能**
  - 在庫推移レポート
  - 出荷実績レポート
  - 入荷実績レポート
  - カスタムレポート作成

## 技術スタック

### フロントエンド
- Next.js 14
- TypeScript
- Chakra UI
- React Icons
- React Query

### 認証
- JWT認証
- ロールベースのアクセス制御

## 開発環境のセットアップ

### 必要条件
- Node.js 18.0.0以上
- npm 9.0.0以上

### インストール手順

1. リポジトリのクローン
```bash
git clone https://github.com/yourusername/tea-logistics.git
cd tea-logistics
```

2. フロントエンドの依存関係のインストール
```bash
cd frontend
npm install
```

3. 開発サーバーの起動
```bash
npm run dev
```

アプリケーションは http://localhost:3000 で起動します。

## 環境変数の設定

フロントエンドの`.env.local`ファイルを作成し、以下の環境変数を設定してください：

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_AUTH_TOKEN_KEY=auth_token
```

## ディレクトリ構造

```
frontend/
├── src/
│   ├── app/              # ページコンポーネント
│   ├── components/       # 共通コンポーネント
│   ├── hooks/           # カスタムフック
│   ├── lib/             # ユーティリティ関数
│   ├── types/           # 型定義
│   └── utils/           # ヘルパー関数
├── public/              # 静的ファイル
└── package.json         # 依存関係の定義
```

## 使用方法

### ログイン
- メールアドレス: admin@example.com
- パスワード: admin123

### 基本操作
1. ダッシュボードで全体の状況を確認
2. 各機能は上部ナビゲーションメニューからアクセス
3. 詳細情報は各一覧の詳細ボタンから確認可能

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。 