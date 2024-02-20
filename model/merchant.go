package model

import "time"

type Merchant struct {
	ID      string    `json:"id" gorm:"primary_key"`
	Balance int64     `json:"balance"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func (m Merchant) TableName() string {
	return "merchants"
}
