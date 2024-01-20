package forecaster

import "HR/pkg/models/user"

// IForecaster - interface for enriching user
type IForecaster interface {
	ForecastUser(u *user.User, enrichUserCh chan *user.EnrichedUser, errCh chan error)
}
