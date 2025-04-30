package main

import (
	"app/config"
	"app/entity"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"os"
)

var db *gorm.DB

func initDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=postgres password=postgre dbname=vds")
	env := os.Getenv("DB_DSN")
	if env != "" {
		dsn = env
	}
	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&entity.Vds{}); err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&entity.Token{}); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// Инициализация базы данных в функции main
	var err error
	db, err = initDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Println(err)
		} else {
			sqlDB.Close()
		}
	}()
	WebServ()
}

func WebServ() {
	cfg := config.Get()
	m := http.NewServeMux()
	m.Handle("/", http.HandlerFunc(handle))
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: m,
	}
	srv.ListenAndServe()
}
