package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const delay = 5 * time.Second

var (
	HTTPHost    = os.Getenv("TESTS_HTTP_HOST")
	PostgresDSN = os.Getenv("TESTS_POSTGRES_DSN")
)

type Event struct {
	ID               int64     `db:"id" json:"id"`
	Title            string    `db:"title" json:"title"`
	BeginDate        time.Time `db:"begin_date" json:"begin_date"`
	EndDate          time.Time `db:"end_date" json:"end_date"`
	Description      string    `db:"description" json:"description"`
	OwnerID          int64     `db:"owner_id" json:"owner_id"`
	NotificationSent bool      `db:"notification_sent" json:"notification_sent"`
}

func init() {
	if HTTPHost == "" {
		HTTPHost = "http://0.0.0.0:8080"
	}

	if PostgresDSN == "" {
		PostgresDSN = "host=0.0.0.0 port=5432 user=calendar password=calendar dbname=calendar sslmode=disable"
	}
}

func TestHTTP(t *testing.T) {
	log.Printf("wait %s for table creation...", delay)
	time.Sleep(delay)

	db, err := sqlx.ConnectContext(context.Background(), "postgres", PostgresDSN)
	if err != nil {
		panic(err)
	}

	httpURLCreate := HTTPHost + "/event/create"
	// httpUrlUpdate := HTTPHost + "/event/update"
	// httpUrlDelete := HTTPHost + "/event/delete"
	// httpUrlDaily := HTTPHost + "/event/daily"
	// httpUrlWeekly := HTTPHost + "/event/weekly"
	// httpUrlMonth := HTTPHost + "/event/month"

	t.Run("test event create", func(t *testing.T) {
		title := fmt.Sprintf("Test_%d", time.Now().Unix())

		event := Event{
			ID:    1,
			Title: title,
		}

		jsonData, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}

		request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, httpURLCreate, bytes.NewBuffer(jsonData))
		if err != nil {
			panic(err)
		}

		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		var response string
		createdMessageText := fmt.Sprintf("Event %q (%d) has been created successfully.", event.Title, event.ID)

		json.NewDecoder(resp.Body).Decode(&response)

		var events []Event

		err = db.Select(&events, "SELECT * FROM app_event WHERE title=$1", title)
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 1, "new event should be added")
		require.Equal(t, http.StatusOK, resp.StatusCode, "response statuscode should be ok")
		require.Equal(t, createdMessageText, response, "response message should be equal")
	})
}
