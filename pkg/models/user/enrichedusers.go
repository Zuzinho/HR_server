package user

// EnrichedUsers - user type for list of *EnrichedUser
type EnrichedUsers []*EnrichedUser

// Append appends user to users
func (users *EnrichedUsers) Append(user *EnrichedUser) {
	*users = append(*users, user)
}
