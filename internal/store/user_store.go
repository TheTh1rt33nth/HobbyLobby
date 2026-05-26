package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

type User struct {
	Id        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	GetUserById(id int) (*User, error)
	GetUserByUsername(name string) (*User, error)
	GetUserByToken(scope, tokenPlainText string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(userId int, user *User) (*User, error)
	DeleteUser(userId int) error
}

func (pg *PostgresUserStore) GetUserById(id int) (*User, error) {
	user := &User{Password: password{}}

	query := `SELECT id, username, email, password_hash, created_at, updated_at
	FROM users 
	WHERE id = $1 AND is_deleted = FALSE`

	err := pg.db.QueryRow(query, id).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pg *PostgresUserStore) GetUserByUsername(name string) (*User, error) {
	user := &User{Password: password{}}

	query := `SELECT id, username, email, password_hash, created_at, updated_at
	FROM users 
	WHERE username = $1 AND is_deleted = FALSE`

	err := pg.db.QueryRow(query, name).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pg *PostgresUserStore) GetUserByToken(scope, tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	query := `SELECT u.id, u.username, u.email, u.password_hash, u.created_at, u.updated_at
	FROM users u
	INNER JOIN tokens t ON u.id = t.user_id
	WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > NOW() AND u.is_deleted = FALSE`

	user := &User{Password: password{}}

	err := pg.db.QueryRow(query, tokenHash[:], scope).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (pg *PostgresUserStore) CreateUser(user *User) (*User, error) {
	query := `INSERT INTO users (username, email, password_hash) 
	VALUES ($1, $2, $3) 
	RETURNING id`

	err := pg.db.QueryRow(query, user.Username, user.Email, user.Password.hash).Scan(&user.Id)
	if err != nil {
		return nil, translatePgError(err)
	}

	// TODO: do not expose password hash here
	return pg.GetUserById(user.Id)
}

func (pg *PostgresUserStore) UpdateUser(userId int, user *User) (*User, error) {
	query := `UPDATE users 
	SET username = $1, email = $2, password_hash = $3, updated_at = NOW() 
	WHERE id = $4 AND is_deleted = FALSE`

	result, err := pg.db.Exec(query, user.Username, user.Email, user.Password.hash, userId)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrNotFound
	}

	// TODO: do not expose password hash here
	return pg.GetUserById(userId)
}

func (pg *PostgresUserStore) DeleteUser(userId int) error {
	query := `UPDATE users 
	SET is_deleted = TRUE, updated_at = NOW() 
	WHERE id = $1 AND is_deleted = FALSE`

	result, err := pg.db.Exec(query, userId)
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
