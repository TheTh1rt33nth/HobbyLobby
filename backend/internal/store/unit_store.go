package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Unit struct {
	Id             int       `json:"id"`
	ProjectId      int       `json:"projectId"`
	Name           string    `json:"name"`
	Quantity       int       `json:"quantity"`
	Status         string    `json:"status"`
	Notes          *string   `json:"notes"`
	PaintProfileId *int      `json:"paintProfileId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ProjectProgress struct {
	TotalUnits int            `json:"totalUnits"`
	ByStatus   map[string]int `json:"byStatus"`
}

type PostgresUnitStore struct {
	db *sql.DB
}

func NewPostgresUnitStore(db *sql.DB) *PostgresUnitStore {
	return &PostgresUnitStore{db: db}
}

type UnitStore interface {
	GetUnitById(ctx context.Context, id int) (*Unit, error)
	GetUnitsByProjectId(ctx context.Context, projectId int) ([]*Unit, error)
	CreateUnit(ctx context.Context, unit *Unit) (*Unit, error)
	UpdateUnit(ctx context.Context, unitId int, unit *Unit) (*Unit, error)
	DeleteUnit(ctx context.Context, unitId int) error
	GetProgressByProjectId(ctx context.Context, projectId int) (*ProjectProgress, error)
	GetProgressByProjectIds(ctx context.Context, projectIds []int) (map[int]*ProjectProgress, error)
}

func (pg *PostgresUnitStore) GetUnitById(ctx context.Context, id int) (*Unit, error) {
	unit := &Unit{}

	query := `SELECT id, project_id, name, quantity, status, notes, paint_profile_id, created_at, updated_at
	FROM units
	WHERE id = $1`

	err := pg.db.QueryRowContext(ctx, query, id).Scan(
		&unit.Id, &unit.ProjectId, &unit.Name, &unit.Quantity, &unit.Status,
		&unit.Notes, &unit.PaintProfileId, &unit.CreatedAt, &unit.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return unit, nil
}

func (pg *PostgresUnitStore) GetUnitsByProjectId(ctx context.Context, projectId int) ([]*Unit, error) {
	query := `SELECT id, project_id, name, quantity, status, notes, paint_profile_id, created_at, updated_at
	FROM units
	WHERE project_id = $1
	ORDER BY name`

	rows, err := pg.db.QueryContext(ctx, query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []*Unit
	for rows.Next() {
		unit := &Unit{}
		err = rows.Scan(
			&unit.Id, &unit.ProjectId, &unit.Name, &unit.Quantity, &unit.Status,
			&unit.Notes, &unit.PaintProfileId, &unit.CreatedAt, &unit.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		units = append(units, unit)
	}

	return units, rows.Err()
}

func (pg *PostgresUnitStore) CreateUnit(ctx context.Context, unit *Unit) (*Unit, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if unit.Status == "" {
		unit.Status = "unassembled"
	}
	if unit.Quantity == 0 {
		unit.Quantity = 1
	}

	query := `INSERT INTO units (project_id, name, quantity, status, notes, paint_profile_id)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query, unit.ProjectId, unit.Name, unit.Quantity, unit.Status, unit.Notes, unit.PaintProfileId).
		Scan(&unit.Id, &unit.CreatedAt, &unit.UpdatedAt)
	if err != nil {
		return nil, translatePgError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return unit, nil
}

func (pg *PostgresUnitStore) UpdateUnit(ctx context.Context, unitId int, unit *Unit) (*Unit, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `UPDATE units
	SET name = $1, quantity = $2, status = $3, notes = $4, paint_profile_id = $5, updated_at = NOW()
	WHERE id = $6
	RETURNING id, project_id, name, quantity, status, notes, paint_profile_id, created_at, updated_at`

	updated := &Unit{}
	err = tx.QueryRowContext(ctx, query, unit.Name, unit.Quantity, unit.Status, unit.Notes, unit.PaintProfileId, unitId).
		Scan(
			&updated.Id, &updated.ProjectId, &updated.Name, &updated.Quantity, &updated.Status,
			&updated.Notes, &updated.PaintProfileId, &updated.CreatedAt, &updated.UpdatedAt,
		)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, translatePgError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return updated, nil
}

func (pg *PostgresUnitStore) DeleteUnit(ctx context.Context, unitId int) error {
	result, err := pg.db.ExecContext(ctx, `DELETE FROM units WHERE id = $1`, unitId)
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

	return nil
}

func (pg *PostgresUnitStore) GetProgressByProjectId(ctx context.Context, projectId int) (*ProjectProgress, error) {
	query := `SELECT status, SUM(quantity) FROM units WHERE project_id = $1 GROUP BY status`

	rows, err := pg.db.QueryContext(ctx, query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	progress := &ProjectProgress{ByStatus: make(map[string]int)}
	for rows.Next() {
		var status string
		var count int
		if err = rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		progress.ByStatus[status] = count
		progress.TotalUnits += count
	}

	return progress, rows.Err()
}

func (pg *PostgresUnitStore) GetProgressByProjectIds(ctx context.Context, projectIds []int) (map[int]*ProjectProgress, error) {
	if len(projectIds) == 0 {
		return map[int]*ProjectProgress{}, nil
	}

	args := make([]any, len(projectIds))
	placeholders := make([]string, len(projectIds))
	for i, id := range projectIds {
		args[i] = id
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(
		`SELECT project_id, status, SUM(quantity) FROM units WHERE project_id IN (%s) GROUP BY project_id, status`,
		strings.Join(placeholders, ","),
	)

	rows, err := pg.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int]*ProjectProgress)
	for rows.Next() {
		var projectId int
		var status string
		var count int
		if err = rows.Scan(&projectId, &status, &count); err != nil {
			return nil, err
		}
		if result[projectId] == nil {
			result[projectId] = &ProjectProgress{ByStatus: make(map[string]int)}
		}
		result[projectId].ByStatus[status] = count
		result[projectId].TotalUnits += count
	}

	return result, rows.Err()
}
