package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

type User struct {
	Id             uuid.UUID `db:"id"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"hashed_password"`
	CreatedAt      time.Time `db:"created_at"`
}

func (u *User) ComparePassword(password string) error {
	passwordBytes, err := base64.StdEncoding.DecodeString(u.HashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(passwordBytes, []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password")
	}

	return nil
}
func (s *UserStore) CreateUser(ctx context.Context, email string, password string) (*User, error) {
	const userDml = `INSERT INTO users (email, hashed_password) VALUES ($1, $2) RETURNING *`

	bcrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	hashedPassword := base64.StdEncoding.EncodeToString(bcrypted)

	var user User
	if err := s.db.GetContext(ctx, &user, userDml, email, hashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user already exists")
		}
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (s *UserStore) ByEmail(ctx context.Context, email string) (*User, error) {
	const query = `SELECT * FROM users WHERE email = $1`

	var user User
	if err := s.db.GetContext(ctx, &user, query, email); err != nil {
		return nil, fmt.Errorf("failed to fetch user by email: %w", err)
	}

	return &user, nil
}
func (s *UserStore) ById(ctx context.Context, id uuid.UUID) (*User, error) {
	const query = `SELECT * FROM users WHERE id = $1`

	var user User
	if err := s.db.GetContext(ctx, &user, query, id); err != nil {
		return nil, fmt.Errorf("failed to fetch user by id %s: %w", id, err)
	}

	return &user, nil
}
