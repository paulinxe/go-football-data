package football_data

import (
	"context"
)

// GetCompetitions fetches the competitions using GET competitions and maps the response.
// mapTo is a pointer to the struct you want to map the response to.
// You can use the CompetitionsList struct found in types.go or pass a custom struct you have.
func (c *Client) GetCompetitions(ctx context.Context, mapTo interface{}) error {
	if mapTo == nil {
		return ErrMapToNil
	}

	body, err := c.get(ctx, "/competitions", nil)
	if err != nil {
		return err
	}

	if err := unmarshal(body, mapTo); err != nil {
		return err
	}

	return nil
}
