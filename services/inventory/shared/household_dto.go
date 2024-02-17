package shared

type CreateHouseholdCommandData struct {
	HouseholdID string `json:"householdId" validate:"required,uuid4"`
	UserID      string `json:"userId" validate:"required,uuid4"`
	Name        string `json:"name" validate:"required,min=3,max=50"`
}

type AddRoomCommandData struct {
	HouseholdID string `param:"householdId" validate:"required,uuid4"`
	RoomID      string `json:"roomId" validate:"required,uuid4"`
	Name        string `json:"name" validate:"required,min=3,max=50"`
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
