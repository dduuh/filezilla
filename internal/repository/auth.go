package repository

import (
	"filezilla/internal/domain"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "users"
	filesTable = "files"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (a *AuthRepository) Create(user domain.User) (int, error) {
	var userId int
	query := fmt.Sprintf(`INSERT INTO %s (email, password) VALUES ($1, $2) RETURNING id`, usersTable)
	row := a.db.QueryRow(query, user.Email, user.Password)

	if err := row.Scan(&userId); err != nil {
		return -1, err
	}

	return userId, nil
}

func (a *AuthRepository) GetUser(email, password string) (int, error) {
	var userId int
	query := fmt.Sprintf(`SELECT id FROM %s WHERE email=$1 AND password=$2`, usersTable)
	err := a.db.Get(&userId, query, email, password)
	if err != nil {
		return -1, err
	}

	return userId, nil
}