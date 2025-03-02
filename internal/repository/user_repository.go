package repository

import (
	"blazperic/radionica/internal/models"
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
        INSERT INTO users (id, username, password, created_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.db.Exec(query, user.ID, user.Username, user.Password, user.CreatedAt)
	return err
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	query := `
        SELECT id, username, password, created_at
        FROM users
        WHERE username = $1
    `
	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
