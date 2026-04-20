# Attendance Management System - Backend API

Go + Echo + PostgreSQLを使用した勤怠管理システムのバックエンドAPIです。

## 技術スタック

- **Go** 1.23+
- **Echo** v4 - Webフレームワーク
- **PostgreSQL** 16 - データベース
- **JWT** - 認証
- **Docker** & **Docker Compose** - コンテナ化
- **Air** - ホットリロード
- **golang-migrate** - データベースマイグレーション

## プロジェクト構成

```
.
├── cmd/
│   └── server/          # メインアプリケーション
├── internal/
│   ├── domain/          # ドメイン層
│   │   ├── entity/      # エンティティ
│   │   └── repository/  # リポジトリインターフェース
│   ├── usecase/         # ユースケース層
│   ├── infrastructure/  # インフラ層
│   │   ├── database/    # DB接続
│   │   ├── repository/  # リポジトリ実装
│   │   └── router/      # ルーティング
│   └── interface/       # インターフェース層
│       ├── handler/     # HTTPハンドラー
│       └── middleware/  # ミドルウェア
├── pkg/                 # 共通パッケージ
│   ├── config/          # 設定管理
│   └── utils/           # ユーティリティ
└── migrations/          # DBマイグレーション
```

## セットアップ

### 前提条件

- Docker & Docker Compose がインストールされていること

### 環境変数の設定

`.env.example`を`.env`にコピーして、必要に応じて値を編集します：

```bash
cp .env.example .env
```

### Docker環境の起動

```bash
# コンテナのビルドと起動
make build
make up

# または一括で
docker-compose up --build
```

### データベースマイグレーション

```bash
# マイグレーション実行
make migrate-up

# マイグレーションロールバック
make migrate-down

# 新しいマイグレーションファイル作成
make migrate-create name=create_some_table
```

## API仕様

### 認証

#### ユーザー登録
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "山田太郎",
  "role": "employee"  // optional: admin, manager, employee
}
```

#### ログイン
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "山田太郎",
    "role": "employee"
  }
}
```

#### 認証情報取得
```
GET /api/v1/auth/me
Authorization: Bearer {token}
```

### 勤怠管理

#### 出勤
```
POST /api/v1/attendances/clock-in
Authorization: Bearer {token}
```

#### 退勤
```
POST /api/v1/attendances/clock-out
Authorization: Bearer {token}
```

#### 休憩開始
```
POST /api/v1/attendances/break-start
Authorization: Bearer {token}
```

#### 休憩終了
```
POST /api/v1/attendances/break-end
Authorization: Bearer {token}
```

#### 今日の勤怠取得
```
GET /api/v1/attendances/today
Authorization: Bearer {token}
```

#### 勤怠履歴取得
```
GET /api/v1/attendances/history?start_date=2024-01-01&end_date=2024-01-31
Authorization: Bearer {token}
```

### 管理者用API

#### 全ユーザーの勤怠取得（管理者/マネージャーのみ）
```
GET /api/v1/admin/attendances?start_date=2024-01-01&end_date=2024-01-31
Authorization: Bearer {token}
```

#### 特定ユーザーの勤怠取得（管理者/マネージャーのみ）
```
GET /api/v1/admin/attendances/user/:user_id?start_date=2024-01-01&end_date=2024-01-31
Authorization: Bearer {token}
```

## 開発コマンド

```bash
# コンテナ起動
make up

# コンテナ停止
make down

# ログ表示
make logs

# クリーンアップ（ボリューム含む）
make clean

# ヘルプ表示
make help
```

## ホットリロード

Airを使用しているため、コードを変更すると自動的にサーバーが再起動します。

## データベース

### テーブル構造

#### users
- id (serial)
- email (varchar)
- password_hash (varchar)
- name (varchar)
- role (varchar)
- created_at (timestamp)
- updated_at (timestamp)

#### attendances
- id (serial)
- user_id (integer)
- date (date)
- clock_in (timestamp)
- clock_out (timestamp)
- break_start (timestamp)
- break_end (timestamp)
- status (varchar)
- notes (text)
- created_at (timestamp)
- updated_at (timestamp)

## セキュリティ

- パスワードはbcryptでハッシュ化
- JWT認証を使用
- ロールベースのアクセス制御（RBAC）
- CORS設定

## 本番環境への展開

本番環境では以下を変更してください：

1. `.env`の`JWT_SECRET`を強力なランダム文字列に変更
2. データベースのパスワードを変更
3. `ENV`を`production`に設定
4. 適切なCORS設定を行う

## ライセンス

MIT
