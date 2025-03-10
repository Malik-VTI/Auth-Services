package config

import (
	"log/slog"
	"os"
	"services-auth/model"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseInit() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic(err.Error())
	// }
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPass := os.Getenv("DB_PASS")
	// dbName := os.Getenv("DB_NAME")

	jsonHandler := slog.NewJSONHandler(os.Stderr, nil)
	myslog := slog.New(jsonHandler)

	dsn := "sqlserver://sa:P@ssw0rd@10.30.100.82:1433?database=dbcollect_datadog"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	myslog.Info("Connected to database")
	db.AutoMigrate(&model.User{}, &model.Company{})
	DB = db
}
