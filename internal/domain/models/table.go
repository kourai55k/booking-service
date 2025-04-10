package models

type Table struct {
	ID          uint
	Number      uint
	Capacity    uint
	IsAvailable bool

	RestaurantID uint
}
