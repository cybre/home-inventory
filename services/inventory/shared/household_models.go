package shared

type UserHouseholdRoom struct {
	HouseholdID string `json:"householdId"`
	RoomID      string `json:"roomId"`
	Name        string `json:"name"`
	ItemCount   int    `json:"itemCount"`
}

type UserHousehold struct {
	UserID      string              `json:"userId"`
	HouseholdID string              `json:"householdId"`
	Name        string              `json:"name"`
	Location    string              `json:"location"`
	Description string              `json:"description"`
	ItemCount   int                 `json:"itemCount"`
	Rooms       []UserHouseholdRoom `json:"rooms"`
	Timestamp   int64               `json:"timestamp"`
}
