package repository

import (
	"database/sql"

	"blazperic/radionica/internal/models"
)

type NewsRepository struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

func (r *NewsRepository) GetAllNews() ([]*models.News, error) {
	query := `
		SELECT id, title, content, user_id, created_at
		FROM news
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []*models.News
	for rows.Next() {
		news := &models.News{}
		err := rows.Scan(&news.ID, &news.Title, &news.Content, &news.UserID, &news.CreatedAt)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, news)
	}
	return newsList, nil
}

func (r *NewsRepository) CreateNews(news *models.News) error {
	query := `
		INSERT INTO news (id, title, content, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, news.ID, news.Title, news.Content, news.UserID, news.CreatedAt)
	return err
}
