package football_data

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/paulinxe/go-football-data/types"

	_ "embed"
)

//go:embed testutil/match/finished.json
var finishedMatch string

//go:embed testutil/match/scheduled.json
var scheduledMatch string

//go:embed testutil/match/logically_invalid.json
var logicallyInvalidMatch string

//go:embed testutil/match/list.json
var matchesList string

func Test_err_is_returned_when_match_mapTo_is_nil(t *testing.T) {
	client := New("api_key")
	if err := client.GetMatch(context.Background(), "1", nil); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_err_is_returned_when_match_get_client_call_fails(t *testing.T) {
	serverHits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		serverHits++
	}))
	defer server.Close()

	mapTo := types.Match{}
	client := New("api_key", WithBaseURL(server.URL))

	if err := client.GetMatch(context.Background(), "1", &mapTo); err == nil {
		t.Fatalf("expected error, got nil")
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}
}

func Test_err_is_returned_when_match_unmarshal_fails(t *testing.T) {
	serverHits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
		serverHits++
	}))
	defer server.Close()

	mapTo := types.Match{}
	client := New("api_key", WithBaseURL(server.URL))

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
		_, _ = w.Write([]byte(logicallyInvalidMatch))
		serverHits++
	}))
	defer server.Close()

	mapTo := types.Match{}
	client := New("api_key", WithBaseURL(server.URL))

	err := client.GetMatch(context.Background(), "1", &mapTo)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected error to be *ValidationError, got %T", err)
	}

	if len(validationErr.Errs) != 3 {
		t.Fatalf("expected 3 validation errors, got %d", len(validationErr.Errs))
	}

	expectedMsgs := []string{
		"winner is required for finished/awarded matches",
		"full time score is required for finished/awarded matches",
		"half time score is required for finished/awarded matches",
	}
	for _, expected := range expectedMsgs {
		var found bool
		for _, err := range validationErr.Errs {
			if err.Error() == expected {
				found = true
				break
			}
		}

		if !found {
			t.Fatalf("expected validation errors to contain %q", expected)
		}
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}
}

func Test_we_can_get_a_finished_match(t *testing.T) {
	serverHits := 0
	calledURL := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(finishedMatch))
		serverHits++
		calledURL = r.URL.String()
	}))
	defer server.Close()

	client := New("api_key", WithBaseURL(server.URL))

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
		ID:          544391,
		HomeTeam:    types.Team{ID: 77, Name: "Athletic Club", ShortName: "Athletic", TLA: "ATH", Crest: "https://crests.football-data.org/77.png"},
		AwayTeam:    types.Team{ID: 86, Name: "Real Madrid CF", ShortName: "Real Madrid", TLA: "RMA", Crest: "https://crests.football-data.org/86.png"},
		Score:       types.Score{Winner: &winner, Duration: "REGULAR", FullTime: types.ScoreTime{Home: &home, Away: &away}, HalfTime: types.ScoreTime{Home: &halfTimeHome, Away: &halfTimeAway}},
		Competition: types.Competition{ID: 2014, Name: "Primera Division", Code: "PD", Type: "LEAGUE", Emblem: "https://crests.football-data.org/laliga.png"},
		UTCDate:     "2025-12-03T18:00:00Z",
		Status:      "FINISHED",
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
		_, _ = w.Write([]byte(scheduledMatch))
		serverHits++
		calledURL = r.URL.String()
	}))
	defer server.Close()

	client := New("api_key", WithBaseURL(server.URL))

	mapTo := types.Match{}
	err := client.GetMatch(context.Background(), "1", &mapTo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedMatch := types.Match{
		ID:          544590,
		HomeTeam:    types.Team{ID: 298, Name: "Girona FC", ShortName: "Girona", TLA: "GIR", Crest: "https://crests.football-data.org/298.png"},
		AwayTeam:    types.Team{ID: 285, Name: "Elche CF", ShortName: "Elche", TLA: "ELC", Crest: "https://crests.football-data.org/285.png"},
		Score:       types.Score{Winner: nil, Duration: "REGULAR", FullTime: types.ScoreTime{Home: nil, Away: nil}, HalfTime: types.ScoreTime{Home: nil, Away: nil}},
		Competition: types.Competition{ID: 2014, Name: "Primera Division", Code: "PD", Type: "LEAGUE", Emblem: "https://crests.football-data.org/laliga.png"},
		UTCDate:     "2026-05-24T00:00:00Z",
		Status:      "SCHEDULED",
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
		_, _ = w.Write([]byte(finishedMatch))
		serverHits++
		calledURL = r.URL.String()
	}))
	defer server.Close()

	client := New("api_key", WithBaseURL(server.URL))

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

// We only test the happy path here because the logic is the same as for a single match.
func Test_we_can_get_a_list_of_matches(t *testing.T) {
	serverHits := 0
	calledURL := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(matchesList))
		serverHits++
		calledURL = r.URL.String()
	}))
	defer server.Close()

	client := New("api_key", WithBaseURL(server.URL))

	mapTo := types.MatchesList{}
	dateFrom, _ := time.Parse(time.DateOnly, "2026-02-12")
	dateTo, _ := time.Parse(time.DateOnly, "2026-02-13")
	filters := MatchesFilter{
		Ids:      []string{"1"},
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
		Status:   "FINISHED",
	}
	err := client.GetMatches(context.Background(), filters, &mapTo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}

	expectedURL := "/matches?dateFrom=2026-02-12T00%3A00%3A00Z&dateTo=2026-02-13T00%3A00%3A00Z&ids=1&status=FINISHED"
	if calledURL != expectedURL {
		t.Fatalf("expected called URL to be %s, got %s", expectedURL, calledURL)
	}

	if len(mapTo.Matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(mapTo.Matches))
	}

	winner := "AWAY_TEAM"
	home := uint(1)
	away := uint(3)
	halfTimeHome := uint(0)
	halfTimeAway := uint(3)
	expectedMatch := types.Match{
		ID:          544214,
		HomeTeam:    types.Team{ID: 298, Name: "Girona FC", ShortName: "Girona", TLA: "GIR", Crest: "https://crests.football-data.org/298.png"},
		AwayTeam:    types.Team{ID: 87, Name: "Rayo Vallecano de Madrid", ShortName: "Rayo Vallecano", TLA: "RAY", Crest: "https://crests.football-data.org/87.png"},
		Score:       types.Score{Winner: &winner, Duration: "REGULAR", FullTime: types.ScoreTime{Home: &home, Away: &away}, HalfTime: types.ScoreTime{Home: &halfTimeHome, Away: &halfTimeAway}},
		Competition: types.Competition{ID: 2014, Name: "Primera Division", Code: "PD", Type: "LEAGUE", Emblem: "https://crests.football-data.org/laliga.png"},
		UTCDate:     "2025-08-15T17:00:00Z",
		Status:      "FINISHED",
	}

	if !reflect.DeepEqual(mapTo.Matches[0], expectedMatch) {
		t.Fatalf("expected match to be %+v, got %+v", expectedMatch, mapTo.Matches[0])
	}
}
