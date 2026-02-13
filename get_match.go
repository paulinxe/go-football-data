package football_data

import (
	"context"
	"fmt"
	"encoding/json"
)

// GetMatch fetches a match using GET matches/{id} and maps the response.
// mapTo is a pointer to the struct you want to map the response to.
// You can use the Match struct found in types.go or pass a custom struct you have.
func (c *Client) GetMatch(ctx context.Context, matchID string, mapTo interface{}) error {
	if mapTo == nil {
		return fmt.Errorf("mapTo cannot be nil") // TODO: use a custom error?
	}
	
	path := fmt.Sprintf("/matches/%s", matchID)
	body, err := c.get(ctx, path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, mapTo); err != nil {
		return fmt.Errorf("failed to unmarshal match: %w", err) // TODO: use a custom error?
	}

	validator := NewValidator()
	errs := validator.ValidateStruct(mapTo)
	if len(errs) > 0 {
		return fmt.Errorf("failed to unmarshal match: %v", errs) // TODO: use a custom error?
	}

	return nil
}