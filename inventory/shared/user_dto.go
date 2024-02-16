package shared

type CreateUserCommandData struct {
	UserID    string `json:"userId" validate:"required,uuid4"`
	FirstName string `json:"firstName" validate:"required,min=2,max=50"`
	LastName  string `json:"lastName" validate:"required,min=2,max=50"`
	Email     string `json:"email" validate:"required,email"`
}

type GenerateOneTimeTokenCommandData struct {
	Email string `json:"email" validate:"required,email"`
}
