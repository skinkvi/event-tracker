package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
}

func (s *Storage) CreateUser(ctx context.Context, user User) (uuid.UUID, error) {
	const fn = "storage.postgres.CreateUser"

	query := `INSERT INTO users (id, email, password) VALUES ($1, $2, $3)`

	var id uuid.UUID

	err := s.db.QueryRow(ctx, query, user.ID, user.Email, user.Password).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("nothing was returned from the postgres: %w", err)
		}

		return uuid.Nil, fmt.Errorf("%s, %w", fn, err)
	}

	return id, nil
}

func (s *Storage) UpdateUser(ctx context.Context, user User) error {
	const fn = "storage.postgres.UpdateUser"

	var exists bool
	existsQ := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := s.db.QueryRow(ctx, existsQ, user.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if !exists {
		return fmt.Errorf("%s: users with id %d not found", fn, user.ID)
	}

	query := `UPDATE users SET email = $1, password = $2 WHERE id = $3`
	_, err = s.db.Exec(ctx, query, user.Email, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	const fn = "storage.postgres.GetUser"

	query := `SELECT email, password FROM users WHERE id = $1`

	var ur User

	err := s.db.QueryRow(ctx, query, id).Scan(&ur.Email, &ur.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, fmt.Errorf("%s: user with id %d not found: %w", fn, id, err)
		}

		return User{}, fmt.Errorf("%s: %w", fn, err)
	}

	return ur, nil
}

func (s *Storage) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const fn = "storage.postgres.DeleteUser"

	var exists bool
	existsQ := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := s.db.QueryRow(ctx, existsQ, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if !exists {
		return fmt.Errorf("%s: users with id %d not found", fn, id)
	}

	query := `DELETE FROM users WHERE id = $1`

	_, err = s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) ListUser(ctx context.Context) ([]User, error) {
	const fn = "storage.postgres.ListUser"

	query := `SELECT id, email, password FROM users`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var usr User

		err := rows.Scan(&usr.ID, &usr.Email, &usr.Password)
		if err != nil {
			return nil, fmt.Errorf("%s: falied to scan user row: %w", fn, err)
		}

		users = append(users, usr)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: error iterating over rows: %w", fn, rows.Err())
	}

	return users, nil
}
