package store

import (
	"context"
	"database/sql"
	"time"
)

type PaintProfile struct {
	Id          int         `json:"id"`
	UserId      int         `json:"userId"`
	Name        string      `json:"name"`
	Description *string     `json:"description"`
	TargetArea  *string     `json:"targetArea"`
	Steps       []PaintStep `json:"steps,omitempty"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type PaintStep struct {
	Id                int     `json:"id"`
	PaintProfileId    int     `json:"paintProfileId"`
	StepOrder         int     `json:"stepOrder"`
	PaintName         string  `json:"paintName"`
	Brand             *string `json:"brand"`
	PaintType         *string `json:"paintType"`
	ApplicationMethod *string `json:"applicationMethod"`
	ColorHex          *string `json:"colorHex"`
	Notes             *string `json:"notes"`
}

type PostgresPaintProfileStore struct {
	db *sql.DB
}

func NewPostgresPaintProfileStore(db *sql.DB) *PostgresPaintProfileStore {
	return &PostgresPaintProfileStore{db: db}
}

type PaintProfileStore interface {
	GetPaintProfileById(ctx context.Context, id int) (*PaintProfile, error)
	GetPaintProfilesByUserId(ctx context.Context, userId int) ([]*PaintProfile, error)
	GetPaintProfilesByProjectId(ctx context.Context, projectId int) ([]*PaintProfile, error)
	CreatePaintProfile(ctx context.Context, profile *PaintProfile) (*PaintProfile, error)
	UpdatePaintProfile(ctx context.Context, profileId int, profile *PaintProfile) (*PaintProfile, error)
	DeletePaintProfile(ctx context.Context, profileId int) error

	CreatePaintStep(ctx context.Context, step *PaintStep) (*PaintStep, error)
	UpdatePaintStep(ctx context.Context, stepId int, step *PaintStep) (*PaintStep, error)
	DeletePaintStep(ctx context.Context, stepId int) error

	AssignToProject(ctx context.Context, projectId, profileId int) error
	UnassignFromProject(ctx context.Context, projectId, profileId int) error
}

func (pg *PostgresPaintProfileStore) GetPaintProfileById(ctx context.Context, id int) (*PaintProfile, error) {
	profile := &PaintProfile{}

	query := `SELECT id, user_id, name, description, target_area, created_at, updated_at
	FROM paint_profiles
	WHERE id = $1 AND is_deleted = FALSE`

	err := pg.db.QueryRowContext(ctx, query, id).Scan(
		&profile.Id, &profile.UserId, &profile.Name, &profile.Description, &profile.TargetArea,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	stepsQuery := `SELECT id, step_order, paint_name, brand, paint_type, application_method, color_hex, notes
	FROM paint_steps
	WHERE paint_profile_id = $1
	ORDER BY step_order`

	rows, err := pg.db.QueryContext(ctx, stepsQuery, profile.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		step := PaintStep{PaintProfileId: profile.Id}
		err = rows.Scan(&step.Id, &step.StepOrder, &step.PaintName, &step.Brand, &step.PaintType, &step.ApplicationMethod, &step.ColorHex, &step.Notes)
		if err != nil {
			return nil, err
		}
		profile.Steps = append(profile.Steps, step)
	}

	return profile, rows.Err()
}

func (pg *PostgresPaintProfileStore) loadStepsForProfile(ctx context.Context, profile *PaintProfile) error {
	stepsQuery := `SELECT id, step_order, paint_name, brand, paint_type, application_method, color_hex, notes
	FROM paint_steps
	WHERE paint_profile_id = $1
	ORDER BY step_order`

	stepRows, err := pg.db.QueryContext(ctx, stepsQuery, profile.Id)
	if err != nil {
		return err
	}
	defer stepRows.Close()

	for stepRows.Next() {
		step := PaintStep{PaintProfileId: profile.Id}
		err = stepRows.Scan(&step.Id, &step.StepOrder, &step.PaintName, &step.Brand, &step.PaintType, &step.ApplicationMethod, &step.ColorHex, &step.Notes)
		if err != nil {
			return err
		}
		profile.Steps = append(profile.Steps, step)
	}
	return stepRows.Err()
}

func (pg *PostgresPaintProfileStore) GetPaintProfilesByUserId(ctx context.Context, userId int) ([]*PaintProfile, error) {
	query := `SELECT id, user_id, name, description, target_area, created_at, updated_at
	FROM paint_profiles
	WHERE user_id = $1 AND is_deleted = FALSE
	ORDER BY name`

	rows, err := pg.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*PaintProfile
	for rows.Next() {
		profile := &PaintProfile{}
		err = rows.Scan(&profile.Id, &profile.UserId, &profile.Name, &profile.Description, &profile.TargetArea, &profile.CreatedAt, &profile.UpdatedAt)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, profile := range profiles {
		if err = pg.loadStepsForProfile(ctx, profile); err != nil {
			return nil, err
		}
	}

	return profiles, nil
}

func (pg *PostgresPaintProfileStore) GetPaintProfilesByProjectId(ctx context.Context, projectId int) ([]*PaintProfile, error) {
	query := `SELECT pp.id, pp.user_id, pp.name, pp.description, pp.target_area, pp.created_at, pp.updated_at
	FROM paint_profiles pp
	INNER JOIN project_paint_profiles ppp ON pp.id = ppp.paint_profile_id
	WHERE ppp.project_id = $1 AND pp.is_deleted = FALSE
	ORDER BY pp.name`

	rows, err := pg.db.QueryContext(ctx, query, projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*PaintProfile
	for rows.Next() {
		profile := &PaintProfile{}
		err = rows.Scan(&profile.Id, &profile.UserId, &profile.Name, &profile.Description, &profile.TargetArea, &profile.CreatedAt, &profile.UpdatedAt)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, profile)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, profile := range profiles {
		if err = pg.loadStepsForProfile(ctx, profile); err != nil {
			return nil, err
		}
	}

	return profiles, nil
}

func (pg *PostgresPaintProfileStore) CreatePaintProfile(ctx context.Context, profile *PaintProfile) (*PaintProfile, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `INSERT INTO paint_profiles (user_id, name, description, target_area)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query, profile.UserId, profile.Name, profile.Description, profile.TargetArea).
		Scan(&profile.Id, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return profile, nil
}

func (pg *PostgresPaintProfileStore) UpdatePaintProfile(ctx context.Context, profileId int, profile *PaintProfile) (*PaintProfile, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `UPDATE paint_profiles
	SET name = $1, description = $2, target_area = $3, updated_at = NOW()
	WHERE id = $4 AND is_deleted = FALSE
	RETURNING id, user_id, name, description, target_area, created_at, updated_at`

	updated := &PaintProfile{}
	err = tx.QueryRowContext(ctx, query, profile.Name, profile.Description, profile.TargetArea, profileId).
		Scan(&updated.Id, &updated.UserId, &updated.Name, &updated.Description, &updated.TargetArea, &updated.CreatedAt, &updated.UpdatedAt)
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

func (pg *PostgresPaintProfileStore) DeletePaintProfile(ctx context.Context, profileId int) error {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `UPDATE paint_profiles SET is_deleted = TRUE, updated_at = NOW() WHERE id = $1 AND is_deleted = FALSE`

	result, err := tx.ExecContext(ctx, query, profileId)
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

	return tx.Commit()
}

func (pg *PostgresPaintProfileStore) CreatePaintStep(ctx context.Context, step *PaintStep) (*PaintStep, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Auto-assign step_order as max + 1 if not provided.
	// An advisory lock on the profile ID serializes concurrent auto-assigns,
	// preventing two requests from racing to read the same MAX and colliding on the UNIQUE constraint.
	if step.StepOrder == 0 {
		if _, err = tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1)`, step.PaintProfileId); err != nil {
			return nil, err
		}
		var maxOrder int
		err = tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(step_order), 0) FROM paint_steps WHERE paint_profile_id = $1`, step.PaintProfileId).Scan(&maxOrder)
		if err != nil {
			return nil, err
		}
		step.StepOrder = maxOrder + 1
	}

	query := `INSERT INTO paint_steps (paint_profile_id, step_order, paint_name, brand, paint_type, application_method, color_hex, notes)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id`

	err = tx.QueryRowContext(ctx, query, step.PaintProfileId, step.StepOrder, step.PaintName, step.Brand, step.PaintType, step.ApplicationMethod, step.ColorHex, step.Notes).
		Scan(&step.Id)
	if err != nil {
		return nil, translatePgError(err)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return step, nil
}

func (pg *PostgresPaintProfileStore) UpdatePaintStep(ctx context.Context, stepId int, step *PaintStep) (*PaintStep, error) {
	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `UPDATE paint_steps
	SET step_order = $1, paint_name = $2, brand = $3, paint_type = $4, application_method = $5, color_hex = $6, notes = $7, updated_at = NOW()
	WHERE id = $8
	RETURNING id, paint_profile_id, step_order, paint_name, brand, paint_type, application_method, color_hex, notes`

	updated := &PaintStep{}
	err = tx.QueryRowContext(ctx, query, step.StepOrder, step.PaintName, step.Brand, step.PaintType, step.ApplicationMethod, step.ColorHex, step.Notes, stepId).
		Scan(&updated.Id, &updated.PaintProfileId, &updated.StepOrder, &updated.PaintName, &updated.Brand, &updated.PaintType, &updated.ApplicationMethod, &updated.ColorHex, &updated.Notes)
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

func (pg *PostgresPaintProfileStore) DeletePaintStep(ctx context.Context, stepId int) error {
	result, err := pg.db.ExecContext(ctx, `DELETE FROM paint_steps WHERE id = $1`, stepId)
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

func (pg *PostgresPaintProfileStore) AssignToProject(ctx context.Context, projectId, profileId int) error {
	_, err := pg.db.ExecContext(ctx,
		`INSERT INTO project_paint_profiles (project_id, paint_profile_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		projectId, profileId,
	)
	return err
}

func (pg *PostgresPaintProfileStore) UnassignFromProject(ctx context.Context, projectId, profileId int) error {
	result, err := pg.db.ExecContext(ctx,
		`DELETE FROM project_paint_profiles WHERE project_id = $1 AND paint_profile_id = $2`,
		projectId, profileId,
	)
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
