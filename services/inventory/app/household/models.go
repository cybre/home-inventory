package household

import "github.com/gocql/gocql"

type UserHouseholdRoomModel struct {
	HouseholdID gocql.UUID `cql:"household_id"`
	RoomID      gocql.UUID `cql:"room_id"`
	Name        string     `cql:"name"`
	ItemCount   int        `cql:"item_count"`
}

type UserHouseholdModel struct {
	UserID      string
	HouseholdID gocql.UUID
	Name        string
	Location    string
	Description string
	ItemCount   int
	Rooms       []UserHouseholdRoomModel
	Timestamp   int64
}
