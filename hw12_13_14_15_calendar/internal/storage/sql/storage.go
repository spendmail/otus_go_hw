package sqlstorage

import (
	"context"
	"log"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
)

const Alias = "sql"

type Config interface {
	GetStorageDSN() string
}

type Storage struct {
	Config Config
	db     *sqlx.DB
}

func New(config Config) *Storage {
	return &Storage{
		Config: config,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "pgx", s.Config.GetStorageDSN())
	if err != nil {
		log.Fatalln(err)
	}

	s.db = db

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	query := `
		INSERT INTO app_event (id, title, begin_date, end_date, description, owner_id) 
		VALUES (nextval('app_event_id_seq'), :title, :begin_date, :end_date, :description, :owner_id)
		RETURNING id
	`

	rows, err := s.db.NamedQueryContext(ctx, query, event)
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		results := make(map[string]interface{})
		err := rows.MapScan(results)
		if err != nil {
			log.Fatal(err)
		}

		event.ID = results["id"].(int64)
	}

	return event, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) (storage.Event, error) {
	query := `
		UPDATE app_event 
		    SET title = :title,
		        begin_date = :begin_date,
		        end_date = :end_date,
		        description = :description,
		        owner_id = :owner_id
		WHERE id = :id
	`

	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		log.Fatalln(err)
	}

	return event, nil
}

func (s *Storage) RemoveEvent(ctx context.Context, event storage.Event) error {
	query := "DELETE FROM app_event WHERE id = :id"
	_, err := s.db.NamedExecContext(ctx, query, event)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func (s *Storage) GetDayAheadEvents(ctx context.Context) ([]storage.Event, error) {
	var events []storage.Event

	query := "SELECT * FROM app_event WHERE begin_date > NOW() AND begin_date < NOW() + interval '1 day'"
	rows, err := s.db.NamedQueryContext(ctx, query, storage.Event{})
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var event storage.Event

		err := rows.StructScan(&event)
		if err != nil {
			log.Fatal(err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) GetWeekAheadEvents(ctx context.Context) ([]storage.Event, error) {
	var events []storage.Event

	query := "SELECT * FROM app_event WHERE begin_date > NOW() AND begin_date < NOW() + interval '1 week'"
	rows, err := s.db.NamedQueryContext(ctx, query, storage.Event{})
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var event storage.Event

		err := rows.StructScan(&event)
		if err != nil {
			log.Fatal(err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) GetMonthAheadEvents(ctx context.Context) ([]storage.Event, error) {
	var events []storage.Event

	query := "SELECT * FROM app_event WHERE begin_date > NOW() AND begin_date < NOW() + interval '1 month'"
	rows, err := s.db.NamedQueryContext(ctx, query, storage.Event{})
	if err != nil {
		log.Fatalln(err)
	}

	for rows.Next() {
		var event storage.Event

		err := rows.StructScan(&event)
		if err != nil {
			log.Fatal(err)
		}

		events = append(events, event)
	}

	return events, nil
}
