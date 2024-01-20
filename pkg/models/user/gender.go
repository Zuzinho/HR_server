package user

// Gender - user type for storing gender value
type Gender string

const (
	Male   Gender = "male"
	Female Gender = "female"
)

// Check checks correction of gender value (is const)
func (gender Gender) Check() error {
	switch gender {
	case Male, Female:
		return nil
	default:
		return newIncorrectGenderError(gender)
	}
}
