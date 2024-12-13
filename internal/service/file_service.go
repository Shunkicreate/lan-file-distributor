package service

import (
	"lan-file-distributor/internal/model"
	"lan-file-distributor/internal/repository"
)

type FileService interface {
	GetFiles(folder string) ([]model.File, error)
}

type fileService struct {
	fileRepo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) FileService {
	return &fileService{fileRepo: repo}
}

func (s *fileService) GetFiles(folder string) ([]model.File, error) {
	return s.fileRepo.ListFiles(folder)
}
