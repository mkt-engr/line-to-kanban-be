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
│   │       ├── usecase.go       # メッセージ処理のユースケース
│   │       └── dto.go           # 入出力DTO
│   ├── adapter/                 # アダプター層（外部I/O実装）
│   │   ├── http/
│   │   │   ├── router.go                # ルーティング設定
│   │   │   ├── hello_handler.go         # GET / (hello world)
│   │   │   ├── health_handler.go        # GET /healthz
│   │   │   ├── line_webhook_handler.go  # POST /line/webhook
│   │   │   ├── admin_handler.go         # POST /admin/tasks
│   │   │   └── kanban_handler.go        # PUT /kanban/status
│   │   └── repository/
│   │       └── memory/
│   │           └── message_repository.go # インメモリリポジトリ実装
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

- **domain 層**: ビジネスルールを持つエンティティ。他の層に依存しない
- **app 層**: ユースケース実装。domain に依存
- **adapter 層**: 外部 I/O（HTTP、DB 等）の実装。app と domain に依存
- **platform 層**: 設定やロガーなど横断的な機能

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

## 環境変数

- `PORT`: サーバーポート番号（デフォルト: 8080）
