package shared

import "fmt"

const (
	UserHouseholdsUserIDParam      = "userId"
	UserHouseholdsHouseholdIDParam = "householdId"
	UserHouseholdsRoomIDParam      = "roomId"
)

var (
	UserHouseholdsRoute = fmt.Sprintf("/user/:%s/households", UserHouseholdsUserIDParam)
	UserHouseholdRoute  = fmt.Sprintf("/user/:%s/households/:%s", UserHouseholdsUserIDParam, UserHouseholdsHouseholdIDParam)

	UserHouseholdRoomsRoute = fmt.Sprintf("/user/:%s/households/:%s/rooms", UserHouseholdsUserIDParam, UserHouseholdsHouseholdIDParam)
	UserHouseholdRoomRoute  = fmt.Sprintf("/user/:%s/households/:%s/rooms/:%s", UserHouseholdsUserIDParam, UserHouseholdsHouseholdIDParam, UserHouseholdsRoomIDParam)
)
