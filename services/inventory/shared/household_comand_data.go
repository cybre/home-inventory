package shared

type CreateHouseholdCommandData struct {
	HouseholdID string `json:"householdId" validate:"required,uuid4"`
	UserID      string `param:"userId" validate:"required"`
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Location    string `json:"location" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=200"`
}

type UpdateHouseholdCommandData struct {
	HouseholdID string `param:"householdId" validate:"required,uuid4"`
	UserID      string `param:"userId" validate:"required"`
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Location    string `json:"location" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"max=200"`
}

type DeleteHouseholdCommandData struct {
	HouseholdID string `param:"householdId" validate:"required,uuid4"`
	UserID      string `param:"userId" validate:"required"`
}

type AddRoomCommandData struct {
	HouseholdID string `param:"householdId" validate:"required,uuid4"`
	UserID      string `param:"userId" validate:"required"`
	RoomID      string `json:"roomId" validate:"required,uuid4"`
	Name        string `json:"name" validate:"required,min=3,max=50"`
}

type UpdateRoomCommandData struct {
	HouseholdID string `param:"householdId" validate:"required,uuid4"`
	UserID      string `param:"userId" validate:"required"`
	RoomID      string `param:"roomId" validate:"required,uuid4"`
	Name        string `json:"name" validate:"required,min=3,max=50"`
}

type DeleteRoomCommandData struct {
	HouseholdID string `param:"householdId" validate:"required,uuid4"`
	UserID      string `param:"userId" validate:"required"`
	RoomID      string `param:"roomId" validate:"required,uuid4"`
}
