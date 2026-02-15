package main

import (
	"context"
	"fmt"
	"log"

	"github.com/paulinxe/go-football-data"
	"github.com/paulinxe/go-football-data/types"
)

func main() {
	client := football_data.New("YOUR_API_KEY_HERE")
	get_competitions(client)
}

func get_competitions(client *football_data.Client) {
	competitions := types.CompetitionsList{}
	err := client.GetCompetitions(context.Background(), &competitions)
	if (err != nil) {
		log.Fatal(err)
		return
	}

	fmt.Println(competitions)
}