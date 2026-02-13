package football_data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/paulinxe/go-football-data/types"
)

// GetMatch fetches a match using GET matches/{id} and maps the response.
// mapTo is a pointer to the struct you want to map the response to.
// You can use the Match struct found in types.go or pass a custom struct you have.
// WARNING: if you pass a custom struct, you are responsible for logically validating it.
// This means that you have to ensure things like the winner is set when the match is finished, etc.
func (c *Client) GetMatch(ctx context.Context, matchID string, mapTo interface{}) error {
	if mapTo == nil {
		return fmt.Errorf("mapTo cannot be nil") // TODO: use a custom error?
	}

	path := fmt.Sprintf("/matches/%s", matchID)
	body, err := c.get(ctx, path, nil)
	if err != nil {
		return err
	}

	if err := unmarshal(body, mapTo); err != nil {
		return err
	}

	match, ok := mapTo.(*types.Match)
	if ok {
		errs := validateMatch(match)
		if len(errs) > 0 {
			return fmt.Errorf("failed to validate match: %v", errs) // TODO: use a custom error?
		}
	}

	return nil
}

type MatchesFilter struct {
	Ids      []string
	Date     *time.Time
	DateFrom *time.Time
	DateTo   *time.Time
	Status   string
}

// GetMatches fetches matches using GET /matches?{filters} and maps the response.
// It forces the caller to pass at least one filter.
// mapTo is a pointer to the struct you want to map the response to.
// You can use the Matches struct found in types.go or pass a custom struct you have.
// WARNING: if you pass a custom struct, you are responsible for logically validating it.
// This means that you have to ensure things like the winner is set when the match is finished, etc.
func (c *Client) GetMatches(ctx context.Context, filters MatchesFilter, mapTo interface{}) error {
	path := "/matches"
	queryParams := url.Values{}

	if len(filters.Ids) > 0 {
		queryParams.Add("ids", strings.Join(filters.Ids, ","))
	}

	if filters.Date != nil {
		queryParams.Add("date", filters.Date.Format(time.RFC3339))
	}

	if filters.DateFrom != nil {
		queryParams.Add("dateFrom", filters.DateFrom.Format(time.RFC3339))
	}

	if filters.DateTo != nil {
		queryParams.Add("dateTo", filters.DateTo.Format(time.RFC3339))
	}

	if filters.Status != "" {
		queryParams.Add("status", filters.Status)
	}

	body, err := c.get(ctx, path, &queryParams)
	if err != nil {
		return err
	}

	if err := unmarshal(body, mapTo); err != nil {
		return err
	}

	list, ok := mapTo.(*types.MatchesList)
	if ok {
		var errs []error
		for _, match := range list.Matches {
			errs = append(errs, validateMatch(&match)...)
		}

		if len(errs) > 0 {
			return fmt.Errorf("failed to validate matches: %w", errors.Join(errs...)) // TODO: use a custom error?
		}
	}

	return nil
}

func unmarshal(body []byte, mapTo interface{}) error {
	if err := json.Unmarshal(body, mapTo); err != nil {
		return fmt.Errorf("failed to unmarshal matches: %w", err) // TODO: use a custom error?
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(mapTo)
	if err != nil {
		return fmt.Errorf("failed to unmarshal matches: %w", err) // TODO: use a custom error?
	}

	return nil
}

// ValidateMatch makes sure the match instance is logically correct.
func validateMatch(match *types.Match) []error {
	if match.Status == "FINISHED" || match.Status == "AWARDED" {
		var errs []error
		if match.Score.Winner == nil {
			errs = append(errs, fmt.Errorf("winner is required for finished/awarded matches"))
		}

		if match.Score.FullTime.Home == nil || match.Score.FullTime.Away == nil {
			errs = append(errs, fmt.Errorf("full time score is required for finished/awarded matches"))
		}

		if match.Score.HalfTime.Home == nil || match.Score.HalfTime.Away == nil {
			errs = append(errs, fmt.Errorf("half time score is required for finished/awarded matches"))
		}

		return errs
	}

	return nil
}
