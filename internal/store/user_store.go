package store

import (
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

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	GetUserById(id int) (*User, error)
	GetUserByUsername(name string) (*User, error)
	CreateUser(user *User) (*User, error)
	UpdateUser(userId int, user *User) (*User, error)
	DeleteUser(userId int) error
}

func (pg *PostgresUserStore) GetUserById(id int) (*User, error) {
	user := &User{Password: password{}}

	query := `SELECT id, username, email, created_at, updated_at
	FROM users 
	WHERE id = $1 AND isDeleted = FALSE`

	err := pg.db.QueryRow(query, id).Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
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

	query := `SELECT id, username, email, created_at, updated_at
	FROM users 
	WHERE username = $1 AND isDeleted = FALSE`

	err := pg.db.QueryRow(query, name).Scan(&user.Id, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
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
		return nil, err
	}

	return pg.GetUserById(user.Id)
}

func (pg *PostgresUserStore) UpdateUser(userId int, user *User) (*User, error) {
	query := `UPDATE users 
	SET username = $1, email = $2, password_hash = $3, updated_at = NOW() 
	WHERE id = $4 AND isDeleted = FALSE`

	result, err := pg.db.Exec(query, user.Username, user.Email, user.Password.hash, userId)
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

	return pg.GetUserById(userId)
}

func (pg *PostgresUserStore) DeleteUser(userId int) error {
	query := `UPDATE users 
	SET isDeleted = TRUE, updated_at = NOW() 
	WHERE id = $1`

	result, err := pg.db.Exec(query, userId)
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

	return nil
}
