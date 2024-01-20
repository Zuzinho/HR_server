package user

// EnrichedUser - struct for storing enriched User (added Age, Gender and Nation)
type EnrichedUser struct {
	User
	Age    int
	Gender Gender
	Nation string
}

// NewEnriched - constructor
func NewEnriched(user *User, age int, gender Gender, nation string) *EnrichedUser {
	return &EnrichedUser{
		User:   *user,
		Age:    age,
		Gender: gender,
		Nation: nation,
	}
}
