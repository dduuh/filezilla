package repository

import (
	"context"
	"filezilla/internal/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

type LogRepository struct {
	db *mongo.Database
}

func NewLogRepository(db *mongo.Database) *LogRepository {
	return &LogRepository{
		db: db,
	}
}

func (l *LogRepository) Log(ctx context.Context, logInput *domain.Log) error {
	if _, err := l.db.Collection("logs").InsertOne(ctx, &logInput); err != nil {
		return err
	}
	return nil
}
