package football_data

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func unmarshal(body []byte, mapTo interface{}) error {
	if err := json.Unmarshal(body, mapTo); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err) // TODO: use a custom error?
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(mapTo)
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err) // TODO: use a custom error?
	}

	return nil
}
