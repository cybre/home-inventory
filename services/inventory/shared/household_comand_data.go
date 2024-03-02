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

type AddItemCommandData struct {
	HouseholdID string `param:"householdId" validate:"required,uuid4"`
	RoomID      string `param:"roomId" validate:"required,uuid4"`
	ItemID      string `json:"itemId" validate:"required,uuid4"`
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Barcode     string `json:"barcode" validate:"required"`
	Quantity    uint   `json:"quantity" validate:"required"`
}

type UpdateItemCommandData struct {
	HouseholdID string `param:"householdId" validate:"required,uuid4"`
	RoomID      string `param:"roomId" validate:"required,uuid4"`
	ItemID      string `param:"itemId" validate:"required,uuid4"`
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Barcode     string `json:"barcode" validate:"required"`
	Quantity    uint   `json:"quantity" validate:"required"`
}
