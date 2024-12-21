package handler

import (
	"bytes"
	"github.com/go-chi/chi/v5"
	"image/jpeg"
	"lan-file-distributor/internal/repository"
	"net/http"
	"strconv"
	"encoding/json"
	"fmt"
)

type FileHandler struct {
	fileRepo repository.FileRepository
}

func NewFileHandler(fileRepo repository.FileRepository) *FileHandler {
	return &FileHandler{fileRepo: fileRepo}
}

func AdjustDimensions(originalWidth, originalHeight uint, width, height uint) (uint, uint) {

	if width == 0 && height == 0 {
		if originalWidth > originalHeight {
			return 1920, 1080
		} else {
			return 1080, 1920
		}
	} else if width == 0 {
		ratio := float64(height) / float64(originalHeight)
		return uint(float64(originalWidth) * ratio), height
	} else if height == 0 {
		ratio := float64(width) / float64(originalWidth)
		return width, uint(float64(originalHeight) * ratio)
	}
	return width, height
}

func (h *FileHandler) GetRandomImages(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.URL.Query().Get("count"))
	if count == 0 {
		count = 1
	}
	width, _ := strconv.ParseUint(r.URL.Query().Get("width"), 10, 32)
	height, _ := strconv.ParseUint(r.URL.Query().Get("height"), 10, 32)

	adjustedWidth, adjustedHeight := AdjustDimensions(1920, 1080, uint(width), uint(height))

	folder := ""

	imageFiles, err := h.fileRepo.GetRandomFiles(folder, count, adjustedWidth, adjustedHeight)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting random files: %v", err), http.StatusInternalServerError)
		return
	}

	// 画像データをエンコード
	for _, imgFile := range imageFiles {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, imgFile.Image, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imgFile.Data = buf.Bytes()
		imgFile.Image = nil // メモリ解放のため
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(imageFiles)
}

func (h *FileHandler) GetImage(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "path")
	width, _ := strconv.ParseUint(r.URL.Query().Get("width"), 10, 32)
	height, _ := strconv.ParseUint(r.URL.Query().Get("height"), 10, 32)

	adjustedWidth, adjustedHeight := AdjustDimensions(1920, 1080, uint(width), uint(height))

	imageFile, err := h.fileRepo.GetFile(path, adjustedWidth, adjustedHeight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, imageFile.Image, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(buf.Bytes())
}

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

	adjustedWidth, adjustedHeight := AdjustDimensions(1920, 1080, request.Width, request.Height)

	imageFiles, err := h.fileRepo.GetFiles(request.Paths, adjustedWidth, adjustedHeight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 画像データをエンコード
	for _, imgFile := range imageFiles {
		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, imgFile.Image, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imgFile.Data = buf.Bytes()
		imgFile.Image = nil // メモリ解放のため
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(imageFiles)
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
