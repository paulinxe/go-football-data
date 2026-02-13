package football_data

import (
	"context"
	"testing"
	"net/http"
	"net/http/httptest"
	"reflect"
	"github.com/paulinxe/go-football-data/types"

	_ "embed"
)

//go:embed testutil/finished_match.json
var finishedMatch string

//go:embed testutil/scheduled_match.json
var scheduledMatch string

//go:embed testutil/logically_invalid_match.json
var logicallyInvalidMatch string

func Test_err_is_returned_when_mapTo_is_nil(t *testing.T) {
	client := NewClient("api_key")
	if err := client.GetMatch(context.Background(), "1", nil); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_err_is_returned_when_get_client_call_fails(t *testing.T) {
	serverHits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		serverHits++
	}))
	defer server.Close()

	mapTo := types.Match{}
	client := NewClient("api_key", WithBaseURL(server.URL))

	if err := client.GetMatch(context.Background(), "1", &mapTo); err == nil {
		t.Fatalf("expected error, got nil")
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}
}

func Test_err_is_returned_when_unmarshal_fails(t *testing.T) {
	serverHits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
		serverHits++
	}))
	defer server.Close()

	mapTo := types.Match{}
	client := NewClient("api_key", WithBaseURL(server.URL))

	if err := client.GetMatch(context.Background(), "1", &mapTo); err == nil {
		t.Fatalf("expected error, got nil")
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}
}

func Test_err_is_returned_when_match_is_logically_incorrect(t *testing.T) {
	serverHits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(logicallyInvalidMatch))
		serverHits++
	}))
	defer server.Close()

	mapTo := types.Match{}
	client := NewClient("api_key", WithBaseURL(server.URL))

	err := client.GetMatch(context.Background(), "1", &mapTo)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	// Most probably there is a better way of doing this
	if err.Error() != "failed to validate match: [winner is required for finished/awarded matches full time score is required for finished/awarded matches half time score is required for finished/awarded matches]" {
		t.Fatalf("expected error to be 'failed to validate match: [winner is required for finished/awarded matches full time score is required for finished/awarded matches half time score is required for finished/awarded matches]', got %s", err)
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}
}

// TODO: use table driven tests for this and the scheduled match test
func Test_we_can_get_a_finished_match(t *testing.T) {
	serverHits := 0
	calledURL := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(finishedMatch))
		serverHits++
		calledURL = r.URL.String()
	}))
	defer server.Close()

	client := NewClient("api_key", WithBaseURL(server.URL))

	mapTo := types.Match{}
	err := client.GetMatch(context.Background(), "1", &mapTo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	winner := "AWAY_TEAM"
	home := uint(0)
	away := uint(3)
	halfTimeHome := uint(0)
	halfTimeAway := uint(2)
	expectedMatch := types.Match{
		ID: 544391,
		HomeTeam: types.Team{ID: 77, Name: "Athletic Club", ShortName: "Athletic", TLA: "ATH", Crest: "https://crests.football-data.org/77.png"},
		AwayTeam: types.Team{ID: 86, Name: "Real Madrid CF", ShortName: "Real Madrid", TLA: "RMA", Crest: "https://crests.football-data.org/86.png"},
		Score: types.Score{Winner: &winner, Duration: "REGULAR", FullTime: types.ScoreTime{Home: &home, Away: &away}, HalfTime: types.ScoreTime{Home: &halfTimeHome, Away: &halfTimeAway}},
		Competition: types.Competition{ID: 2014, Name: "Primera Division", Code: "PD", Type: "LEAGUE", Emblem: "https://crests.football-data.org/laliga.png"},
		UTCDate: "2025-12-03T18:00:00Z",
		Status: "FINISHED",
	}

	if !reflect.DeepEqual(mapTo, expectedMatch) {
		t.Fatalf("expected match to be %+v, got %+v", expectedMatch, mapTo)
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}

	if calledURL != "/matches/1" {
		t.Fatalf("expected called URL to be /matches/1, got %s", calledURL)
	}
}

func Test_we_can_get_a_scheduled_match(t *testing.T) {
	serverHits := 0
	calledURL := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(scheduledMatch))
		serverHits++
		calledURL = r.URL.String()
	}))
	defer server.Close()

	client := NewClient("api_key", WithBaseURL(server.URL))

	mapTo := types.Match{}
	err := client.GetMatch(context.Background(), "1", &mapTo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedMatch := types.Match{
		ID: 544590,
		HomeTeam: types.Team{ID: 298, Name: "Girona FC", ShortName: "Girona", TLA: "GIR", Crest: "https://crests.football-data.org/298.png"},
		AwayTeam: types.Team{ID: 285, Name: "Elche CF", ShortName: "Elche", TLA: "ELC", Crest: "https://crests.football-data.org/285.png"},
		Score: types.Score{Winner: nil, Duration: "REGULAR", FullTime: types.ScoreTime{Home: nil, Away: nil}, HalfTime: types.ScoreTime{Home: nil, Away: nil}},
		Competition: types.Competition{ID: 2014, Name: "Primera Division", Code: "PD", Type: "LEAGUE", Emblem: "https://crests.football-data.org/laliga.png"},
		UTCDate: "2026-05-24T00:00:00Z",
		Status: "SCHEDULED",
	}

	if !reflect.DeepEqual(mapTo, expectedMatch) {
		t.Fatalf("expected match to be %+v, got %+v", expectedMatch, mapTo)
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}

	if calledURL != "/matches/1" {
		t.Fatalf("expected called URL to be /matches/1, got %s", calledURL)
	}
}

type CustomMatch struct {
	ID uint `json:"id" validate:"required"`
}
func Test_we_can_get_a_match_with_a_custom_struct(t *testing.T) {
	serverHits := 0
	calledURL := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(finishedMatch))
		serverHits++
		calledURL = r.URL.String()
	}))
	defer server.Close()

	client := NewClient("api_key", WithBaseURL(server.URL))

	mapTo := CustomMatch{}
	err := client.GetMatch(context.Background(), "1", &mapTo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedMatch := CustomMatch{
		ID: 544391,
	}

	if !reflect.DeepEqual(mapTo, expectedMatch) {
		t.Fatalf("expected match to be %+v, got %+v", expectedMatch, mapTo)
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}

	if calledURL != "/matches/1" {
		t.Fatalf("expected called URL to be /matches/1, got %s", calledURL)
	}
}