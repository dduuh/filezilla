package repository

import (
	"filezilla/internal/domain"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type FileRepository struct {
	db *sqlx.DB
}

func NewFileRepository(db *sqlx.DB) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func (f *FileRepository) CreateFile(file domain.File) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s
	(file_size, created_at, user_id, storage_url) VALUES ($1, $2, $3, $4) RETURNING id`,
		filesTable)

	row := f.db.QueryRow(query, file.FileSize, file.CreatedAt, file.UserId, file.StorageUrl)
	err := row.Scan(&file.Id)
	if err != nil {
		return -1, nil
	}

	return file.Id, nil
}

func (f *FileRepository) GetFiles() ([]domain.File, error) {
	var files []domain.File
	query := fmt.Sprintf(`SELECT file_size, created_at, user_id, storage_url FROM %s`, filesTable)
	err := f.db.Select(&files, query)
	if err != nil {
		return nil, err
	}

	return files, nil
}