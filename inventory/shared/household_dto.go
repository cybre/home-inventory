package shared

type CreateHouseholdCommandData struct {
	HouseholdID string `json:"householdId"`
	UserID      string `json:"userId"`
	Name        string `json:"name"`
}

type AddRoomCommandData struct {
	HouseholdID string `param:"householdId"`
	RoomID      string `json:"roomId"`
	Name        string `json:"name"`
}

type AddItemCommandData struct {
	HouseholdID string `param:"householdId"`
	RoomID      string `param:"roomId"`
	ItemID      string `json:"itemId"`
	Name        string `json:"name"`
	Barcode     string `json:"barcode"`
	Quantity    uint   `json:"quantity"`
}

type UpdateItemCommandData struct {
	HouseholdID string `param:"householdId"`
	RoomID      string `param:"roomId"`
	ItemID      string `param:"itemId"`
	Name        string `json:"name"`
	Barcode     string `json:"barcode"`
	Quantity    uint   `json:"quantity"`
}
