package main

import (
	"dbtest/model"
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

// connect starts the connection with the database.
func connect() (*gorm.DB, error) {

	return gorm.Open(sqlite.Open("car.db"), nil)
}

func run() error {

	db, err := connect()

	if err != nil {
		return fmt.Errorf("could not connect to db: %w", err)
	}

	car := model.Car{
		Id:   40,
		Name: "opel",
	}

	err = db.Create(&car).Error
	if err != nil {
		return fmt.Errorf("could not create car: %s", err.Error())
	}

	err = db.Where("id = ?", 40).Delete(&model.Car{}).Error
	if err != nil {
		return fmt.Errorf("could not create car: %s", err.Error())
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("[error] %s\n", err.Error())
		os.Exit(1)
	}
}
