package models

import (
	"gorm.io/gorm"
)

type IvrDetail struct {
	ID          int
	Name        string
	Description string
}

func (IvrDetail) TableName() string {
	return "ivr_details"
}

func GetAll(db *gorm.DB, details *[]IvrDetail) (err error) {
	if result := db.Find(&details); result.Error != nil {
		return err
	}
	return nil
}
