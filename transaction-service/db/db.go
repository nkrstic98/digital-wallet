package db

import (
	"fmt"
	"golang.org/x/exp/slog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"transaction-service/config"
	"transaction-service/models"
)

var DB *gorm.DB

func OpenConnection(config config.Config) (*gorm.DB, error) {
	var err error
	connectionString := getDatabaseConnectionString(config.Database)

	DB, err = gorm.Open(postgres.Open(connectionString))
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to connect to the Database, error %s", err.Error()))
		return nil, err
	}

	slog.Info("🚀 Connected Successfully to the Database")

	return DB, nil
}

func ReinitDatabase() error {
	if DB.Migrator().HasTable(&models.User{}) {
		err := DB.Migrator().DropTable(&models.User{})
		if err != nil {
			return err
		}
	}
	if DB.Migrator().HasTable(&models.Transaction{}) {
		err := DB.Migrator().DropTable(&models.Transaction{})
		if err != nil {
			return err
		}
	}

	if err := DB.AutoMigrate(&models.User{}, &models.Transaction{}); err != nil {
		slog.Error(fmt.Sprintf("Failed to reinit database, error %s", err.Error()))
		return err
	}

	slog.Info("👍 Database reinit completed successfully!")
	return nil
}

func CloseConnection() error {
	defer func() {
		v := recover()
		if v != nil {
			panic(v)
		} else {
			slog.Info("Databased connection closed successfully")
		}
	}()

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func getDatabaseConnectionString(config config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.DB,
		config.Pwd,
	)
}
