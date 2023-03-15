package storage

import "time"

type Event struct {
	ID                   int64     `db:"id" json:"id"`
	Title                string    `db:"title" json:"title"`
	BeginDate            time.Time `db:"begin_date" json:"begin_date"`
	EndDate              time.Time `db:"end_date" json:"end_date"`
	Description          string    `db:"description" json:"description"`
	OwnerID              int64     `db:"owner_id" json:"owner_id"`
	NotificationSent     bool      `db:"notification_sent" json:"notification_sent"`
	NotificationReceived bool      `db:"notification_received" json:"notification_received"`
}
