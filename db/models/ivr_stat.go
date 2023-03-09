package models

import (
	"gorm.io/gorm"
)

func (IVR) TableName() string {
	return "ivr_stats"
}

type IVR struct {
	ID        int
	Ivr       string
	Sip       string
	GroupDate string
	Amount    int
}

func CreateBulk(db *gorm.DB, ivrs []IVR) (err error) {

	err = db.Create(&ivrs).Error
	if err != nil {
		return err
	}
	return nil
}
