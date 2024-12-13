package router

import (
	"lan-file-distributor/internal/handler"
	"lan-file-distributor/internal/repository"
	"lan-file-distributor/internal/service"
	"lan-file-distributor/pkg/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	// ミドルウェア
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 依存性の初期化
	fileRepo := repository.NewFileRepository(cfg.NasMountPath)
	fileService := service.NewFileService(fileRepo)
	fileHandler := handler.NewFileHandler(fileService)

	// ルート定義
	r.Get("/files/{folder}", fileHandler.ListFiles)

	return r
}
