package service

import (
	"time"

	"blazperic/radionica/internal/models"
	"blazperic/radionica/internal/repository"

	"github.com/google/uuid"
)

type CirriculumService struct {
	repo *repository.CirriculumRepository
}

func NewCirriculumService(repo *repository.CirriculumRepository) *CirriculumService {
	return &CirriculumService{repo: repo}
}

func (s *CirriculumService) GetAllCirriculum() ([]*models.Cirriculum, error) {
	return s.repo.GetAllCirriculum()
}

func (s *CirriculumService) CreateCirriculum(title, description string, week int, userID uuid.UUID) (*models.Cirriculum, error) {
	cirriculum := &models.Cirriculum{
		ID:        uuid.New(),
		Title:     title,
		Week:      week,
		Content:   description,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateCirriculum(cirriculum); err != nil {
		return nil, err
	}
	return cirriculum, nil
}
