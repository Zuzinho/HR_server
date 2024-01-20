package forecaster

import "HR/pkg/models/user"

// ageResp stored response from agify
type ageResp struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

// genderResp stored response from genderize
type genderResp struct {
	Count       int         `json:"count"`
	Name        string      `json:"name"`
	Gender      user.Gender `json:"gender"`
	Probability float64     `json:"probability"`
}

// nationResp stored response from nationalize
type nationResp struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}
