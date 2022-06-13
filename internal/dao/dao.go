package dao

import "gorm.io/gorm"

type Dao struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Dao {
	return &Dao{db: db}
}
