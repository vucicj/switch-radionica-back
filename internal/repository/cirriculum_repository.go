package repository

import (
	"database/sql"

	"blazperic/radionica/internal/models"
)

type CirriculumRepository struct {
	db *sql.DB
}

func NewCirriculumRepository(db *sql.DB) *CirriculumRepository {
	return &CirriculumRepository{db: db}
}

func (r *CirriculumRepository) GetAllCirriculum() ([]*models.Cirriculum, error) {
	query := `
		SELECT id, title, week, description, user_id, created_at
		FROM cirriculum
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cirriculaList []*models.Cirriculum = make([]*models.Cirriculum, 0)
	for rows.Next() {
		cirriculum := &models.Cirriculum{}
		err := rows.Scan(&cirriculum.ID, &cirriculum.Title, &cirriculum.Week, &cirriculum.Content, &cirriculum.UserID, &cirriculum.CreatedAt)
		if err != nil {
			return nil, err
		}
		cirriculaList = append(cirriculaList, cirriculum)
	}
	return cirriculaList, nil
}

func (r *CirriculumRepository) CreateCirriculum(cirriculum *models.Cirriculum) error {
	query := `
		INSERT INTO cirriculum (id, title, week, description, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, cirriculum.ID, cirriculum.Title, cirriculum.Week, cirriculum.Content, cirriculum.UserID, cirriculum.CreatedAt)
	return err
}
