package db

import "gorm.io/gorm"

type GormDB *gorm.DB

type Store struct {
	*gorm.DB
}
