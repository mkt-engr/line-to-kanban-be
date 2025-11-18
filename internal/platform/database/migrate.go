package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations はマイグレーションを実行します
// CockroachDB向けにpgxpoolを使用してマイグレーションSQLを直接実行
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsPath string) error {
	// schema_migrationsテーブルを作成（存在しない場合）
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INT PRIMARY KEY,
			dirty BOOLEAN NOT NULL DEFAULT FALSE
		)
	`
	if _, err := pool.Exec(ctx, createTableSQL); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// 現在適用されているバージョンを取得
	var currentVersion int
	err := pool.QueryRow(ctx, "SELECT COALESCE(MAX(version), 0) FROM schema_migrations WHERE NOT dirty").Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
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
		if version <= currentVersion {
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

		// バージョンを記録
		if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version, dirty) VALUES ($1, FALSE)", version); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to record migration version: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return nil
}
