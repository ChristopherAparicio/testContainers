package user

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type User struct {
	Email        string
	PasswordHash string
	LastLogin    time.Time
}

type SqlUserRepository struct {
	db *sql.DB
}

func NewSqlUserRepository(db *sql.DB) *SqlUserRepository {
	return &SqlUserRepository{
		db: db,
	}
}

func (s *SqlUserRepository) CreateUser(ctx context.Context, user *User) (*User, error) {
	result, err := s.db.ExecContext(ctx,
		`INSERT INTO userAuthentication (email, password_hash, last_login) VALUES ($1, $2, $3)`,
		user.Email,
		user.PasswordHash,
		user.LastLogin,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("failed to insert user")
	}

	return user, nil
}

func (s *SqlUserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT email, password_hash, last_login FROM userAuthentication WHERE email = $1`,
		email,
	)

	user := &User{}
	err := row.Scan(&user.Email, &user.PasswordHash, &user.LastLogin)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *SqlUserRepository) IsPasswordMatch(ctx context.Context, email, password string) (bool, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT email FROM userAuthentication WHERE email = $1 AND password_hash = $2`,
		email,
		password,
	)

	var found bool
	err := row.Scan(&found)
	if err != nil {
		return false, err
	}

	return found, nil
}

func (s *SqlUserRepository) DeleteUserByEmail(ctx context.Context, email string) error {
	result, err := s.db.ExecContext(ctx,
		`DELETE FROM userAuthentication WHERE email = $1`,
		email,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to delete user")
	}

	return nil
}
