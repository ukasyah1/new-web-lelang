package schedules

import "time"

type Contact struct{ Name, Phone string }
type KPKNL struct{ ID, Name string }
type Schedule struct {
	AssetID     string
	AuctionDate *time.Time
	Timezone    string
	KPKNL       KPKNL
	Address     string
	Contacts    []Contact
	AuctionLink string
}
type Query struct {
	KPKNLID            string
	StartDate, EndDate *time.Time
	Page, Limit        int
}
type Result struct {
	Total     int64
	Schedules []Schedule
}
