# 茶物流管理システム

茶葉の在庫管理、出荷管理、入荷管理を一元化して行うための物流管理システムです。

## 機能

- **在庫管理**
  - 商品の在庫状況の確認
  - 在庫の入出荷記録
  - 在庫アラート設定
  - リアルタイム在庫更新

- **出荷管理**
  - 出荷指示の作成・編集
  - 出荷状況の追跡
  - 配送スケジュール管理
  - 顧客情報管理
  - QRコードによる出荷確認

- **入荷管理**
  - 発注管理
  - 入荷予定の管理
  - 仕入先情報管理
  - 入荷実績の記録
  - バーコードスキャン対応

- **レポート機能**
  - 在庫推移レポート
  - 出荷実績レポート
  - 入荷実績レポート
  - カスタムレポート作成
  - データエクスポート機能

## 技術スタック

### フロントエンド
- Next.js 14
- TypeScript
- Chakra UI
- React Query
- React Icons
- Chart.js
- date-fns

### バックエンド
- Go 1.21
- NestJS
- PostgreSQL
- Prisma ORM
- JWT認証
- OpenAPI/Swagger

### インフラ
- Docker
- GitHub Actions
- Nginx
- AWS (予定)

## 開発環境のセットアップ

### 必要条件
- Node.js 18.0.0以上
- npm 9.0.0以上
- Go 1.21以上
- PostgreSQL 15以上
- Docker (オプション)

### インストール手順

1. リポジトリのクローン
```bash
git clone https://github.com/io0323/tea-logistics.git
cd tea-logistics
```

2. フロントエンドのセットアップ
```bash
cd frontend
npm install
npm run dev
```

3. バックエンドのセットアップ
```bash
# PostgreSQLの起動確認
# Goバックエンド
cd ../
go mod download
go run cmd/api/main.go

# NestJSバックエンド
cd backend
npm install
npx prisma migrate dev
npm run start:dev
```

## 環境変数の設定

### フロントエンド (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_AUTH_TOKEN_KEY=auth_token
```

### バックエンド (.env)
```env
# PostgreSQL設定
DATABASE_URL="postgresql://user:password@localhost:5432/tea_logistics?schema=public"

# JWT設定
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION=24h

# サーバー設定
PORT=8080
ENV=development
```

## ディレクトリ構造

```
.
├── frontend/                 # フロントエンドアプリケーション
│   ├── src/
│   │   ├── app/            # ページコンポーネント
│   │   ├── components/     # 共通コンポーネント
│   │   ├── hooks/         # カスタムフック
│   │   ├── lib/           # ユーティリティ関数
│   │   ├── types/         # 型定義
│   │   └── utils/         # ヘルパー関数
│   └── public/            # 静的ファイル
│
├── backend/                 # NestJSバックエンド
│   ├── src/
│   │   ├── auth/          # 認証モジュール
│   │   ├── prisma/        # Prismaサービス
│   │   └── common/        # 共通モジュール
│   └── prisma/            # Prismaスキーマ
│
├── cmd/                    # Goのエントリーポイント
├── internal/              # 内部パッケージ
├── pkg/                   # 公開パッケージ
│   ├── handlers/         # HTTPハンドラー
│   ├── models/          # データモデル
│   ├── repository/      # データアクセス層
│   └── services/        # ビジネスロジック
│
└── docs/                  # ドキュメント
```

## API仕様

APIの詳細な仕様は以下のURLで確認できます：
- Go Backend: `http://localhost:8080/swagger/index.html`
- NestJS Backend: `http://localhost:3001/api`

## 認証情報

デフォルトの管理者アカウント：
- メールアドレス: admin@example.com
- パスワード: admin123

## 開発ガイドライン

1. コーディング規約
   - Go: [Effective Go](https://golang.org/doc/effective_go)
   - TypeScript: [Google TypeScript Style Guide](https://google.github.io/styleguide/tsguide.html)

2. コミットメッセージ
   - 形式: `type(scope): description`
   - 例: `feat(auth): implement JWT authentication`

3. ブランチ戦略
   - main: 本番環境
   - develop: 開発環境
   - feature/*: 機能開発
   - fix/*: バグ修正

## テスト

### フロントエンド
```bash
cd frontend
npm run test
```

### バックエンド
```bash
# Goテスト
go test ./...

# NestJSテスト
cd backend
npm run test
```

## デプロイ

現在、以下のデプロイ方法を計画中：
1. AWS ECS (コンテナ)
2. AWS RDS (PostgreSQL)
3. AWS S3 (静的アセット)
4. CloudFront (CDN)

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 貢献

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'feat: add amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## お問い合わせ

バグ報告や機能リクエストは[Issue](https://github.com/io0323/tea-logistics/issues)にて受け付けています。 