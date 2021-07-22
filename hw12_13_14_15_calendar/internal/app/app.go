package app

import (
	"time"

	"github.com/spendmail/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Storage interface {
	AddEvent(event storage.Event) error
	UpdateEvent(event storage.Event) error
	RemoveEvent(eventID string) error
	DailyEvents(date time.Time) ([]storage.Event, error)
	WeeklyEvents(date time.Time) ([]storage.Event, error)
	MonthEvents(date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger,
		storage,
	}
}

//func (a *App) CreateEvent(ctx context.Context, id, title string) error {
//	// TODO
//	return nil
//	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
//}

// TODO

//func (a *App) CreateEvent(title string, date int64) error {
//	id, err := uuid.NewRandom()
//	if err != nil {
//		return fmt.Errorf("can not create unique id, %w", err)
//	}
//
//	event := storage.Event{ID: id.String(), Title: title, Date: date, OwnerID: ownerID}
//
//	return a.storage.AddEvent(event)
//}
//
//func (a *App) UpdateEvent(id, title string, date int64, description string, durationUntil int64, ownerID string, noticeBefore int64) error {
//	event := storage.Event{ID: id, Title: title, Date: date, DurationUntil: durationUntil, Description: description, OwnerID: ownerID, NoticeBefore: noticeBefore}
//
//	return a.storage.UpdateEvent(event)
//}
//
//func (a *App) RemoveEvent(id string) error {
//	return a.storage.RemoveEvent(id)
//}
//
//func (a *App) DailyEvents(date time.Time) ([]storage.Event, error) {
//	return a.storage.DailyEvents(date)
//}
//
//func (a *App) WeeklyEvents(date time.Time) ([]storage.Event, error) {
//	return a.storage.WeeklyEvents(date)
//}
//
//func (a *App) MonthEvents(date time.Time) ([]storage.Event, error) {
//	return a.storage.MonthEvents(date)
//}
