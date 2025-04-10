package service

import (
	"fmt"

	"github.com/kourai55k/booking-service/internal/domain/models"
)

type TableRepository interface {
	CreateTable(*models.Table) (uint, error)
	GetTablesByRestaurantID(uint) ([]*models.Table, error)
	GetAvailableTablesByRestaurantID(uint) ([]*models.Table, error)
	GetTableByID(uint) (*models.Table, error)
	UpdateTable(*models.Table) error
	DeleteTable(uint) error
}

type RestaurantRepository interface {
	CreateRestaurant(*models.Restaurant) (uint, error)
	GetRestaurants() ([]*models.Restaurant, error)
	GetRestaurantByID(uint) (*models.Restaurant, error)
	UpdateRestraunt(*models.Restaurant) error
	DeleteRestraunt(uint) error
}

type RestaurantService struct {
	tableRepo      TableRepository
	restaurantRepo RestaurantRepository
}

// Tables management
func (s *RestaurantService) CreateTable(table *models.Table) (uint, error) {
	const op = "RestaurantService.CreateTable"

	id, err := s.tableRepo.CreateTable(table)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *RestaurantService) GetTablesByRestaurantID(restaurantID uint) ([]*models.Table, error) {
	const op = "RestaurantService.GetTablesByRestaurantID"

	tables, err := s.tableRepo.GetTablesByRestaurantID(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tables, nil
}

func (s *RestaurantService) GetAvailableTablesByRestaurantID(restaurantID uint) ([]*models.Table, error) {
	const op = "RestaurantService.GetAvailableTablesByRestaurantID"

	tables, err := s.tableRepo.GetAvailableTablesByRestaurantID(restaurantID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tables, nil
}

func (s *RestaurantService) GetTableByID(id uint) (*models.Table, error) {
	const op = "RestaurantService.GetTableByID"

	table, err := s.tableRepo.GetTableByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return table, nil
}

func (s *RestaurantService) UpdateTable(table *models.Table) error {
	const op = "RestaurantService.UpdateUser"

	err := s.tableRepo.UpdateTable(table)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *RestaurantService) DeleteTable(id uint) error {
	const op = "RestaurantService.UpdateUser"

	err := s.tableRepo.DeleteTable(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Restaurants management
func (s *RestaurantService) CreateRestaurant(restaurant *models.Restaurant) (uint, error) {
	const op = "RestaurantService.CreateRestaurant"

	id, err := s.restaurantRepo.CreateRestaurant(restaurant)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *RestaurantService) GetRestaurants() ([]*models.Restaurant, error) {
	const op = "RestaurantService.GetRestaurants"

	restaurants, err := s.restaurantRepo.GetRestaurants()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return restaurants, nil
}

func (s *RestaurantService) GetRestaurantByID(id uint) (*models.Restaurant, error) {
	const op = "RestaurantService.GetRestaurantByID"

	restaurant, err := s.restaurantRepo.GetRestaurantByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return restaurant, nil
}

func (s *RestaurantService) UpdateRestraunt(restaurant *models.Restaurant) error {
	const op = "RestaurantService.UpdateRestraunt"

	err := s.restaurantRepo.UpdateRestraunt(restaurant)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *RestaurantService) DeleteRestraunt(id uint) error {
	const op = "RestaurantService.DeleteRestraunt"

	err := s.restaurantRepo.DeleteRestraunt(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *RestaurantService) IsOwnerOfRestaurant(userID, restaurantID uint) (bool, error) {
	const op = "RestaurantService.IsOwnerOfRestaurant"

	// Retrieve the restaurant by its ID
	restaurant, err := s.restaurantRepo.GetRestaurantByID(restaurantID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	// Check if the restaurant's owner ID matches the user ID
	if restaurant.OwnerID == userID {
		return true, nil
	}

	return false, nil
}
