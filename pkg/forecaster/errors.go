package forecaster

import "fmt"

// NoCountryError - error with empty probability countries set
type NoCountryError struct {
	name string
}

func newNoCountryError(name string) NoCountryError {
	return NoCountryError{
		name: name,
	}
}

func (err NoCountryError) Error() string {
	return fmt.Sprintf("no country for surname '%s'", err.name)
}

var NoCountryErr = NoCountryError{}
