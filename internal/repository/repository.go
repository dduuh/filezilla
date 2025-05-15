package repository

import (
	"context"
	"filezilla/internal/domain"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authentication interface {
	Create(user domain.User) (int, error)
	GetUser(email, password string) (int, error)
}

type Files interface {
	CreateFile(file domain.File) (int, error)
	GetFiles() ([]domain.File, error)
}

type Logs interface {
	Log(ctx context.Context, logInput *domain.Log) error
}

type Repository struct {
	Authentication
	Files
	Logs
}

func NewRepository(db *sqlx.DB, mongoDb *mongo.Database) *Repository {
	return &Repository{
		Authentication: NewAuthRepository(db),
		Files:          NewFileRepository(db),
		Logs:           NewLogRepository(mongoDb),
	}
}
