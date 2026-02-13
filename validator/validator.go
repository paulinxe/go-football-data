package validator

import (
	"fmt"
	"reflect"
	"github.com/paulinxe/go-football-data/types"
)

type Validator struct{}

func New() *Validator {
	return &Validator{}
}

// ValidateStruct makes sure the unmarshalling of the JSON payload is correct.
func (v *Validator) ValidateStruct(object interface{}) []error {
	var errs []error
	
	value := reflect.ValueOf(object)
	typeOf := reflect.TypeOf(object)

	// If it's a pointer to a struct, get the underlying element
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		typeOf = typeOf.Elem()
	}

	// Iterate through all fields of the struct
	for i := 0; i < value.NumField(); i++ {
		field := typeOf.Field(i)
		value := value.Field(i)
		
		tag := field.Tag.Get("validate")
		if tag == "required" {
			if !v.required(value.Interface()) {
				errs = append(errs, fmt.Errorf("field %s is required but was empty", field.Name))
			}
		}
	}

	return errs
}

// ValidateMatch makes sure the match instance is logically correct.
func (v *Validator) ValidateMatch(match *types.Match) []error {
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

func (v *Validator) required(val interface{}) bool {
	if val == nil {
		return false
	}
	rv := reflect.ValueOf(val)
	
	// Handle cases where an interface holds a nil pointer
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return false
	}

	switch rv.Kind() {
	case reflect.String:
		return len(rv.String()) > 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() != 0
	default:
		return !rv.IsZero()
	}
}