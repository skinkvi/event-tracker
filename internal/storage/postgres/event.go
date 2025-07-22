package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Event struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Type   string
	// Это типо само событие к примеру "праздник!"
	Payload string
}

func (s *Storage) CreateEvent(ctx context.Context, event Event) (uuid.UUID, error) {
	const fn = "storage.postgres.CreateEvent"

	query := `INSERT INTO events (user_id, type, payload) VALUES ($1, $2, $3)`

	var id uuid.UUID

	err := s.db.QueryRow(ctx, query, event.UserID, event.Type, event.Payload).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, fmt.Errorf("%s: nothing was returned from the postgres: %w", fn, err)
		}

		return uuid.Nil, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event Event) error {
	const fn = "storage.postgres.UpdateEvent"

	var exists bool
	existsQ := `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`
	err := s.db.QueryRow(ctx, existsQ, event.ID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if !exists {
		return fmt.Errorf("%s: event with id %d not found", fn, event.ID)
	}

	query := `UPDATE events SET user_id = $1, type = $2, payload = $3 WHERE id = $4`
	_, err = s.db.Exec(ctx, query, event.UserID, event.Type, event.Payload, event.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, id int) (Event, error) {
	const fn = "storage.postres.GetEvent"

	query := `SELECT user_id, type, payload FROM events WHERE id = $1`

	var ev Event

	err := s.db.QueryRow(ctx, query, id).Scan(ev.UserID, ev.Type, ev.Payload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Event{}, fmt.Errorf("%s: nothing was returned from the postgres: %w", fn, err)
		}

		return Event{}, fmt.Errorf("%s: %w", fn, err)
	}

	return ev, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int) error {
	const fn = "storage.postgres.DeleteEvent"

	var exists bool
	existsQ := `SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)`
	err := s.db.QueryRow(ctx, existsQ, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if !exists {
		return fmt.Errorf("%s: event with id %d not found", fn, id)
	}

	query := `DELETE FROM events WHERE id = $1`

	_, err = s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) ListEvent(ctx context.Context) ([]Event, error) {
	const fn = "storage.postgres.ListEvent"

	query := `SELECT id, user_id, type, payload FROM events`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var evnt Event

		err := rows.Scan(&evnt.ID, &evnt.UserID, &evnt.Type, &evnt.Payload)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		events = append(events, evnt)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", fn, rows.Err())
	}

	return events, nil
}
