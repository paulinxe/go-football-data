package football_data

import (
	"errors"
	"fmt"
	"strings"
)

// ErrMapToNil is returned when a Get method is called with a nil mapTo argument.
// Callers can detect it with: errors.Is(err, football_data.ErrMapToNil).
var ErrMapToNil = errors.New("mapTo cannot be nil")

// UnmarshalError wraps an error that occurred during JSON unmarshaling or struct validation.
// Callers can use errors.As to get the underlying error:
//
//	var unmarshalErr *UnmarshalError
//	if errors.As(err, &unmarshalErr) {
//	    log.Println(unmarshalErr.Err)
//	}
type UnmarshalError struct {
	Err error
}

func (e *UnmarshalError) Error() string {
	return fmt.Sprintf("failed to unmarshal: %v", e.Err)
}

func (e *UnmarshalError) Unwrap() error {
	return e.Err
}

// ValidationError holds one or more validation errors (e.g. logical match validation).
// Callers can use errors.As to inspect the list:
//
//	var validationErr *ValidationError
//	if errors.As(err, &validationErr) {
//	    for _, e := range validationErr.Errs {
//	        log.Println(e)
//	    }
//	}
type ValidationError struct {
	Errs []error
}

func (e *ValidationError) Error() string {
	if len(e.Errs) == 0 {
		return "validation failed"
	}

	msgs := make([]string, len(e.Errs))
	for i, err := range e.Errs {
		msgs[i] = err.Error()
	}

	return "failed to validate: [" + strings.Join(msgs, " ") + "]"
}

func (e *ValidationError) Unwrap() error {
	if len(e.Errs) == 0 {
		return nil
	}

	return errors.Join(e.Errs...)
}
