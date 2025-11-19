package validator

import (
	"fmt"
	"reflect"
	"strings"
)

// Validator provides simple validation capabilities
type Validator struct{}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{}
}

// Validate performs basic validation on a struct
// For production, consider using github.com/go-playground/validator/v10
func (v *Validator) Validate(s interface{}) []string {
	var errors []string
	
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return errors
	}
	
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		
		// Check for required tag
		if tag := field.Tag.Get("validate"); strings.Contains(tag, "required") {
			if isEmpty(value) {
				errors = append(errors, fmt.Sprintf("%s is required", field.Name))
			}
		}
		
		// Check for email tag
		if tag := field.Tag.Get("validate"); strings.Contains(tag, "email") {
			if value.Kind() == reflect.String && value.String() != "" {
				if !isValidEmail(value.String()) {
					errors = append(errors, fmt.Sprintf("%s must be a valid email", field.Name))
				}
			}
		}
	}
	
	return errors
}

// isEmpty checks if a reflect.Value is empty
func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}
	return false
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// ValidateRequired checks if a string is not empty
func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidateEmail checks if a string is a valid email format
func ValidateEmail(email string) error {
	if !isValidEmail(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidateMinLength checks if a string meets minimum length
func ValidateMinLength(value string, min int, fieldName string) error {
	if len(value) < min {
		return fmt.Errorf("%s must be at least %d characters", fieldName, min)
	}
	return nil
}

// ValidateMaxLength checks if a string doesn't exceed maximum length
func ValidateMaxLength(value string, max int, fieldName string) error {
	if len(value) > max {
		return fmt.Errorf("%s must not exceed %d characters", fieldName, max)
	}
	return nil
}
