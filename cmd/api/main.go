package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	httpAdapter "line-to-kanban-be/internal/adapter/http"
	lineAdapter "line-to-kanban-be/internal/adapter/line"
	"line-to-kanban-be/internal/platform/config"
	"line-to-kanban-be/internal/platform/logger"
)

func main() {
	// .envファイルの読み込み（存在しない場合はスキップ）
	_ = godotenv.Load()

	// 設定とロガーの初期化
	cfg := config.Load()
	lineConfig := config.LoadLineConfig()
	appLogger := logger.New()

	// LINE Botクライアントの初期化（環境変数が必須）
	if lineConfig.ChannelSecret == "" || lineConfig.ChannelAccessToken == "" {
		log.Fatal("LINE credentials are required. Please set LINE_CHANNEL_SECRET and LINE_CHANNEL_ACCESS_TOKEN environment variables.")
	}

	lineClient, err := lineAdapter.NewClient(lineConfig)
	if err != nil {
		log.Fatalf("LINE Botクライアントの初期化に失敗しました: %v", err)
	}
	lineWebhookHandler := lineAdapter.NewWebhookHandler(lineClient)
	appLogger.Info("LINE Bot client initialized successfully")

	// ルーターの作成
	router := httpAdapter.NewRouter(lineWebhookHandler)

	// サーバー起動
	addr := fmt.Sprintf(":%s", cfg.Port)
	appLogger.Info(fmt.Sprintf("Server starting on port %s", cfg.Port))

	if err := http.ListenAndServe(addr, router); err != nil {
		appLogger.Error(fmt.Sprintf("Server failed to start: %v", err))
	}
}
