# API Documentation

このディレクトリには、Attendance Management APIのOpenAPI仕様書とSwagger UIが含まれています。

## ファイル構成

```
docs/
├── openapi.yaml              # メインのOpenAPI仕様書
├── index.html                # Swagger UIのエントリーポイント
├── README.md                 # このファイル
└── openapi/                  # OpenAPI定義ファイル（分割管理）
    ├── paths/                # エンドポイント定義
    │   ├── health.yaml
    │   ├── auth-*.yaml       # 認証関連
    │   ├── attendances-*.yaml # 勤怠管理関連
    │   └── admin-*.yaml      # 管理者用
    ├── schemas/              # データモデル定義
    │   ├── User.yaml
    │   ├── Attendance.yaml
    │   └── Error.yaml
    └── responses/            # 共通レスポンス定義
        ├── UnauthorizedError.yaml
        └── ForbiddenError.yaml
```

## 特徴

- **実装との整合性**: Go実装と完全に一致したAPI仕様
- **モジュール化**: パス、スキーマ、レスポンスを分離して管理
- **視認性向上**: ファイル分割により各エンドポイントの仕様が明確
- **メンテナンス性**: 新しいエンドポイント追加時は対応するファイルを追加するだけ

## ローカルでのプレビュー

ローカル環境でドキュメントを確認する場合：

```bash
# このディレクトリに移動
cd docs

# Pythonの簡易HTTPサーバーを起動
python3 -m http.server 8000

# ブラウザで以下のURLを開く
# http://localhost:8000
```

または、他のHTTPサーバーを使用：

```bash
# Node.jsのhttp-serverを使用する場合
npx http-server -p 8000

# PHPを使用する場合
php -S localhost:8000
```

## GitHub Pagesでの公開

このドキュメントは、mainブランチにpushすると自動的にGitHub Pagesにデプロイされます。

公開URL: `https://<username>.github.io/<repository-name>/`

## OpenAPI仕様書の編集

### 新しいエンドポイントを追加する場合

1. `openapi/paths/` に新しいYAMLファイルを作成
2. `openapi.yaml` の `paths` セクションに参照を追加

例：
```yaml
# openapi/paths/users-list.yaml を作成

# openapi.yaml に追加
paths:
  /api/v1/users:
    $ref: './openapi/paths/users-list.yaml'
```

### 新しいスキーマを追加する場合

1. `openapi/schemas/` に新しいYAMLファイルを作成
2. `openapi.yaml` の `components.schemas` セクションに参照を追加

例：
```yaml
# openapi/schemas/Department.yaml を作成

# openapi.yaml に追加
components:
  schemas:
    Department:
      $ref: './openapi/schemas/Department.yaml'
```

### 推奨エディタ・ツール

- [Swagger Editor](https://editor.swagger.io/) - オンラインエディタ
- [VS Code](https://code.visualstudio.com/) + [OpenAPI (Swagger) Editor](https://marketplace.visualstudio.com/items?itemName=42Crunch.vscode-openapi) 拡張機能
- [Stoplight Studio](https://stoplight.io/studio) - GUI編集ツール

## API仕様の検証

OpenAPI仕様書の妥当性を検証：

```bash
# npm経由でswagger-cliをインストール
npm install -g @apidevtools/swagger-cli

# 仕様書を検証
cd docs
swagger-cli validate openapi.yaml
```

または、オンラインで検証：
https://editor.swagger.io/ にopenapi.yamlの内容を貼り付け

## Swagger UIの設定

### 基本設定

`index.html` 内の `SwaggerUIBundle` の設定を変更することで、表示をカスタマイズできます：

```javascript
const ui = SwaggerUIBundle({
    url: "openapi.yaml",           // OpenAPI仕様書のURL
    dom_id: '#swagger-ui',         // 表示先のDOM ID
    deepLinking: true,             // URLでエンドポイントに直接リンク
    docExpansion: "list",          // デフォルトの展開状態
    filter: true,                  // 検索フィルター有効化
    tryItOutEnabled: true,         // Try it out機能有効化
    // その他のオプション...
});
```

### 利用可能なオプション

詳細は[Swagger UIのドキュメント](https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/)を参照してください。

### テーマのカスタマイズ

Swagger UIのテーマをカスタマイズする場合は、`index.html` の `<style>` タグ内にCSSを追加してください。

## API仕様とコードの同期

実装とOpenAPI仕様書の整合性を保つため、以下を推奨します：

1. **エンドポイント追加時**: 実装と同時にOpenAPI定義も作成
2. **レスポンス変更時**: OpenAPI定義も更新
3. **定期的な検証**: 実装とOpenAPIの差異をチェック

## トラブルシューティング

### Swagger UIで仕様が読み込まれない場合

1. ブラウザのコンソールでエラーを確認
2. `$ref` のパスが正しいか確認（相対パスに注意）
3. YAMLの構文エラーがないか確認

### GitHub Pagesで表示されない場合

1. リポジトリの Settings > Pages で GitHub Actions が有効か確認
2. Actions タブでデプロイのステータスを確認
3. `docs/` ディレクトリの変更が main ブランチにpushされているか確認

## 参考リンク

- [OpenAPI Specification 3.0.3](https://spec.openapis.org/oas/v3.0.3)
- [Swagger UI Documentation](https://swagger.io/docs/open-source-tools/swagger-ui/)
- [OpenAPI Best Practices](https://oai.github.io/Documentation/best-practices.html)
