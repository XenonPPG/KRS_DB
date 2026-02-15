package initializers

import (
	"fmt"
	"time"

	descModels "github.com/XenonPPG/KRS_CONTRACTS/models"

	"DB/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(config models.Config) error {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		config.PostgresUser, config.PostgresPassword, config.PostgresHost, config.PostgresPort, config.PostgresDB)

	var err error
	// try 5 times to connect DB
	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Printf("Attempt %d failed: %v\n", i+1, err)
		time.Sleep(time.Second * 3)
	}
	if err != nil {
		return err
	}

	err = DB.AutoMigrate(&descModels.User{}, &descModels.Note{})
	return err
}
