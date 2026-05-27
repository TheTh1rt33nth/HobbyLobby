package test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/store"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = store.Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("Failed to run migrations on test database: %v", err)
	}

	_, err = db.Exec("TRUNCATE TABLE projects CASCADE")
	if err != nil {
		t.Fatalf("Failed to truncate projects table: %v", err)
	}

	return db
}

func strPtr(s string) *string { return &s }

func TestCreateHobbyProject(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	projectStore := store.NewPostgresHobbyProjectStore(db)

	tests := []struct {
		name        string
		project     *store.HobbyProject
		expectError bool
	}{
		{
			name: "Valid project",
			project: &store.HobbyProject{
				Name:        "Test Project",
				Description: strPtr("This is a test project decription"),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdProject, err := projectStore.CreateHobbyProject(context.Background(), tt.project)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.project.Name, createdProject.Name)
			assert.Equal(t, tt.project.Description, createdProject.Description)

			retrievedProject, err := projectStore.GetHobbyProjectById(context.Background(), createdProject.Id)
			require.NoError(t, err)
			assert.Equal(t, createdProject.Id, retrievedProject.Id)
			assert.Equal(t, createdProject.Name, retrievedProject.Name)
			assert.Equal(t, createdProject.Description, retrievedProject.Description)
		})
	}
}
