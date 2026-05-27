package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/TheTh1rt33nth/HobbyLobby/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	InsertToken(ctx context.Context, token *tokens.Token) error
	CreateToken(ctx context.Context, userId int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteTokensForUser(ctx context.Context, userId int, scope string) error
	DeleteAndCreateToken(ctx context.Context, userId int, ttl time.Duration, scope string) (*tokens.Token, error)
}

func (pg *PostgresTokenStore) InsertToken(ctx context.Context, token *tokens.Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) 
	VALUES ($1, $2, $3, $4)`
	_, err := pg.db.ExecContext(ctx, query, token.Hash, token.UserId, token.Expiry, token.Scope)

	return err
}

func (pg *PostgresTokenStore) CreateToken(ctx context.Context, userId int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(int(userId), ttl, scope)
	if err != nil {
		return nil, err
	}

	err = pg.InsertToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (pg *PostgresTokenStore) DeleteTokensForUser(ctx context.Context, userId int, scope string) error {
	query := `DELETE FROM tokens 
	WHERE user_id = $1 AND scope = $2`
	_, err := pg.db.ExecContext(ctx, query, userId, scope)

	return err
}

func (pg *PostgresTokenStore) DeleteAndCreateToken(ctx context.Context, userId int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}

	tx, err := pg.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`, userId, scope)
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`,
		token.Hash, token.UserId, token.Expiry, token.Scope)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return token, nil
}
