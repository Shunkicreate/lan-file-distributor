package handler

import (
	"encoding/json"
	"lan-file-distributor/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(fs service.FileService) *FileHandler {
	return &FileHandler{fileService: fs}
}

func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	folder := chi.URLParam(r, "folder")
	files, err := h.fileService.GetFiles(folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSONレスポンス
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}
