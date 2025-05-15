package service

import (
	"context"
	"filezilla/internal/domain"
	"filezilla/internal/repository"
)

type LogService struct {
	repo repository.Logs
}

func NewLogService(repo repository.Logs) *LogService {
	return &LogService{
		repo: repo,
	}
}

func (l *LogService) Log(ctx context.Context, logInput *domain.Log) error {
	return l.repo.Log(ctx, logInput)
}
