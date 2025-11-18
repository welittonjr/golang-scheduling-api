package user

type UserInput struct {
	ID int `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Password string `json:"passowrd"`
	Role  string `json:"role"`
}

type UserOutput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}