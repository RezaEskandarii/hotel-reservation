package repositories

import (
	"errors"
	"reservation-api/internal/commons"
	"reservation-api/internal/config"
	"reservation-api/internal/dto"
	"reservation-api/internal/message_keys"
	"reservation-api/internal/models"
	"reservation-api/internal/services/common_services"
	"reservation-api/pkg/database/connection_resolver"
	"sync"
)

type HotelRepository struct {
	ConnectionResolver  *connection_resolver.TenantConnectionResolver
	FileTransferService common_services.FileTransferer
}

func NewHotelRepository(r *connection_resolver.TenantConnectionResolver, fileTransferService common_services.FileTransferer) *HotelRepository {

	return &HotelRepository{
		ConnectionResolver:  r,
		FileTransferService: fileTransferService,
	}
}

func (r *HotelRepository) Create(hotel *models.Hotel, tenantID uint64) (*models.Hotel, error) {

	db := r.ConnectionResolver.GetDB(tenantID)

	if tx := db.Create(&hotel); tx.Error != nil {
		return nil, tx.Error
	}

	if hotel.Thumbnails != nil && len(hotel.Thumbnails) > 0 {

		var wg sync.WaitGroup
		errorsCh := make(chan error, 0)

		for _, file := range hotel.Thumbnails {
			if file != nil {
				wg.Add(1)
				go func() {
					result, err := r.FileTransferService.Upload(config.HotelsBucketName, "", file, &wg)
					if err != nil {
						errorsCh <- err
						return
					}
					thumbnail := models.Thumbnail{
						VersionID:  result.VersionID,
						HotelId:    hotel.Id,
						BucketName: result.BucketName,
						FileName:   result.FileName,
						FileSize:   result.FileSize,
					}

					if err := db.Create(&thumbnail).Error; err != nil {
						errorsCh <- err
					}
				}()
			}
		}
		select {
		case err := <-errorsCh:
			return nil, err
		default:

		}
		wg.Wait()
		close(errorsCh)
	}

	return hotel, nil
}

func (r *HotelRepository) Update(hotel *models.Hotel, tenantID uint64) (*models.Hotel, error) {

	db := r.ConnectionResolver.GetDB(tenantID)

	if tx := db.Updates(&hotel); tx.Error != nil {
		return nil, tx.Error
	}

	return hotel, nil
}

func (r *HotelRepository) Find(id uint64, tenantID uint64) (*models.Hotel, error) {

	model := models.Hotel{}
	db := r.ConnectionResolver.GetDB(tenantID)

	if tx := db.Where("id=?", id).Preload("Grades").Find(&model); tx.Error != nil {
		return nil, tx.Error
	}

	if model.Id == 0 {
		return nil, nil
	}

	return &model, nil
}

func (r *HotelRepository) FindAll(input *dto.PaginationFilter) (*commons.PaginatedResult, error) {
	db := r.ConnectionResolver.GetDB(input.TenantID)
	return paginatedList(&models.Hotel{}, db, input)
}

func (r HotelRepository) Delete(id uint64, tenantID uint64) error {

	db := r.ConnectionResolver.GetDB(tenantID)

	if query := db.Model(&models.Hotel{}).Where("id=?", id).Delete(&models.Hotel{}); query.Error != nil {
		return query.Error
	}

	return nil
}

func (r *HotelRepository) hasRepeatData(hotel *models.Hotel, tenantID uint64) error {

	var countByName int64 = 0
	db := r.ConnectionResolver.GetDB(tenantID)

	if tx := *db.Model(&models.Hotel{}).Where(&models.Hotel{Name: hotel.Name}).Count(&countByName); tx.Error != nil {
		return tx.Error
	}

	if countByName > 0 {

		return errors.New(message_keys.HotelRepeatPostalCode)
	}
	return nil
}
