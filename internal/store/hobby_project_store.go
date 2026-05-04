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
	CreateHobbyProject(project *HobbyProject) (*HobbyProject, error)
	UpdateHobbyProject(project *HobbyProject) (*HobbyProject, error)
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

func (pg *PostgresHobbyProjectStore) CreateHobbyProject(project *HobbyProject) (*HobbyProject, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `INSERT INTO projects (name, description) 
	VALUES ($1, $2) 
	RETURNING id`

	err = tx.QueryRow(query, project.Name, project.Description).Scan(&project.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (pg *PostgresHobbyProjectStore) UpdateHobbyProject(project *HobbyProject) (*HobbyProject, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `UPDATE projects 
	SET name = $1, description = $2 
	WHERE id = $3`

	_, err = tx.Exec(query, project.Name, project.Description, project.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return project, nil
}
