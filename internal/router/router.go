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
	r.Route("/api/files", func(r chi.Router) {
		// ディレクトリ内のファイル一覧を取得
		r.Get("/list/*", fileHandler.ListFiles)
		
		// 指定枚数のランダムな画像を取得 
		r.Get("/random", fileHandler.GetRandomImages)
		
		// 特定の画像を取得（リサイズオプション付き）
		r.Get("/image/*", fileHandler.GetImage)
		
		// 複数の特定画像を取得（リサイズオプション付き）
		r.Post("/batch", fileHandler.GetMultipleImages)
		
		// ディレクトリ内の全ファイルパスを取得
		r.Get("/paths/*", fileHandler.GetFilePaths)
	})

	return r
}
