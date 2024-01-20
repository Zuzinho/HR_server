package user

// User - struct for storing basic type of user
type User struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}
