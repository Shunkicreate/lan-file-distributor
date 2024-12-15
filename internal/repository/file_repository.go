package repository

import (
	"fmt"
	"github.com/nfnt/resize"
	"github.com/pixiv/go-libjpeg/jpeg"
	"image"
	_ "image/jpeg" // JPEGサポートのみ残す
	"io/ioutil"
	"lan-file-distributor/internal/model"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type FileRepository interface {
	ListFiles(folder string) ([]model.File, error)
	GetFile(path string, width, height uint) (image.Image, error)
	GetFiles(paths []string, width, height uint) ([]image.Image, error)
	GetFilePaths(folder string) ([]string, error)
	GetRandomFiles(folder string, count int, width, height uint) ([]image.Image, error)
}

type fileRepository struct {
	basePath string
}

// NewFileRepository creates a new instance of FileRepository
func NewFileRepository(basePath string) FileRepository {
	return &fileRepository{basePath: basePath}
}

// ListFiles retrieves a list of files in the specified folder.
func (r *fileRepository) ListFiles(folder string) ([]model.File, error) {
	fullPath := filepath.Join(r.basePath, defaultFolder(folder, "/nas"))
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var result []model.File
	for _, file := range files {
		result = append(result, model.File{
			Name: file.Name(),
			Path: filepath.Join(folder, file.Name()),
			Size: file.Size(),
		})
	}
	return result, nil
}

// GetFile retrieves and optionally resizes a single image file.
func (r *fileRepository) GetFile(path string, width, height uint) (image.Image, error) {
	if !isSupportedImage(path) {
		return nil, fmt.Errorf("unsupported image format (only JPG/JPEG supported): %s", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	options := &jpeg.DecoderOptions{
		DCTMethod:              jpeg.DCTIFast,
		DisableFancyUpsampling: true,
		DisableBlockSmoothing:  true,
		ScaleTarget: image.Rectangle{
			Max: image.Point{
				X: int(width),
				Y: int(height),
			},
		},
	}

	img, err := jpeg.Decode(file, options)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	if (width > 0 || height > 0) && 
	   (uint(img.Bounds().Dx()) != width || uint(img.Bounds().Dy()) != height) {
		return resize.Resize(width, height, img, resize.Lanczos3), nil
	}

	return img, nil
}

// GetFiles retrieves and optionally resizes multiple images concurrently.
func (r *fileRepository) GetFiles(paths []string, width, height uint) ([]image.Image, error) {
	var (
		images  = make([]image.Image, len(paths))
		errChan = make(chan error, len(paths))
		wg      sync.WaitGroup
	)

	for i, path := range paths {
		wg.Add(1)
		go func(index int, filePath string) {
			defer wg.Done()

			img, err := r.GetFile(filePath, width, height)
			if err != nil {
				errChan <- fmt.Errorf("failed to process file %s: %v", filePath, err)
				return
			}
			images[index] = img
		}(i, path)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return images, nil
}

// GetFilePaths retrieves all image file paths in a specified folder.
func (r *fileRepository) GetFilePaths(folder string) ([]string, error) {
	fullPath := filepath.Join(r.basePath, defaultFolder(folder, os.Getenv("NAS_PATH")))

	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var paths []string
	for _, file := range files {
		if !file.IsDir() && isSupportedImage(file.Name()) {
			fullFilePath := filepath.Join(fullPath, file.Name())
			paths = append(paths, fullFilePath)
		}
	}

	return paths, nil
}

// GetRandomFiles retrieves a specified number of random images.
func (r *fileRepository) GetRandomFiles(folder string, count int, width, height uint) ([]image.Image, error) {
	paths, err := r.GetFilePaths(folder)
	if err != nil {
		return nil, fmt.Errorf("failed to get file paths: %v", err)
	}

	if count > len(paths) {
		count = len(paths)
	}

	selectedPaths := randomSample(paths, count)
	images, err := r.GetFiles(selectedPaths, width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to get files: %v", err)
	}

	return images, nil
}

// Helper functions

// defaultFolder returns a default folder if none is specified.
func defaultFolder(folder, defaultPath string) string {
	if folder == "" {
		return defaultPath
	}
	return folder
}

// isSupportedImage checks if a file is a supported image format (JPG/JPEG only).
func isSupportedImage(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".jpg" || ext == ".jpeg"
}

// randomSample selects a random sample of file paths.
func randomSample(paths []string, count int) []string {
	perm := rand.Perm(len(paths))
	selected := make([]string, count)
	for i := 0; i < count; i++ {
		selected[i] = paths[perm[i]]
	}
	return selected
}
