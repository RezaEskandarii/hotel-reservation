package repositories

import (
	"gorm.io/gorm"
	"hotel-reservation/internal/commons"
	"hotel-reservation/internal/dto"
	"hotel-reservation/internal/models"
	"hotel-reservation/pkg/application_loger"
)

type CityRepository struct {
	DB *gorm.DB
}

//type CityRepository interface {
//	Create(city *models.City) (*models.City, error)
//	Update(city *models.City) (*models.City, error)
//	Find(city *models.City) (*models.City, error)
//	FindAll(input *dto.PaginationInput) (commons.PaginatedList, error)
//	Delete(id uint64) error
//}

func NewCityRepository(db *gorm.DB) *CityRepository {
	return &CityRepository{
		DB: db,
	}
}

func (r *CityRepository) Create(city *models.City) (*models.City, error) {

	if tx := r.DB.Create(&city); tx.Error != nil {
		application_loger.LogError(tx.Error)
		return nil, tx.Error
	}

	return city, nil
}

func (r *CityRepository) Update(city *models.City) (*models.City, error) {

	if tx := r.DB.Updates(&city); tx.Error != nil {
		application_loger.LogError(tx.Error)
		return nil, tx.Error
	}

	return city, nil
}

func (r *CityRepository) Find(id uint64) (*models.City, error) {

	model := models.City{}
	if tx := r.DB.Where("id=?", id).Find(&model); tx.Error != nil {
		application_loger.LogError(tx.Error)
		return nil, tx.Error
	}

	if model.Id == 0 {
		return nil, nil
	}

	return &model, nil
}

func (r *CityRepository) FindAll(input *dto.PaginationInput) (*commons.PaginatedList, error) {

	return finAll(&models.City{}, r.DB, input)
}
