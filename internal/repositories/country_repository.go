package repositories

import (
	"gorm.io/gorm"
	"hotel-reservation/internal/commons"
	"hotel-reservation/internal/dto"
	"hotel-reservation/internal/models"
	"hotel-reservation/pkg/application_loger"
)

type CountryRepository struct {
	DB *gorm.DB
}

func NewCountryRepository(db *gorm.DB) *CountryRepository {
	return &CountryRepository{
		DB: db,
	}
}

func (r *CountryRepository) Create(country *models.Country) (*models.Country, error) {

	if tx := r.DB.Create(&country); tx.Error != nil {
		application_loger.LogError(tx.Error)
		return nil, tx.Error
	}

	return country, nil
}

func (r *CountryRepository) Update(country *models.Country) (*models.Country, error) {

	if tx := r.DB.Updates(&country); tx.Error != nil {
		application_loger.LogError(tx.Error)
		return nil, tx.Error
	}

	return country, nil
}

func (r *CountryRepository) Find(id uint64) (*models.Country, error) {

	model := models.Country{}
	if tx := r.DB.Where("id=?", id).Find(&model); tx.Error != nil {
		application_loger.LogError(tx.Error)
		return nil, tx.Error
	}

	if model.Id == 0 {
		return nil, nil
	}

	return &model, nil
}

func (r *CountryRepository) FindAll(input *dto.PaginationInput) (*commons.PaginatedList, error) {

	list := make([]models.Country, 0)
	result := &commons.PaginatedList{}

	result.Data = list
	return result, nil
}
