package shared

type CreateUserCommandData struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type GenerateOneTimeTokenCommandData struct {
	Email string `json:"email"`
}
