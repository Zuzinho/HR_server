package forecaster

import (
	"HR/pkg/env"
	"HR/pkg/models/user"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const (
	Agify       = "https://api.agify.io"
	Genderize   = "https://api.genderize.io"
	Nationalize = "https://api.nationalize.io"
)

// Forecaster - implementor of IForecaster interface
type Forecaster struct {
	config *env.LinksConfig
}

// NewForecaster - constructor
func NewForecaster(config *env.LinksConfig) *Forecaster {
	return &Forecaster{
		config: config,
	}
}

// ForecastUser gets probable age, gender and nation by LinksConfig
// And writes EnrichedUser to enrichUserCh or error to errCh
func (fc *Forecaster) ForecastUser(u *user.User, enrichUserCh chan *user.EnrichedUser, errCh chan error) {
	var age int
	var gender user.Gender
	var nation string

	var err error

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()

		age, err = fc.forecastAge(u.Name)

		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()

		gender, err = fc.forecastGender(u.Name)

		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		defer wg.Done()

		nation, err = fc.forecastNation(u.Surname)

		if err != nil {
			errCh <- err
		}
	}()

	wg.Wait()

	enrichUserCh <- user.NewEnriched(u, age, gender, nation)
}

// forecastAge forecasts age by name
func (fc *Forecaster) forecastAge(name string) (int, error) {
	u, err := url.Parse(fc.config.AgifyLink)
	if err != nil {
		return 0, err
	}

	q := u.Query()
	q.Add("name", name)
	u.RawQuery = q.Encode()

	r := new(ageResp)

	err = fc.forecast(u, r)
	if err != nil {
		return 0, err
	}

	log.Printf("name '%s', propably age %d", r.Name, r.Age)

	return r.Age, nil
}

// forecastGender forecasts gender by name
func (fc *Forecaster) forecastGender(name string) (user.Gender, error) {
	u, err := url.Parse(fc.config.GenderizeLink)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("name", name)
	u.RawQuery = q.Encode()

	r := new(genderResp)

	err = fc.forecast(u, r)
	if err != nil {
		return "", err
	}

	log.Printf("name '%s', propably gender %s", r.Name, r.Gender)

	err = r.Gender.Check()
	if err != nil {
		return "", err
	}

	return r.Gender, nil
}

// forecastNation forecasts nation by surname
func (fc *Forecaster) forecastNation(surname string) (string, error) {
	u, err := url.Parse(fc.config.NationalizeLink)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("name", surname)
	u.RawQuery = q.Encode()

	r := new(nationResp)

	err = fc.forecast(u, r)
	if err != nil {
		return "", err
	}

	log.Printf("name '%s', propably countries %v", r.Name, r.Country)

	if len(r.Country) == 0 {
		return "", newNoCountryError(r.Name)
	}

	return r.Country[0].CountryID, nil
}

// forecast fills form by response of GET method to u
func (fc *Forecaster) forecast(u *url.URL, form any) error {
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(buf, form)

	return err
}
