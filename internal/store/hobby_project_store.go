package store

import "database/sql"

type HobbyProject struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PostgresHobbyProjectStore struct {
	db *sql.DB
}

func NewPostgresHobbyProjectRepository(db *sql.DB) *PostgresHobbyProjectStore {
	return &PostgresHobbyProjectStore{db: db}
}

type HobbyProjectStore interface {
	GetHobbyProjectById(id int) (*HobbyProject, error)
}

func (pg *PostgresHobbyProjectStore) GetHobbyProjectById(id int) (*HobbyProject, error) {
	project := &HobbyProject{}

	query := `SELECT id, name, description FROM projects WHERE id = $1`

	err := pg.db.QueryRow(query, id).Scan(&project.ID, &project.Name, &project.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return project, nil
}
