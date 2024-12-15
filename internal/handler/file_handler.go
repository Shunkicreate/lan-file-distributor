package handler

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"image/jpeg"
	"lan-file-distributor/internal/repository"
	// "math/rand"
	"net/http"
	"strconv"
	// "time"
	"encoding/json"
	"fmt"
)

type FileHandler struct {
	fileRepo repository.FileRepository
}

func NewFileHandler(fileRepo repository.FileRepository) *FileHandler {
	return &FileHandler{fileRepo: fileRepo}
}

// ランダムな画像を取得
func (h *FileHandler) GetRandomImages(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))
	if count == 0 {
		count = 1
	}
	width, _ := strconv.ParseUint(r.URL.Query().Get("width"), 10, 32)
	if width == 0 {
		width = 1080 // スマートフォンの一般的な横幅
	}
	height, _ := strconv.ParseUint(r.URL.Query().Get("height"), 10, 32)
	if height == 0 {
		height = 1920/4 // スマートフォンの一般的な縦幅
	}
	folder := ""

	images, err := h.fileRepo.GetRandomFiles(folder, count, uint(width), uint(height))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting random files: %v", err), http.StatusInternalServerError)
		return
	}

	// 画像をJPEGとしてエンコードして返す
	var imageBuffers [][]byte
	for _, img := range images {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, img, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imageBuffers = append(imageBuffers, buf.Bytes())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"images": imageBuffers,
	})
}

// 特定の画像を取得
func (h *FileHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "path")
	width, _ := strconv.ParseUint(r.URL.Query().Get("width"), 10, 32)
	height, _ := strconv.ParseUint(r.URL.Query().Get("height"), 10, 32)

	img, err := h.fileRepo.GetFile(path, uint(width), uint(height))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// ファイル一覧を取得
func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	folder := ""
	files, err := h.fileRepo.ListFiles(folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(files)
}

// 複数の特定画像を取得
func (h *FileHandler) GetMultipleImages(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Paths  []string `json:"paths"`
		Width  uint     `json:"width"`
		Height uint     `json:"height"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	images, err := h.fileRepo.GetFiles(request.Paths, request.Width, request.Height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var imageBuffers [][]byte
	for _, img := range images {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, img, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imageBuffers = append(imageBuffers, buf.Bytes())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"images": imageBuffers,
	})
}

// ファイルパス一覧を取得
func (h *FileHandler) GetFilePaths(w http.ResponseWriter, r *http.Request) {
	folder := chi.URLParam(r, "folder")
	paths, err := h.fileRepo.GetFilePaths(folder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"paths": paths,
	})
}
