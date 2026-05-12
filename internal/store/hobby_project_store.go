package store

import (
	"database/sql"
	"time"
)

type HobbyProject struct {
	Id          int       `json:"id"`
	UserId      int       `json:"userId"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	GameSystem  *string   `json:"gameSystem"`
	Faction     *string   `json:"faction"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PostgresHobbyProjectStore struct {
	db *sql.DB
}

func NewPostgresHobbyProjectStore(db *sql.DB) *PostgresHobbyProjectStore {
	return &PostgresHobbyProjectStore{db: db}
}

type HobbyProjectStore interface {
	GetHobbyProjectById(id int) (*HobbyProject, error)
	GetHobbyProjectsByUserId(userId int) ([]*HobbyProject, error)
	CreateHobbyProject(project *HobbyProject) (*HobbyProject, error)
	UpdateHobbyProject(projectId int, project *HobbyProject) (*HobbyProject, error)
	DeleteHobbyProject(projectId int) error
}

func (pg *PostgresHobbyProjectStore) GetHobbyProjectById(id int) (*HobbyProject, error) {
	project := &HobbyProject{}

	query := `SELECT id, user_id, name, description, game_system, faction, status, created_at, updated_at
	FROM projects
	WHERE id = $1 AND is_deleted = FALSE`

	err := pg.db.QueryRow(query, id).Scan(
		&project.Id, &project.UserId, &project.Name, &project.Description,
		&project.GameSystem, &project.Faction, &project.Status,
		&project.CreatedAt, &project.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (pg *PostgresHobbyProjectStore) GetHobbyProjectsByUserId(userId int) ([]*HobbyProject, error) {
	query := `SELECT id, user_id, name, description, game_system, faction, status, created_at, updated_at
	FROM projects
	WHERE user_id = $1 AND is_deleted = FALSE
	ORDER BY created_at DESC`

	rows, err := pg.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*HobbyProject
	for rows.Next() {
		project := &HobbyProject{}
		err = rows.Scan(
			&project.Id, &project.UserId, &project.Name, &project.Description,
			&project.GameSystem, &project.Faction, &project.Status,
			&project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, rows.Err()
}

func (pg *PostgresHobbyProjectStore) CreateHobbyProject(project *HobbyProject) (*HobbyProject, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if project.Status == "" {
		project.Status = "unassembled"
	}

	query := `INSERT INTO projects (user_id, name, description, game_system, faction, status)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at`

	err = tx.QueryRow(query, project.UserId, project.Name, project.Description, project.GameSystem, project.Faction, project.Status).
		Scan(&project.Id, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return nil, translatePgError(err)
	}

	if err = tx.Commit(); err != nil {
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
	SET name = $1, description = $2, game_system = $3, faction = $4, status = $5, updated_at = NOW()
	WHERE id = $6 AND is_deleted = FALSE
	RETURNING id, user_id, name, description, game_system, faction, status, created_at, updated_at`

	updated := &HobbyProject{}
	err = tx.QueryRow(query, project.Name, project.Description, project.GameSystem, project.Faction, project.Status, projectId).
		Scan(
			&updated.Id, &updated.UserId, &updated.Name, &updated.Description,
			&updated.GameSystem, &updated.Faction, &updated.Status,
			&updated.CreatedAt, &updated.UpdatedAt,
		)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return updated, nil
}

func (pg *PostgresHobbyProjectStore) DeleteHobbyProject(projectId int) error {
	tx, err := pg.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `UPDATE projects 
	SET is_deleted = TRUE, updated_at = NOW() 
	WHERE id = $1 AND is_deleted = FALSE`

	result, err := tx.Exec(query, projectId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
