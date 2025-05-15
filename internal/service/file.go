package service

import (
	"filezilla/internal/domain"
	"filezilla/internal/repository"
)

type FileService struct {
	repo repository.Files
}

func NewFileService(repo repository.Files) *FileService {
	return &FileService{
		repo: repo,
	}
}

func (f *FileService) CreateFile(file domain.File) (int, error) {
	return f.repo.CreateFile(file)
}

func (f *FileService) GetFiles() ([]domain.File, error) {
	return f.repo.GetFiles()
}