package sqlstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
)

const Alias = "sql"

var (
	ErrDatabaseConnect     = errors.New("unable to connect to database")
	ErrDatabaseClose       = errors.New("unable to close database")
	ErrCreateEvent         = errors.New("create event error")
	ErrUpdateEvent         = errors.New("update event error")
	ErrRemoveEvent         = errors.New("removing event error")
	ErrGetDayAheadEvents   = errors.New("getting day events error")
	ErrGetWeekAheadEvents  = errors.New("getting week events error")
	ErrGetMonthAheadEvents = errors.New("getting month events error")
	ErrGetComingEvents     = errors.New("getting coming events error")
)

type Config interface {
	GetStorageDSN() string
}

type Storage struct {
	Config Config
	db     *sqlx.DB
}

// New returns a new sql storage instance.
func New(config Config) *Storage {
	return &Storage{
		Config: config,
	}
}

// Connect is trying to connect to database server.
func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "pgx", s.Config.GetStorageDSN())
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseConnect, err)
	}

	s.db = db

	return nil
}

// Close breaks the database connection.
func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		err = fmt.Errorf("%w: %v", ErrDatabaseClose, err)
	}

	return err
}

// CreateEvent saves event into a sql storage.
func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	query := `
		INSERT INTO app_event (id, title, begin_date, end_date, description, owner_id) 
		VALUES (nextval('app_event_id_seq'), :title, :begin_date, :end_date, :description, :owner_id)
		RETURNING id
	`

	rows, err := s.db.NamedQueryContext(ctx, query, event)
	if err != nil {
		return storage.Event{}, fmt.Errorf("%w: %v", ErrCreateEvent, err)
	}

	defer rows.Close()

	for rows.Next() {
		results := make(map[string]interface{})
		err := rows.MapScan(results)
		if err != nil {
			return storage.Event{}, fmt.Errorf("%w: %v", ErrCreateEvent, err)
		}

		event.ID = results["id"].(int64)
	}

	return event, nil
}

// UpdateEvent updates event in sql storage if exists.
func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	query := `
		UPDATE app_event 
		    SET title = :title,
		        begin_date = :begin_date,
		        end_date = :end_date,
		        description = :description,
		        owner_id = :owner_id,
		        notification_sent = :notification_sent
		WHERE id = :id
	`

	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return storage.Event{}, fmt.Errorf("%w: %v", ErrUpdateEvent, err)
	}

	return event, nil
}

// RemoveEvent removes event from sql storage if exists.
func (s *Storage) RemoveEvent(ctx context.Context, event storage.Event) error {
	query := "DELETE FROM app_event WHERE id = :id"
	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrRemoveEvent, err)
	}

	return nil
}

// GetDayAheadEvents returns a day events slice.
func (s *Storage) GetDayAheadEvents(ctx context.Context) ([]storage.Event, error) {
	var events []storage.Event

	query := "SELECT * FROM app_event WHERE begin_date > NOW() AND begin_date < NOW() + interval '1 day'"
	rows, err := s.db.NamedQueryContext(ctx, query, storage.Event{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetDayAheadEvents, err)
	}

	defer rows.Close()

	for rows.Next() {
		var event storage.Event

		err := rows.StructScan(&event)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrGetDayAheadEvents, err)
		}

		events = append(events, event)
	}

	return events, nil
}

// GetWeekAheadEvents returns a week events slice.
func (s *Storage) GetWeekAheadEvents(ctx context.Context) ([]storage.Event, error) {
	var events []storage.Event

	query := "SELECT * FROM app_event WHERE begin_date > NOW() AND begin_date < NOW() + interval '1 week'"
	rows, err := s.db.NamedQueryContext(ctx, query, storage.Event{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetWeekAheadEvents, err)
	}

	defer rows.Close()

	for rows.Next() {
		var event storage.Event

		err := rows.StructScan(&event)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrGetWeekAheadEvents, err)
		}

		events = append(events, event)
	}

	return events, nil
}

// GetMonthAheadEvents returns a month events slice.
func (s *Storage) GetMonthAheadEvents(ctx context.Context) ([]storage.Event, error) {
	var events []storage.Event

	query := "SELECT * FROM app_event WHERE begin_date > NOW() AND begin_date < NOW() + interval '1 month'"
	rows, err := s.db.NamedQueryContext(ctx, query, storage.Event{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetMonthAheadEvents, err)
	}

	defer rows.Close()

	for rows.Next() {
		var event storage.Event

		err := rows.StructScan(&event)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrGetMonthAheadEvents, err)
		}

		events = append(events, event)
	}

	return events, nil
}

// GetComingEvents returns events slice, that need to be notified.
func (s *Storage) GetComingEvents(ctx context.Context) ([]storage.Event, error) {
	var events []storage.Event

	query := "SELECT * FROM app_event WHERE begin_date > NOW() AND begin_date < NOW() + interval '1 day' AND notification_sent is FALSE"
	rows, err := s.db.NamedQueryContext(ctx, query, storage.Event{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetComingEvents, err)
	}

	defer rows.Close()

	for rows.Next() {
		var event storage.Event

		err := rows.StructScan(&event)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrGetComingEvents, err)
		}

		events = append(events, event)
	}

	return events, nil
}
