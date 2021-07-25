package storage

import "errors"

var (
	ErrDbConnect = errors.New("unable to connect to database")
	ErrDbPrepare = errors.New("unable to prepare query")
	ErrDbExec    = errors.New("unable to execute query")

// ErrCantCreateEvent      = errors.New("can not create event")
// ErrCantUpdateEvent      = errors.New("can not update event")
// ErrCantRemoveEvent      = errors.New("can not remove event")
// ErrCantConnectToStorage = errors.New("can not connect to storage")
)
