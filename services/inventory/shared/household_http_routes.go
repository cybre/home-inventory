package shared

import "fmt"

const (
	UserHouseholdsUserIDParam      = "userId"
	UserHouseholdsHouseholdIDParam = "householdId"
)

var (
	UserHouseholdsRoute = fmt.Sprintf("/user/:%s/households", UserHouseholdsUserIDParam)
	UserHouseholdRoute  = fmt.Sprintf("/user/:%s/households/:%s", UserHouseholdsUserIDParam, UserHouseholdsHouseholdIDParam)
)
