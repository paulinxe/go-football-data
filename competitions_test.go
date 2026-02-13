package football_data

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	_ "embed"
	"github.com/paulinxe/go-football-data/types"
)

//go:embed testutil/competition/list.json
var competitionsList string

func Test_err_is_returned_when_competitions_mapTo_is_nil(t *testing.T) {
	client := New("api_key")
	if err := client.GetCompetitions(context.Background(), nil); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_err_is_returned_when_competitions_get_client_call_fails(t *testing.T) {
	serverHits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		serverHits++
	}))
	defer server.Close()

	mapTo := types.CompetitionsList{}
	client := New("api_key", WithBaseURL(server.URL))

	if err := client.GetCompetitions(context.Background(), &mapTo); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_err_is_returned_when_competitions_unmarshal_fails(t *testing.T) {
	serverHits := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("{}"))
		serverHits++
	}))
	defer server.Close()

	mapTo := types.CompetitionsList{}
	client := New("api_key", WithBaseURL(server.URL))

	if err := client.GetCompetitions(context.Background(), &mapTo); err == nil {
		t.Fatalf("expected error, got nil")
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}
}

func Test_we_can_get_a_list_of_competitions(t *testing.T) {
	serverHits := 0
	calledURL := ""
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(competitionsList))
		serverHits++
		calledURL = r.URL.String()
	}))
	defer server.Close()

	client := New("api_key", WithBaseURL(server.URL))

	mapTo := types.CompetitionsList{}
	err := client.GetCompetitions(context.Background(), &mapTo)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if serverHits != 1 {
		t.Fatalf("expected 1 server hit, got %d", serverHits)
	}

	expectedURL := "/competitions"
	if calledURL != expectedURL {
		t.Fatalf("expected called URL to be %s, got %s", expectedURL, calledURL)
	}

	if len(mapTo.Competitions) != 1 {
		t.Fatalf("expected 1 competition, got %d", len(mapTo.Competitions))
	}

	expectedCompetition := types.Competition{
		ID:     2013,
		Name:   "Campeonato Brasileiro Série A",
		Code:   "BSA",
		Type:   "LEAGUE",
		Emblem: "https://crests.football-data.org/bsa.png",
	}

	if !reflect.DeepEqual(mapTo.Competitions[0], expectedCompetition) {
		t.Fatalf("expected competition to be %+v, got %+v", expectedCompetition, mapTo.Competitions[0])
	}
}
