package household

import "github.com/gocql/gocql"

type UserHouseholdRoomModel struct {
	HouseholdID gocql.UUID `cql:"household_id"`
	RoomID      gocql.UUID `cql:"room_id"`
	Name        string     `cql:"name"`
	Order       uint       `cql:"sort_order"`
	Timestamp   int64      `cql:"tstamp"`
}

type UserHouseholdModel struct {
	UserID      string
	HouseholdID gocql.UUID
	Name        string
	Location    string
	Description string
	Rooms       []UserHouseholdRoomModel
	Timestamp   int64
	Order       uint
}
