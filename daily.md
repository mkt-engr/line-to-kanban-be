# 開発日誌

## 2025-11-18

### 概要

LINEメッセージをCockroachDBに保存する機能を実装しました。

### 完了したタスク

#### 1. 環境変数の設定

- `.env`と`.env.example`に`DATABASE_URL`を追加
- データベース接続文字列を環境変数で管理

#### 2. データベーススキーマ設計

`migrations/001_create_messages_table.sql`を作成:

- `message_status` ENUM型を定義 (todo/in_progress/done)
- `messages`テーブル
  - `id`: UUID (主キー、自動生成)
  - `content`: TEXT (メッセージ内容)
  - `status`: message_status (ステータス、デフォルトtodo)
  - `created_at`: TIMESTAMP (作成日時)
  - `updated_at`: TIMESTAMP (更新日時)
- `created_at`にインデックスを追加（降順検索の高速化）
- `status`にインデックスを追加（ステータス検索の高速化）

#### 3. 依存関係の追加

- `github.com/jackc/pgx/v5`: PostgreSQL/CockroachDB用のGoドライバー
- Go標準の`database/sql`パッケージを使用

#### 4. リポジトリ実装

`internal/adapter/repository/cockroachdb/message_repository.go`を作成:

- `Save()`: メッセージをDBに保存
- `FindByID()`: IDでメッセージを検索
- `FindAll()`: 全メッセージを取得（作成日時降順）
- `UpdateStatus()`: 未実装（今回のスコープ外）

#### 5. 設定管理の追加

`internal/platform/config/database.go`を作成:

- `DATABASE_URL`環境変数を読み込むヘルパー関数

#### 6. アプリケーション統合

`cmd/api/main.go`を更新:

- データベース接続の初期化
- 接続テスト（Ping）
- リポジトリのDI
- LINE Webhookハンドラーへのリポジトリ注入

`internal/adapter/line/webhook_handler.go`を更新:

- メッセージリポジトリをフィールドに追加
- テキストメッセージ受信時にデータベースへ保存
- 保存エラー時もユーザーへの返信は継続

#### 7. マイグレーションスクリプト

`scripts/migrate.sh`を作成:

- マイグレーションSQLを実行するシェルスクリプト
- psqlコマンドまたはCloud Console経由で実行可能

#### 8. sqlc導入とコード生成

sqlcとは:
- SQLからタイプセーフなGoコードを自動生成するツール
- データベースクエリを型安全に実行できる
- 手動でScanを書く必要がなくなる
  - Scanとは: SQLクエリ結果をGoの変数にマッピングする処理
  - 通常は`row.Scan(&id, &content, &status, ...)`のように手動で書く必要がある
  - カラムの順番を間違えるとバグになり、型も手動で合わせる必要がある
  - sqlcは自動生成したコード内でScan処理を行うため、開発者は関数を呼ぶだけでOK
- SQLファイルからGoの関数とデータ型を生成

sqlc.yamlの設定内容:
```yaml
version: "2"
sql:
  - engine: "postgresql"           # PostgreSQL互換（CockroachDB）
    queries: "queries/"             # SQLクエリファイルの場所
    schema: "migrations/"           # スキーマ定義（マイグレーションファイル）
    gen:
      go:
        package: "db"                      # 生成されるGoパッケージ名
        out: "internal/adapter/repository/db"  # 出力先ディレクトリ
        sql_package: "pgx/v5"              # pgx v5ドライバーを使用
        emit_json_tags: true               # JSONタグを付ける
        emit_interface: true               # インターフェースを生成
        emit_empty_slices: true            # 空スライスを生成
        emit_pointers_for_null_types: true # NULL許容型をポインタに
```

コード生成方法:
```bash
sqlc generate
```

このコマンドで以下が自動生成されます:
- `internal/adapter/repository/db/db.go`: インターフェース定義
- `internal/adapter/repository/db/models.go`: データ型（Message構造体など）
- `internal/adapter/repository/db/messages.sql.go`: クエリ実行関数
  - `CreateMessage()` （今回はこれのみ使用）

生成されたコードの使い方:

1. データベース接続プールを作成:
```go
import "github.com/jackc/pgx/v5/pgxpool"

ctx := context.Background()
dbPool, err := pgxpool.New(ctx, databaseURL)
```

2. sqlc Queriesインスタンスを作成:
```go
import "line-to-kanban-be/internal/adapter/repository/db"

queries := db.New(dbPool)
```

3. 生成された関数を使ってクエリ実行:
```go
// メッセージ作成
savedMsg, err := queries.CreateMessage(ctx, db.CreateMessageParams{
    Content: "メッセージ内容",
    Status:  db.MessageStatusTodo,
})
```

sqlcが自動生成した型:
- `db.MessageStatus`: ENUM型 (todo/in_progress/done)
- `db.Message`: messagesテーブルの構造体
- `db.CreateMessageParams`: CreateMessage関数のパラメータ構造体

### アーキテクチャの変更点

依存関係フロー:

```
cmd/api/main.go
  |
  v
adapter/line/webhook_handler.go
  | (依存)
  v
domain/message/repository.go (インターフェース)
  ^
  | (実装)
adapter/repository/cockroachdb/message_repository.go
```

クリーンアーキテクチャの原則に従い、domainレイヤーがインターフェースを定義し、adapterレイヤーが具体的な実装を提供する構造を維持しました。

### 次のステップ

1. マイグレーションSQLを実行してテーブル作成
2. sqlcでクエリファイルを作成
3. sqlc generateでコード生成
4. 生成されたコードを使ってリポジトリを書き換え
5. アプリケーションをビルド＆起動してDB接続確認
6. LINEからメッセージを送信してDB保存をテスト

### 技術スタック

- **データベース**: CockroachDB (PostgreSQL互換)
- **ドライバー**: pgx v5
- **言語**: Go 1.23
- **アーキテクチャ**: クリーンアーキテクチャ
- **コード生成**: sqlc
