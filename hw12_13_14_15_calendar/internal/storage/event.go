package storage

type Event struct {
	ID          string `db:"id"`
	Title       string `db:"title"`
	BeginDate   int64  `db:"begin_date"`
	EndDate     int64  `db:"end_date"`
	Description string `db:"description"`
	OwnerID     string `db:"owner_id"`
}
