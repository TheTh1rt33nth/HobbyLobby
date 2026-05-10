package store

import (
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
	InsertToken(token *tokens.Token) error
	CreateToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteTokensForUser(userId int, scope string) error
}

func (pg *PostgresTokenStore) InsertToken(token *tokens.Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) 
	VALUES ($1, $2, $3, $4)`
	_, err := pg.db.Exec(query, token.Hash, token.UserId, token.Expiry, token.Scope)

	return err
}

func (pg *PostgresTokenStore) CreateToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(int(userId), ttl, scope)
	if err != nil {
		return nil, err
	}

	err = pg.InsertToken(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (pg *PostgresTokenStore) DeleteTokensForUser(userId int, scope string) error {
	query := `DELETE FROM tokens 
	WHERE user_id = $1 AND scope = $2`
	_, err := pg.db.Exec(query, userId, scope)

	return err
}
