package models

import (
	"time"
)

// TimeRange defines a start and end time for a single day
type TimeRange struct {
	Open  time.Time
	Close time.Time
}

// OpeningHours defines the weekly operating hours for a restaurant
type OpeningHours struct {
	Monday    TimeRange
	Tuesday   TimeRange
	Wednesday TimeRange
	Thursday  TimeRange
	Friday    TimeRange
	Saturday  TimeRange
	Sunday    TimeRange
}

// Table represents a single dining table in the restaurant
type Table struct {
	ID           uint // Table ID
	Number       int  // Table number or identifier
	Capacity     int  // How many people this table can accommodate
	IsAvailable  bool // If the table is currently available for booking
	RestaurantID uint // Foreign key linking to the Restaurant this table belongs to
}

// Restaurant represents the restaurant entity
type Restaurant struct {
	ID           uint   // Restaurant ID
	Name         string // Restaurant name
	Description  string
	Address      string       // Could be a city or address
	OpeningHours OpeningHours // Weekly opening hours
	Tables       []Table      // A slice of tables available in the restaurant
}
