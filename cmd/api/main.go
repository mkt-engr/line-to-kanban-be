package main

import (
	"context"
	// "database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
	// _ "github.com/lib/pq"

	httpAdapter "line-to-kanban-be/internal/adapter/http"
	lineAdapter "line-to-kanban-be/internal/adapter/line"
	"line-to-kanban-be/internal/adapter/repository/db"
	"line-to-kanban-be/internal/platform/config"
	// "line-to-kanban-be/internal/platform/database"
	"line-to-kanban-be/internal/platform/logger"
)

func main() {
	// .envファイルの読み込み（存在しない場合はスキップ）
	_ = godotenv.Load()

	// 設定とロガーの初期化
	cfg := config.Load()
	lineConfig := config.LoadLineConfig()
	dbConfig := config.LoadDatabaseConfig()
	appLogger := logger.New()

	// LINE Botクライアントの初期化（環境変数が必須）
	if lineConfig.ChannelSecret == "" || lineConfig.ChannelAccessToken == "" {
		log.Fatal("LINE credentials are required. Please set LINE_CHANNEL_SECRET and LINE_CHANNEL_ACCESS_TOKEN environment variables.")
	}

	lineClient, err := lineAdapter.NewClient(lineConfig)
	if err != nil {
		log.Fatalf("LINE Botクライアントの初期化に失敗しました: %v", err)
	}
	appLogger.Info("LINE Bot client initialized successfully")

	// データベース接続の初期化
	if dbConfig.URL == "" {
		log.Fatal("Database URL is required. Please set DATABASE_URL environment variable.")
	}

	// TODO: CockroachDBのマイグレーション対応（pg_advisory_lock()の問題を解決後に有効化）
	// マイグレーション用にdatabase/sql接続を作成
	// sqlDB, err := sql.Open("postgres", dbConfig.URL)
	// if err != nil {
	// 	log.Fatalf("データベース接続に失敗しました: %v", err)
	// }
	// defer sqlDB.Close()

	// マイグレーションを実行
	// appLogger.Info("Running database migrations...")
	// if err := database.RunMigrations(sqlDB, "migrations"); err != nil {
	// 	log.Fatalf("マイグレーション実行に失敗しました: %v", err)
	// }
	// appLogger.Info("Database migrations completed successfully")

	// sqlc用にpgxpool接続を作成
	ctx := context.Background()
	dbPool, err := pgxpool.New(ctx, dbConfig.URL)
	if err != nil {
		log.Fatalf("データベース接続プールの作成に失敗しました: %v", err)
	}
	defer dbPool.Close()

	// データベース接続テスト
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("データベースへのPingに失敗しました: %v", err)
	}
	appLogger.Info("Database connection established successfully")

	// sqlcで生成されたQueriesの初期化
	queries := db.New(dbPool)

	// LINE Webhookハンドラーの初期化
	lineWebhookHandler := lineAdapter.NewWebhookHandler(lineClient, queries)
	appLogger.Info("Webhook handler initialized successfully")

	// ルーターの作成
	router := httpAdapter.NewRouter(lineWebhookHandler)

	// サーバー起動
	addr := fmt.Sprintf(":%s", cfg.Port)
	appLogger.Info(fmt.Sprintf("Server starting on port %s", cfg.Port))

	if err := http.ListenAndServe(addr, router); err != nil {
		appLogger.Error(fmt.Sprintf("Server failed to start: %v", err))
	}
}
