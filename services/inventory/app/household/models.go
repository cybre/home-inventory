package household

type UserHouseholdRoomModel struct {
	RoomID    string `cql:"room_id"`
	Name      string `cql:"name"`
	ItemCount int    `cql:"item_count"`
}

type UserHouseholdModel struct {
	UserID      string
	HouseholdID string
	Name        string
	Location    string
	Description string
	ItemCount   int
	Rooms       []UserHouseholdRoomModel
	Timestamp   int64
}
