# Attendance Management System - Backend API

Go + Echo + PostgreSQLを使用した勤怠管理システムのバックエンドAPIです。

## API仕様書

📖 [API Documentation (Swagger UI)](https://ryuki20.github.io/Attendance-Backend/)

OpenAPI 3.0仕様書: [docs/openapi.yaml](docs/openapi.yaml)

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

## GitHub Pagesでのドキュメント公開

API仕様書をGitHub Pagesで公開するための設定手順：

### 1. GitHub Pagesの有効化

1. GitHubリポジトリの`Settings`に移動
2. 左メニューから`Pages`を選択
3. `Source`を`GitHub Actions`に変更

### 2. 自動デプロイ

`docs/`ディレクトリ内のファイルを変更してmainブランチにpushすると、GitHub Actionsが自動的にデプロイします。

ワークフローファイル: `.github/workflows/deploy-docs.yml`

### 3. ドキュメントへのアクセス

デプロイ後、以下のURLでアクセスできます：
```
https://<username>.github.io/<repository-name>/
```

例: `https://myuto.github.io/Attendance-Backend/`

### 4. ローカルでのプレビュー

ローカルでドキュメントをプレビューする場合：

```bash
# シンプルなHTTPサーバーを起動
cd docs
python3 -m http.server 8000

# ブラウザで http://localhost:8000 を開く
```

## ライセンス

MIT
