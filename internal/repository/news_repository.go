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
		SELECT id, title, content,image_path,category, user_id, created_at
		FROM news
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []*models.News = make([]*models.News, 0)
	for rows.Next() {
		news := &models.News{}
		err := rows.Scan(&news.ID, &news.Title, &news.Content, &news.ImagePath, &news.Category, &news.UserID, &news.CreatedAt)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, news)
	}
	return newsList, nil
}

func (r *NewsRepository) CreateNews(news *models.News) error {
	query := `
		INSERT INTO news (id, title, content,image_path,category, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, news.ID, news.Title, news.Content, news.ImagePath, news.Category, news.UserID, news.CreatedAt)
	return err
}
