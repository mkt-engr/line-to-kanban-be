package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"line-to-kanban-be/internal/adapter/repository/db"
)

// RunMigrations はマイグレーションを実行します
// CockroachDB向けにpgxpoolを使用してマイグレーションSQLを直接実行
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsPath string) error {
	queries := db.New(pool)

	// 現在適用されているバージョンを取得
	// schema_migrationsテーブルが存在しない場合は-1を返す
	currentVersion, err := queries.GetCurrentMigrationVersion(ctx)
	if err != nil {
		// テーブルが存在しない場合は、バージョン-1から開始（000_init_schema_migrations.up.sqlを実行するため）
		currentVersion = -1
	}

	// マイグレーションファイルを読み込み
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// .upファイルのみを抽出してソート
	var upFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".up.sql") {
			upFiles = append(upFiles, file.Name())
		}
	}
	sort.Strings(upFiles)

	// 各マイグレーションファイルを実行
	for _, fileName := range upFiles {
		// ファイル名からバージョン番号を抽出（例: 001_create_messages_table.up.sql -> 1）
		versionStr := strings.Split(fileName, "_")[0]
		var version int
		if _, err := fmt.Sscanf(versionStr, "%d", &version); err != nil {
			continue // バージョン番号が取得できない場合はスキップ
		}

		// 既に適用済みの場合はスキップ
		if int32(version) <= currentVersion {
			continue
		}

		// マイグレーションファイルを読み込み
		filePath := filepath.Join(migrationsPath, fileName)
		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", fileName, err)
		}

		// トランザクション内でマイグレーションを実行
		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// マイグレーションSQLを実行
		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to execute migration %s: %w", fileName, err)
		}

		// バージョンを記録（sqlcの生成コードを使用）
		txQueries := queries.WithTx(tx)
		if err := txQueries.InsertMigrationVersion(ctx, int32(version)); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to record migration version: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return nil
}
