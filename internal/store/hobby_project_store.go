package store

import "database/sql"

type HobbyProject struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PostgresHobbyProjectStore struct {
	db *sql.DB
}

func NewPostgresHobbyProjectStore(db *sql.DB) *PostgresHobbyProjectStore {
	return &PostgresHobbyProjectStore{db: db}
}

type HobbyProjectStore interface {
	GetHobbyProjectById(id int) (*HobbyProject, error)
	CreateHobbyProject(project *HobbyProject) (*HobbyProject, error)
	UpdateHobbyProject(projectId int, project *HobbyProject) (*HobbyProject, error)
	DeleteHobbyProject(projectId int) error
}

func (pg *PostgresHobbyProjectStore) GetHobbyProjectById(id int) (*HobbyProject, error) {
	project := &HobbyProject{}

	query := `SELECT id, name, description 
	FROM projects 
	WHERE id = $1 AND isDeleted = FALSE`

	err := pg.db.QueryRow(query, id).Scan(&project.Id, &project.Name, &project.Description)
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

	err = tx.QueryRow(query, project.Name, project.Description).Scan(&project.Id)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (pg *PostgresHobbyProjectStore) UpdateHobbyProject(projectId int, project *HobbyProject) (*HobbyProject, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `UPDATE projects 
	SET name = $1, description = $2 
	WHERE id = $3 AND isDeleted = FALSE`

	result, err := tx.Exec(query, project.Name, project.Description, projectId)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (pg *PostgresHobbyProjectStore) DeleteHobbyProject(projectId int) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `UPDATE projects 
	SET isDeleted = TRUE 
	WHERE id = $1 AND isDeleted = FALSE`

	result, err := tx.Exec(query, projectId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
