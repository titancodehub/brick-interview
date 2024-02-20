package postgres

import (
	"fmt"
	"github.com/titancodehub/brick-interview/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// All config are hardcoded to make it simpler
func CreateConnection() (*gorm.DB, error) {
	cfg := config.GetDBConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	return db, nil
}
