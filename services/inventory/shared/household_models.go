package shared

type UserHousehold struct {
	UserID      string `json:"userId"`
	HouseholdID string `json:"householdId"`
	Name        string `json:"name"`
	Location    string `json:"location"`
	Description string `json:"description"`
}
