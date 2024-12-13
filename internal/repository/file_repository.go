package repository

import (
	"io/ioutil"
	"lan-file-distributor/internal/model"
)

type FileRepository interface {
	ListFiles(folder string) ([]model.File, error)
}

type fileRepository struct {
	basePath string
}

func NewFileRepository(basePath string) FileRepository {
	return &fileRepository{basePath: basePath}
}

func (r *fileRepository) ListFiles(folder string) ([]model.File, error) {
	files, err := ioutil.ReadDir(r.basePath + "/" + folder)
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
