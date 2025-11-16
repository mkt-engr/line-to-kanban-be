package main

import (
	"fmt"
	"net/http"

	httpAdapter "line-to-kanban-be/internal/adapter/http"
	"line-to-kanban-be/internal/platform/config"
	"line-to-kanban-be/internal/platform/logger"
)

func main() {
	// 設定とロガーの初期化
	cfg := config.Load()
	log := logger.New()

	// ルーターの作成
	router := httpAdapter.NewRouter()

	// サーバー起動
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Info(fmt.Sprintf("Server starting on port %s", cfg.Port))

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Error(fmt.Sprintf("Server failed to start: %v", err))
	}
}
