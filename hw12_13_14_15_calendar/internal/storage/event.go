package storage

import "time"

type Event struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	BeginDate   time.Time `db:"begin_date"`
	EndDate     time.Time `db:"end_date"`
	Description string    `db:"description"`
	OwnerID     int64     `db:"owner_id"`
}
