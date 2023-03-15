package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var (
	HTTPHost    = os.Getenv("TESTS_HTTP_HOST")
	PostgresDSN = os.Getenv("TESTS_POSTGRES_DSN")
)

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

func init() {
	if HTTPHost == "" {
		HTTPHost = "http://0.0.0.0:8080"
	}

	if PostgresDSN == "" {
		PostgresDSN = "host=0.0.0.0 port=5432 user=calendar password=calendar dbname=calendar sslmode=disable"
	}
}

//nolint:funlen,gocognit
func TestHTTP(t *testing.T) {
	// Waiting for a few seconds in order to create a schema.
	t.Logf("Waiting for the database preparaion...\n")
	time.Sleep(5 * time.Second)

	// DB connection.
	db, err := sqlx.ConnectContext(context.Background(), "postgres", PostgresDSN)
	if err != nil {
		panic(err)
	}

	creatingURL := HTTPHost + "/event/create"
	updatingURL := HTTPHost + "/event/update"
	removingURL := HTTPHost + "/event/remove"
	dayEventsURL := HTTPHost + "/event/day"
	weekEventsURL := HTTPHost + "/event/week"
	monthEventsURL := HTTPHost + "/event/month"

	t.Run("test event crud", func(t *testing.T) {
		// CREATING REQUEST.

		// Initial event.
		event := Event{
			ID:                   1,
			Title:                "Title",
			BeginDate:            time.Now().Add(time.Hour * 12),
			EndDate:              time.Now().Add(time.Hour * 12),
			Description:          "Description",
			OwnerID:              1,
			NotificationSent:     false,
			NotificationReceived: false,
		}

		// Receiving event entities.
		var events []Event

		// Marshalling event in order to send a request.
		jsonData, err := json.Marshal(event)
		if err != nil {
			panic(err)
		}

		// Creation request.
		request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, creatingURL, bytes.NewBuffer(jsonData))
		if err != nil {
			panic(err)
		}

		// Sending the creating request.
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()

		// Checking created rows in database.
		err = db.Select(&events, "SELECT * FROM app_event WHERE title=$1", event.Title)
		t.Logf("Event creating checking...\n")
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 1, "new event should be added")
		require.Equal(t, http.StatusOK, response.StatusCode, "response status code should be ok")

		// UPDATING REQUEST.

		// Updating event entity.
		event.Title = "New Title"

		// Marshalling event in order to send a request.
		jsonData, err = json.Marshal(event)
		if err != nil {
			panic(err)
		}

		// Creating the updating request.
		request, err = http.NewRequestWithContext(context.Background(), http.MethodPost, updatingURL, bytes.NewBuffer(jsonData))
		if err != nil {
			panic(err)
		}

		// Sending the updating request.
		response, err = http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()

		// Checking updated rows in database.
		events = []Event{}
		t.Logf("Event updating checking...\n")
		err = db.Select(&events, "SELECT * FROM app_event WHERE title=$1", event.Title)
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 1, "updated event not found")
		require.Equal(t, http.StatusOK, response.StatusCode, "response status code should be ok")

		// DAY EVENTS REQUEST.

		// Creating the request.
		request, err = http.NewRequestWithContext(context.Background(), http.MethodPost, dayEventsURL, nil)
		if err != nil {
			panic(err)
		}

		// Sending the request.
		response, err = http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()

		// Getting response body.
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Unmarshalling response body.
		events = []Event{}
		err = json.Unmarshal(body, &events)
		if err != nil {
			log.Fatal(err)
		}

		// Checking day events number.
		t.Logf("Daily event receiving checking...\n")
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 1, "day events number must be 1")
		require.Equal(t, http.StatusOK, response.StatusCode, "response status code should be ok")

		// WEEK EVENTS REQUEST.

		// Creating the request.
		request, err = http.NewRequestWithContext(context.Background(), http.MethodPost, weekEventsURL, nil)
		if err != nil {
			panic(err)
		}

		// Sending the request.
		response, err = http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()

		// Getting response body.
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Unmarshalling response body.
		events = []Event{}
		err = json.Unmarshal(body, &events)
		if err != nil {
			log.Fatal(err)
		}

		// Checking day events number.
		t.Logf("Weekly event receiving checking...\n")
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 1, "day events number must be 1")
		require.Equal(t, http.StatusOK, response.StatusCode, "response status code should be ok")

		// MONTH EVENTS REQUEST.

		// Creating the request.
		request, err = http.NewRequestWithContext(context.Background(), http.MethodPost, monthEventsURL, nil)
		if err != nil {
			panic(err)
		}

		// Sending the request.
		response, err = http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()

		// Getting response body.
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Unmarshalling response body.
		events = []Event{}
		err = json.Unmarshal(body, &events)
		if err != nil {
			log.Fatal(err)
		}

		// Checking day events number.
		t.Logf("Month event receiving checking...\n")
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 1, "day events number must be 1")
		require.Equal(t, http.StatusOK, response.StatusCode, "response status code should be ok")

		// RECEIVING REQUEST.

		// Waiting for event receiving
		t.Logf("Waiting for notification receiving...\n")
		time.Sleep(5 * time.Second)

		// Checking removed rows in database.
		events = []Event{}
		err = db.Select(&events, "SELECT * FROM app_event WHERE notification_received is true")
		t.Logf("Notification receiving checking...\n")
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 1, "event notification is not received")
		require.Equal(t, http.StatusOK, response.StatusCode, "response status code should be ok")

		// REMOVING REQUEST.

		// Creating the removing request.
		request, err = http.NewRequestWithContext(context.Background(), http.MethodPost, removingURL, bytes.NewBuffer(jsonData))
		if err != nil {
			panic(err)
		}

		// Sending the removing request.
		response, err = http.DefaultClient.Do(request)
		if err != nil {
			panic(err)
		}

		defer response.Body.Close()

		// Checking removed rows in database.
		events = []Event{}
		err = db.Select(&events, "SELECT * FROM app_event WHERE title=$1", event.Title)
		t.Logf("Event removing checking...\n")
		require.NoError(t, err, "should be without errors")
		require.Len(t, events, 0, "event is not removed")
		require.Equal(t, http.StatusOK, response.StatusCode, "response status code should be ok")
	})
}
