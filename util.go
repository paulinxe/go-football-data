package football_data

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

func unmarshal(body []byte, mapTo interface{}) error {
	if err := json.Unmarshal(body, mapTo); err != nil {
		return &UnmarshalError{Err: err}
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(mapTo)
	if err != nil {
		return &UnmarshalError{Err: err}
	}

	return nil
}
