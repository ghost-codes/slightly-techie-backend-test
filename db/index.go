package db

import (
	"fmt"

	models "ghost-codes/slightly-techie-blog/db/models"

	// "gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	*gorm.DB
}

func NewGorm(dbConnLink string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbConnLink), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("DB could not be initialized: %s", err)
	}

	migrations := Migrations{
		DB: db,
		Models: []interface{}{
			&models.User{},
			&models.Post{},
		},
	}

	err = RunMigrations(migrations)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Migrations struct {
	DB     *gorm.DB
	Models []interface{}
}

// RunMigrations runs migrations
func RunMigrations(migrations Migrations) error {
	for _, model := range migrations.Models {
		err := migrations.DB.AutoMigrate(model)
		if err != nil {
			return fmt.Errorf("Error running migrations, %s", err)
		}
	}
	return nil
}
