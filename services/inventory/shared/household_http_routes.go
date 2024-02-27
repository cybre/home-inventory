package shared

import "fmt"

const (
	UserHouseholdsUserIDParam = "userId"
)

var (
	UserHouseholdsRoute = fmt.Sprintf("/user/:%s/households", UserHouseholdsUserIDParam)
)
