# 開発日誌

## 2025-11-22

### 概要

LINEからタスクを削除する機能を実装しました。その後、クリーンアーキテクチャを採用し、タスク作成機能をusecase経由にリファクタリングしました。さらに、保守性向上のためrepository層を4ファイルに分割しました。

### 完了したタスク

#### 2. クリーンアーキテクチャへのリファクタリング（タスク作成機能）

タスク作成機能をクリーンアーキテクチャに準拠するよう、usecase経由に変更しました。

実装内容:

1. domain層の修正:
   - `internal/domain/message/entity.go`: Status型を`todo/in_progress/done`に変更（DBスキーマに合わせる）
   - `Message`エンティティに`UserID`フィールドを追加
   - `NewMessage`関数に`userID`パラメータを追加

2. Repository実装の作成:
   - `internal/adapter/repository/message_repository.go`: 新規作成
   - domainの`Repository`インターフェースをsqlc実装でラップ
   - 型変換: `db.Message` ⇔ `domain.Message`
   - pgtype.UUIDとstring間の変換処理を実装
   - `Save`, `FindByUserID`, `Delete`メソッドを実装

3. app層（usecase）の修正:
   - `internal/app/message/usecase.go`: `CreateMessage`を修正してuserIDに対応
   - `internal/app/message/dto.go`: `CreateMessageRequest`と`MessageResponse`に`UserID`を追加

4. adapter層の修正:
   - `internal/adapter/line/webhook_handler.go`: タスク作成処理をusecase経由に変更
   - `queries.CreateMessage`の直接呼び出しから`usecase.CreateMessage`に変更
   - コマンド処理（一覧・削除）は引き続きqueries直接呼び出し（次回対応予定）

5. DI（依存性注入）の組み立て:
   - `cmd/api/main.go`: Repository → Usecase → Handlerの順で初期化
   - `NewMessageRepository(queries)`でRepository層を作成
   - `NewUsecase(messageRepo)`でUsecase層を作成
   - `NewWebhookHandler(lineClient, queries, messageUsecase)`でHandlerを作成

依存関係の流れ:
```
webhook_handler → usecase → repository → sqlc → DB
     (adapter)      (app)    (adapter)   (生成)
                      ↓
                   domain
```

メリット:
- テスト容易性: usecaseのテストでRepositoryをモック化可能
- ビジネスロジックの集約: タスク作成ロジックがusecase層に集約
- 型安全性: domain型で統一され、DB詳細が隠蔽される
- 保守性: DB変更時はrepository層のみ修正すれば良い

次回対応予定:
- ステータス更新機能の実装

#### 4. 一覧表示・削除機能のusecase化とusecase層の分割

一覧表示と削除機能をusecase経由に変更し、さらにusecase層もrepositoryと同じパターンで分割しました。

実装内容:

1. usecaseに新しいメソッドを追加:
   - `ListMessagesByUser`: ユーザーのメッセージ一覧を取得
   - `DeleteMessage`: メッセージを削除

2. webhook_handlerをusecase経由に変更:
   - `handleListCommand`: queries.ListMessagesByUser → usecase.ListMessagesByUser
   - `handleDeleteCommand`: queries.DeleteMessage → usecase.DeleteMessage
   - queriesフィールドを完全に削除

3. usecase層のファイル分割:
   - `usecase.go`: 構造体定義とコンストラクタ（15行）
   - `usecase_read.go`: Read系メソッド - ListMessagesByUser（約20行）
   - `usecase_write.go`: Write系メソッド - CreateMessage, UpdateMessageStatus, DeleteMessage（約25行）

メリット:
- 全てのDB操作がusecase経由に統一された
- repository層と同じパターンで分割し、一貫性が向上
- 各ファイルが15-25行程度で非常に読みやすい
- テストが書きやすい構造

#### 3. リポジトリ層のファイル分割

保守性向上のため、repository層を機能ごとに4ファイルに分割しました。

実装内容:

1. message_repository.goの分割:
   - 元の103行のファイルを4つに分割
   - 将来的に全メソッド実装で200行超になることを防ぐ

2. 新しいファイル構成:
   - `message_repository.go`: 構造体定義とコンストラクタ（15行）
   - `message_read.go`: Read系メソッド - FindByID, FindByUserID（約30行）
   - `message_write.go`: Write系メソッド - Save, UpdateStatus, Delete（約30行）
   - `message_converter.go`: 型変換関数 - toMessage, toDBStatus, toUUID（約35行）

3. 不要なメソッドを削除:
   - `FindAll`メソッドを削除（使用予定なし）
   - `GetAllMessages`ユースケースを削除

4. 設計方針:
   - 各ファイル15-40行程度で管理しやすく保つ
   - テストファイルは_testサフィックス（例: message_read_test.go）
   - Goの標準的なテスト規約に準拠
   - Read/Write分割により、関連する機能をグループ化

5. CLAUDE.mdの更新:
   - ディレクトリ構成を実際の構造に合わせて更新
   - リポジトリ層のファイル分割の設計方針を追記
   - データベースと環境変数の情報を追加

メリット:
- 各ファイルが小さく、理解しやすい
- テストが書きやすい（機能ごとにテストファイルを分割可能）
- 将来的なメソッド追加でもファイルサイズが適切に保たれる
- Git差分が見やすい（関連する変更が同じファイルに集約）

#### 1. 削除機能の実装

LINEで「削除1」または「削除 1」のようなコマンドを送信すると、一覧の1番目のタスクを削除できる機能を実装しました。

実装内容:

1. SQLクエリの追加:
   - `queries/messages.sql`に`DeleteMessage`クエリを追加
   - `id`と`user_id`の両方でフィルタリング（他人のタスクを削除できないように）

2. コマンド処理の分離:
   - 正規表現で「削除」コマンドを検知: `^削除\s*(\d+)$`
   - スペースの有無に対応（`削除1`、`削除 1`、`削除  1`すべてマッチ）
   - コマンドはDBに保存しない設計

3. Webhookハンドラーのリファクタリング:
   - `isCommand()`: コマンド判定関数
   - `handleCommand()`: コマンドルーター
   - `handleListCommand()`: 一覧表示処理（既存ロジックを移動）
   - `handleDeleteCommand()`: 削除処理（新規実装）
   - `replyError()`: エラー返信ヘルパー

4. 削除処理の流れ:
   - 番号を抽出・バリデーション
   - `ListMessagesByUser`でユーザーのタスク一覧を取得
   - N番目のタスクIDを特定
   - `DeleteMessage`で物理削除
   - 成功メッセージを返信

5. エラーハンドリング:
   - 0以下の数字: 「正しい番号を指定してください」
   - 範囲外の番号: 「指定されたタスクが見つかりません」
   - DB接続エラー: 「削除に失敗しました」

メリット:
- コマンドがDBに残らずクリーン
- スペースの有無を気にせず使いやすい
- スケーラブルな設計（個人用アプリでは全件取得でも問題なし）

### 設計方針と今後の拡張

#### スケーラビリティの考慮

現在の設計（全件取得）で対応可能な規模:
- ユーザー数: 数百万人
- 1ユーザーあたりのタスク数: 数千件
- 全件取得のコスト: 数ミリ秒（実用上問題なし）

理由:
- 各ユーザーのタスク数は独立
- データベースの負荷は線形に増加
- 個人用カンバンアプリでは1ユーザーあたり数百〜数千件が現実的

#### 将来の拡張オプション（必要になった場合のみ）

1. Redisセッションキャッシュ:
   - 一覧表示時にタスクリストをキャッシュ
   - 削除時はキャッシュから取得
   - 適用時期: 1ユーザーあたり数万件のタスクが発生した場合

2. Webフロントエンドでの直接指定:
   - ブラウザでカンバン表示
   - 削除ボタンでUUIDを直接指定
   - API: `DELETE /api/messages/{uuid}`

3. 論理削除:
   - `deleted_at`カラムを追加
   - 削除したタスクの復元機能
   - 削除履歴の表示

4. その他のコマンド拡張:
   - 完了コマンド: `^完了\s*(\d+)$`
   - 編集コマンド: `^編集\s*(\d+)\s+(.+)$`
   - ステータス変更コマンド

## 2025-11-19

### 概要

マイグレーション処理をsqlcで型安全に管理するよう改善し、messagesテーブルにuser_idカラムを追加しました。

### 完了したタスク

#### 1. マイグレーション処理のsqlc化

すべてのマイグレーション関連のSQLをファイルで管理し、sqlcで型安全なコードを生成するよう改善しました。

実装内容:

1. schema_migrationsテーブルの作成をSQLファイル化:
   - `migrations/000_init_schema_migrations.up.sql`を作成
   - Goコードにハードコードしていたテーブル定義を削除

2. マイグレーション管理クエリをsqlcで定義:
   - `queries/schema_migrations.sql`を作成
   - `GetCurrentMigrationVersion`: 現在のバージョンを取得
   - `InsertMigrationVersion`: 新しいバージョンを記録
   - 型キャスト(`::int`)でsqlcが正しい型(int32)を生成

3. migrate.goでsqlcの生成コードを使用:
   - 生SQLクエリをすべて削除
   - sqlcで生成された関数を使用
   - 型安全性が向上

メリット:
- すべてのSQLをファイルで管理（Goコードにハードコードしない）
- sqlcで型安全なバージョン管理
- SQLとGoコードの明確な分離
- 保守性の向上

#### 2. messagesテーブルにuser_idカラムを追加

どのLINEユーザーがメッセージを送信したかを追跡できるようにしました。

実装内容:

1. マイグレーションファイルを作成:
   - `migrations/002_add_user_id_to_messages.up.sql`
   - `user_id TEXT NOT NULL DEFAULT ''`カラムを追加
   - `user_id`にインデックスを追加（ユーザー別検索の高速化）

2. sqlcクエリを更新:
   - `queries/messages.sql`の`CreateMessage`に`user_id`パラメータを追加
   - sqlcでコード再生成

3. Webhookハンドラーを更新:
   - `internal/adapter/line/webhook_handler.go`
   - LINEから取得した`userID`をデータベースに保存

メリット:
- メッセージの送信者を追跡可能
- ユーザー別のメッセージ検索が高速化
- 将来的にユーザー別のタスク管理が可能

### 技術的な学び

1. SQLの型キャスト:
   - PostgreSQL/CockroachDBでは`::int`で型キャストが可能
   - TypeScriptの`as`とは異なり、実行時に型変換が行われるため安全
   - sqlcが正しい型を推論するために有効

2. マイグレーションのベストプラクティス:
   - すべてのSQLをファイルで管理
   - バージョン管理もsqlcで型安全に
   - Goコードには制御フローのみを記述

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

#### 4. 設定管理の追加

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

#### 9. 自動マイグレーション機能の実装（CockroachDB対応版）

CockroachDBは`pg_advisory_lock()`をサポートしていないため、golang-migrateの代わりにpgxpool + sqlcを使った独自のマイグレーション実装を作成しました。

pg_advisory_lock()とは:
- PostgreSQLのアドバイザリロック機能
- 複数のプロセスが同時にマイグレーションを実行するのを防ぐための排他制御
- golang-migrateなどのマイグレーションツールが内部で使用
- CockroachDBは分散データベースのため、この関数をサポートしていない
- 代わりにCockroachDBは独自のトランザクション管理機構を持つ

今回の対応:
- golang-migrateを使わず、独自のマイグレーション実装を作成
- トランザクション内でマイグレーションを実行することで排他制御を実現
- CockroachDBのトランザクション機能を活用

実装内容:

1. マイグレーションSQLファイル:
   - `migrations/000_init_schema_migrations.up.sql`: schema_migrationsテーブルの作成
   - `migrations/001_create_messages_table.up.sql`: messagesテーブルの作成
   - すべてのマイグレーションロジックをSQLファイルで管理

2. マイグレーション管理用のクエリ（sqlc使用）:
   - `queries/schema_migrations.sql`にクエリを定義
   - `GetCurrentMigrationVersion`: 現在のバージョンを取得
   - `InsertMigrationVersion`: 新しいバージョンを記録
   - sqlcで型安全なGoコードを自動生成

3. マイグレーション実行ロジック:
   - `internal/platform/database/migrate.go`
   - pgxpoolを使ってマイグレーションSQLを直接実行
   - sqlcの生成コードでバージョン管理
   - トランザクション内で実行して原子性を保証

4. マイグレーションの仕組み:
   - `migrations/`ディレクトリから`.up.sql`ファイルを読み込み
   - ファイル名からバージョン番号を抽出（例: 001_create_messages_table.up.sql → 1）
   - sqlcで現在のバージョンを取得
   - 現在のバージョンより新しいマイグレーションのみ実行
   - 実行後、sqlcでバージョンを記録

5. マイグレーションファイルの命名規則:
   - `{version}_{description}.up.sql`形式
   - 例: `000_init_schema_migrations.up.sql`, `001_create_messages_table.up.sql`

6. main.goで起動時にマイグレーション実行:
   - pgxpool接続でマイグレーション実行
   - 同じpgxpool接続でsqlcのクエリ実行
   - 実行済みマイグレーションは自動でスキップ

メリット:
- すべてのSQLをファイルで管理（Goコードにハードコードしない）
- sqlcで型安全なバージョン管理
- CockroachDBで正常に動作（pg_advisory_lock()不要）
- 手動でマイグレーションSQLを実行する必要がなくなった
- マイグレーション履歴を自動管理
- 複数回起動しても安全（冪等性）
- トランザクションで原子性を保証
- golang-migrateの依存関係が不要

### アーキテクチャの変更点

sqlcの導入により、手動でリポジトリ実装を書く必要がなくなりました。

依存関係フロー:

```
cmd/api/main.go
  |
  v
adapter/line/webhook_handler.go
  | (依存)
  v
adapter/repository/db (sqlc自動生成)
```

sqlcが型安全なインターフェース(Querier)と実装を自動生成するため、domainレイヤーは不要になりました。

### 技術スタック

- データベース: CockroachDB (PostgreSQL互換)
- ドライバー: pgx v5
- 言語: Go 1.23
- コード生成: sqlc
- マイグレーション: golang-migrate/migrate
