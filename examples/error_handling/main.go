package main

import (
	"context"
	"fmt"
	"errors"
	"log"
	"github.com/paulinxe/go-football-data"
	"github.com/paulinxe/go-football-data/types"
)

func main() {
	invalidApiKey()
	mapToIsNil()
}

func invalidApiKey() {
	client := football_data.New("INVALID_API_KEY")
	competitions := types.CompetitionsList{}
	err := client.GetCompetitions(context.Background(), &competitions)
	if (err != nil) {
		var httpErr *football_data.HTTPError
		if errors.As(err, &httpErr) {
			switch httpErr.StatusCode {
			case 400:
				fmt.Println("Our API key is invalid. The API returns a 400")
				log.Println(err)
			default:
				fmt.Println("An unknown error occurred.")
				log.Println(err)
			}
		}

		return
	}
}

func mapToIsNil() {
	client := football_data.New("YOUR_API_KEY_HERE")
	err := client.GetCompetitions(context.Background(), nil)
	if (err != nil) {
		if errors.Is(err, football_data.ErrMapToNil) {
			fmt.Println("The mapTo argument is nil.")
			log.Println(err)
		}
	}
}