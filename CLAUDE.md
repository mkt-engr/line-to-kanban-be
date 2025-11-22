# line-to-kanban-be

LINE 連携カンバンアプリケーションのバックエンド API

## ディレクトリ構成

```
line-to-kanban-be/
├── go.mod                       # Goモジュール定義
├── cmd/
│   └── api/
│       └── main.go              # エントリーポイント（DI・サーバー起動）
├── internal/
│   ├── domain/                  # ドメイン層（ビジネスロジック、依存なし）
│   │   └── message/
│   │       ├── entity.go        # メッセージエンティティ
│   │       └── repository.go    # リポジトリインターフェース
│   ├── app/                     # ユースケース層（アプリケーションロジック）
│   │   └── message/
│   │       ├── usecase.go       # Usecase構造体定義
│   │       ├── usecase_read.go  # Read系ユースケース
│   │       ├── usecase_write.go # Write系ユースケース
│   │       └── dto.go           # 入出力DTO
│   ├── adapter/                 # アダプター層（外部I/O実装）
│   │   ├── http/
│   │   │   ├── router.go                # ルーティング設定
│   │   │   └── line_webhook_handler.go  # POST /line/webhook
│   │   ├── line/
│   │   │   ├── client.go                # LINE Bot クライアント
│   │   │   └── webhook_handler.go       # LINE Webhook処理
│   │   └── repository/
│   │       ├── message_repository.go    # リポジトリ構造体定義
│   │       ├── message_read.go          # Read系メソッド（FindByID, FindByUserID）
│   │       ├── message_write.go         # Write系メソッド（Save, UpdateStatus, Delete）
│   │       ├── message_converter.go     # 型変換（domain ⇔ DB）
│   │       └── db/                      # sqlc生成コード
│   │           ├── models.go
│   │           ├── messages.sql.go
│   │           └── querier.go
│   └── platform/                # 横断的関心事
│       ├── config/
│       │   └── config.go        # 環境変数・設定管理
│       └── logger/
│           └── logger.go        # ロガー
└── app                          # ビルド成果物（.gitignore対象）
```

## アーキテクチャ

クリーンアーキテクチャを採用。依存関係は以下の方向：

```
adapter → app → domain
```

- domain 層: ビジネスルールを持つエンティティ。他の層に依存しない
- app 層: ユースケース実装。domain に依存
- adapter 層: 外部 I/O（HTTP、DB 等）の実装。app と domain に依存
- platform 層: 設定やロガーなど横断的な機能

### ファイル分割の設計方針

repository層とusecase層は機能ごとにファイルを分割し、Goの標準的なテスト規約に準拠：

Repository層のファイル構成:
- message_repository.go: 構造体定義とコンストラクタ（15行）
- message_read.go: Read系メソッド（約30行）
- message_write.go: Write系メソッド（約30行）
- message_converter.go: 型変換関数（約35行）

Usecase層のファイル構成:
- usecase.go: Usecase構造体定義とコンストラクタ（15行）
- usecase_read.go: Read系ユースケース（約20行）
- usecase_write.go: Write系ユースケース（約25行）

設計方針:
- 各ファイル15-40行程度で管理しやすく保つ
- テストファイルは_testサフィックス（例: message_read_test.go）
- 1ファイルに全メソッドを含めると将来的に100-200行超になるため分割
- Read/Write分割により、関連する機能をグループ化
- repository層とusecase層で統一したパターンを採用

## API エンドポイント

| メソッド | パス             | 説明                      |
| -------- | ---------------- | ------------------------- |
| GET      | `/`              | Hello World（動作確認用） |
| GET      | `/healthz`       | ヘルスチェック            |
| POST     | `/line/webhook`  | LINE からの Webhook 受信  |
| POST     | `/admin/tasks`   | 管理タスク作成            |
| PUT      | `/kanban/status` | カンバンステータス更新    |

## ビルド・実行

### ビルド

```bash
go build -o app cmd/api/main.go
```

### 実行

```bash
./app
```

デフォルトでポート 8080 で起動。環境変数`PORT`で変更可能。

### テスト

```bash
curl http://localhost:8080              # hello world
curl http://localhost:8080/healthz      # ヘルスチェック
```

## 依存関係

- `github.com/google/uuid` - UUID 生成（メッセージ ID 用）
- `github.com/line/line-bot-sdk-go/v7` - LINE Bot SDK
- `github.com/jackc/pgx/v5` - PostgreSQL ドライバー
- `github.com/joho/godotenv` - 環境変数読み込み

## データベース

- PostgreSQL を使用
- sqlc でクエリから型安全なコードを自動生成
- マイグレーション管理

## 環境変数

- `PORT`: サーバーポート番号（デフォルト: 8080）
- `LINE_CHANNEL_SECRET`: LINE Bot チャンネルシークレット
- `LINE_CHANNEL_ACCESS_TOKEN`: LINE Bot アクセストークン
- `DATABASE_URL`: PostgreSQL 接続URL
- to memorize