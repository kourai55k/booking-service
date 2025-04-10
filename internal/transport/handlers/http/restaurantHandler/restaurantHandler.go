package restauranthandler

import "github.com/kourai55k/booking-service/internal/domain/models"

type RestaurantService interface {
	CreateRestaurant(*models.Restaurant) (uint, error)
	GetRestaurants() ([]*models.Restaurant, error)
	GetRestaurantByID(uint) (*models.Restaurant, error)
	UpdateRestraunt(*models.Restaurant) error
	DeleteRestraunt(uint) error

	CreateTable(*models.Table) (uint, error)
	GetTablesByRestaurantID(uint) ([]*models.Table, error)
	GetAvailableTablesByRestaurantID(uint) ([]*models.Table, error)
	GetTableByID(uint) (*models.Table, error)
	UpdateTable(*models.Table) error
	DeleteTable(uint) error

	IsOwnerOfRestaurant(uint, uint) (bool, error)
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type RestraurantHandler struct {
	restaurantService RestaurantService
	logger            Logger
}

func NewRestaurantHandler(restaurantService RestaurantService, logger Logger) *RestraurantHandler {
	return &RestraurantHandler{restaurantService: restaurantService, logger: logger}
}
