package main

import (
	"log"
	"net/http"

	"lan-file-distributor/internal/router"
	"lan-file-distributor/pkg/config"
)

func main() {
	// 設定の読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// ルーターの初期化
	r := router.NewRouter(cfg)

	// サーバー起動
	log.Printf("Starting server on %s", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
