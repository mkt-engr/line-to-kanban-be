# Database Connection

## ライブラリ

### sqlc

公式サイト: https://sqlc.dev/

概要:
sqlc は SQL からタイプセーフな Go コードを自動生成するツールです。

主な特徴:
- SQL ファイルを書くと、Go のコードを自動生成
- コンパイル時に型チェックが可能
- ORM ではなく、生の SQL を使用
- PostgreSQL、MySQL、SQLite に対応

使い方:
1. SQL クエリを `.sql` ファイルに記述
2. `sqlc generate` コマンドで Go コードを生成
3. 生成されたコードを使ってデータベース操作

例:
```sql
-- queries/users.sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC;
```

sqlc が以下のような Go コードを自動生成:
```go
func (q *Queries) GetUser(ctx context.Context, id int64) (User, error)
func (q *Queries) ListUsers(ctx context.Context) ([]User, error)
```

メリット:
- タイプセーフ: コンパイル時にエラーを検出
- パフォーマンス: 生の SQL なので高速
- 学習コスト低: SQL を書くだけ
- メンテナンス性: SQL とコードが分離

---

### pgx

公式サイト: https://github.com/jackc/pgx

概要:
pgx は PostgreSQL 専用の高性能な Go ドライバーです。

主な特徴:
- PostgreSQL の標準ドライバー `database/sql` よりも高速
- PostgreSQL 固有の機能をフルサポート
- 接続プーリング機能を内蔵
- プリペアドステートメント、バッチ処理に対応

使い方:
```go
import "github.com/jackc/pgx/v5/pgxpool"

// 接続プールの作成
pool, err := pgxpool.New(context.Background(), "postgres://user:pass@localhost:5432/dbname")
if err != nil {
    log.Fatal(err)
}
defer pool.Close()

// クエリ実行
var name string
err = pool.QueryRow(context.Background(), "SELECT name FROM users WHERE id=$1", 123).Scan(&name)
```

メリット:
- 高速: `database/sql` より約 2 倍高速
- PostgreSQL 特化: JSON、配列、UUID など PostgreSQL の型を完全サポート
- 接続プール: 自動的にコネクションを管理
- エラーハンドリング: PostgreSQL のエラーを詳細に取得可能

pgx/v5 と pgx/v4 の違い:
- v5 は Go 1.19 以降が必要
- v5 はジェネリクス対応でより型安全
- v5 はパフォーマンスが向上

---

## sqlc と pgx の組み合わせ

sqlc は pgx ドライバーと組み合わせて使うことができます。

設定例 (sqlc.yaml):
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "queries/"
    schema: "schema.sql"
    gen:
      go:
        package: "db"
        out: "internal/adapter/db"
        sql_package: "pgx/v5"
```

この設定により、sqlc が生成するコードが pgx/v5 を使用するようになります。

メリット:
- sqlc のタイプセーフな API
- pgx の高速なパフォーマンス
- PostgreSQL の機能をフルに活用

---

## このプロジェクトでの使用予定

- データベース: PostgreSQL (または開発時は SQLite)
- ドライバー: pgx/v5
- コード生成: sqlc
- 配置:
  - SQL ファイル: `queries/`
  - スキーマ: `migrations/`
  - 生成コード: `internal/adapter/db/`

---

## 参考リンク

- [sqlc Documentation](https://docs.sqlc.dev/)
- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [sqlc + pgx Tutorial](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html)
