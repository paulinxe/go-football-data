package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/paulinxe/go-football-data"
	"github.com/paulinxe/go-football-data/types"
)

func main() {
	client := football_data.New("YOUR_API_KEY_HERE")
	getMatches(client)
	getMatch(client)
}

func getMatches(client *football_data.Client) {
	matches := types.MatchesList{}
	competitionId := 2014
	dateFrom, _ := time.Parse(time.DateOnly, "2026-02-12")
	dateTo, _ := time.Parse(time.DateOnly, "2026-02-13")
	filters := football_data.MatchesFilter{
		DateFrom: &dateFrom,
		DateTo:   &dateTo,
	}

	err := client.GetMatches(context.Background(), uint(competitionId), filters, &matches)
	if (err != nil) {
		log.Fatal(err)
		return
	}

	fmt.Println(matches)
}

func getMatch(client *football_data.Client) {
	match := types.Match{}
	matchId := "544391"
	err := client.GetMatch(context.Background(), matchId, &match)
	if (err != nil) {
		log.Fatal(err)
		return
	}

	fmt.Println(match)
}