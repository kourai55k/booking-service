package models

type Restaurant struct {
	ID           uint
	Name         string
	Description  string
	Address      string
	OpeningHours []OpeningHours

	OwnerID uint
}

type OpeningHours struct {
	DayOfWeek string
	OpenTime  string
	CloseTime string

	RestaurantID uint
}
