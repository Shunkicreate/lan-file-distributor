package service

import (
	"lan-file-distributor/internal/model"
	"lan-file-distributor/internal/repository"
)

type FileService struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) *FileService {
	return &FileService{repo: repo}
}

func (s *FileService) ListFiles(folder string) ([]model.File, error) {
	return s.repo.ListFiles(folder)
}

func (s *FileService) GetFile(path string, width uint, height uint) (*model.ImageFile, error) {
	return s.repo.GetFile(path, width, height)
}

func (s *FileService) GetFiles(paths []string, width uint, height uint) ([]*model.ImageFile, error) {
	return s.repo.GetFiles(paths, width, height)
}

func (s *FileService) GetFilePaths(folder string) ([]string, error) {
	return s.repo.GetFilePaths(folder)
}

func (s *FileService) GetRandomFiles(folder string, count int, width uint, height uint) ([]*model.ImageFile, error) {
	return s.repo.GetRandomFiles(folder, count, width, height)
}
