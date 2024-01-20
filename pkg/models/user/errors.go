package user

import "fmt"

// IncorrectGenderError - error with incorrect gender value
type IncorrectGenderError struct {
	gender Gender
}

func newIncorrectGenderError(gender Gender) IncorrectGenderError {
	return IncorrectGenderError{
		gender: gender,
	}
}

func (err IncorrectGenderError) Error() string {
	return fmt.Sprintf("incorrect gender '%s'", err.gender)
}

// UnknownFieldNameError - error with unknown field name
type UnknownFieldNameError struct {
	queryName FieldName
}

func newUnknownFieldNameError(queryName FieldName) UnknownFieldNameError {
	return UnknownFieldNameError{
		queryName: queryName,
	}
}

func (err UnknownFieldNameError) Error() string {
	return fmt.Sprintf("unknown field name '%s'", err.queryName)
}

var (
	IncorrectGenderErr  = IncorrectGenderError{}
	UnknownFieldNameErr = UnknownFieldNameError{}
)
