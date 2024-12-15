package repository

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	_ "image/gif"  // GIFサポートを追加
	_ "image/jpeg" // JPEGサポート
	_ "image/png"  // PNGサポート
	"io/ioutil"
	"lan-file-distributor/internal/model"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"sync"
)

type FileRepository interface {
	ListFiles(folder string) ([]model.File, error)
	GetFile(path string, width uint, height uint) (image.Image, error)
	GetFiles(paths []string, width uint, height uint) ([]image.Image, error)
	GetFilePaths(folder string) ([]string, error)
	GetRandomFiles(folder string, count int, width uint, height uint) ([]image.Image, error)
}

type fileRepository struct {
	basePath string
}

func NewFileRepository(basePath string) FileRepository {
	return &fileRepository{basePath: basePath}
}

func (r *fileRepository) ListFiles(folder string) ([]model.File, error) {
	if folder == "" {
		folder = "/nas"
	}
	fullPath := filepath.Join(r.basePath, folder)
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var result []model.File
	for _, f := range files {
		result = append(result, model.File{
			Name: f.Name(),
			Path: folder + "/" + f.Name(),
			Size: f.Size(),
		})
	}
	return result, nil
}

func (r *fileRepository) GetFile(path string, width uint, height uint) (image.Image, error) {
	fullPath := filepath.Join(r.basePath, path)
	
	// ファイル拡張子をチェック
	ext := strings.ToLower(filepath.Ext(fullPath))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	// ファイルを開く
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 画像をデコード
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image (%s format): %v", format, err)
	}

	// リサイズが指定されている場合
	if width > 0 || height > 0 {
		// リサイズ処理
		resizedImg := resize.Resize(width, height, img, resize.Lanczos3)
		return resizedImg, nil
	}

	return img, nil
}

func (r *fileRepository) GetFiles(paths []string, width uint, height uint) ([]image.Image, error) {
	var images = make([]image.Image, len(paths))
	errChan := make(chan error, len(paths))
	var wg sync.WaitGroup

	for i, path := range paths {
		wg.Add(1)
		go func(index int, filePath string) {
			defer wg.Done()
			img, err := r.GetFile(filePath, width, height)
			if err != nil {
				errChan <- err
				return
			}
			images[index] = img
		}(i, path)
	}

	wg.Wait()
	close(errChan)

	// エラーチェック
	if len(errChan) > 0 {
		return nil, <-errChan // 最初のエラーを返す
	}

	return images, nil
}

func (r *fileRepository) GetFilePaths(folder string) ([]string, error) {
	if folder == "" {
		folder = os.Getenv("NAS_PATH")
	}
	fullPath := filepath.Join(r.basePath, folder)
	
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, f := range files {
		if !f.IsDir() {
			ext := strings.ToLower(filepath.Ext(f.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" {
				path := filepath.Join(folder, f.Name())
				paths = append(paths, path)
			}
		}
	}

	return paths, nil
}

func (r *fileRepository) GetRandomFiles(folder string, count int, width uint, height uint) ([]image.Image, error) {
	// まず全てのファイルパスを取得
	paths, err := r.GetFilePaths(folder)
	if err != nil {
		return nil, err
	}

	// ファイル数が要求数より少ない場合は、利用可能な最大数に調整
	if len(paths) < count {
		count = len(paths)
	}

	// ランダムに指定枚数を選択
	rand.Seed(time.Now().UnixNano())
	selectedPaths := make([]string, count)
	perm := rand.Perm(len(paths))
	for i := 0; i < count; i++ {
		selectedPaths[i] = paths[perm[i]]
	}

	// 選択された画像を取得
	images, err := r.GetFiles(selectedPaths, width, height)
	if err != nil {
		return nil, err
	}

	return images, nil
}
