package service

import (
	"time"

	"blazperic/radionica/internal/models"
	"blazperic/radionica/internal/repository"

	"github.com/google/uuid"
)

type NewsService struct {
	repo *repository.NewsRepository
}

func NewNewsService(repo *repository.NewsRepository) *NewsService {
	return &NewsService{repo: repo}
}

func (s *NewsService) GetAllNews() ([]*models.News, error) {
	return s.repo.GetAllNews()
}

func (s *NewsService) CreateNews(title, content string, userID uuid.UUID) (*models.News, error) {
	news := &models.News{
		ID:        uuid.New(),
		Title:     title,
		Content:   content,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateNews(news); err != nil {
		return nil, err
	}
	return news, nil
}
