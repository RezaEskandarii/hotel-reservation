package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

// City city struct
type City struct {
	BaseModel
	Name       string    `json:"name" valid:"required"`
	Alias      string    `json:"alias" valid:"required"`
	ProvinceId uint64    `json:"province_id" valid:"required"`
	Province   *Province `json:"province,omitempty" gorm:"foreignkey:ProvinceId" valid:"-"`
}

func (c *City) Validate() (bool, error) {

	return govalidator.ValidateStruct(c)
}

func (c *City) BeforeCreate(tx *gorm.DB) error {
	_, err := c.Validate()

	if err != nil {
		tx.AddError(err)
		return err
	}

	return nil
}
