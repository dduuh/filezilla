package domain

import "time"

type File struct {
	Id         int       `json:"-" db:"id"`
	UserId     int       `json:"user_id" db:"user_id"`
	FileSize   int64     `json:"file_size" db:"file_size"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	StorageUrl string    `json:"storage_url" db:"storage_url"`
}
