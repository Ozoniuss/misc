package main

import (
	"dbtest/model"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// connect starts the connection with the database.
func connect() (*gorm.DB, error) {
	host := "localhost"
	port := 5432
	user := "root"
	dbname := "defaultdb"
	password := "somefatguy"

	conn, err := gorm.Open(postgres.Open(
		fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", user, password, host, port, dbname),
	))
	return conn, err
}

func run() error {

	db, err := connect()

	if err != nil {
		return fmt.Errorf("could not connect to db: %w", err)
	}

	var entities = make([]model.Entity, 0)

	// In this case I've used Debug() just to print the query.
	err = db.Debug().Order("name DESC").Order("salary ASC").Order("uuid ASC").Limit(5).
		Find(&entities).Error

	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	for _, e := range entities {
		e.PrintInfo()
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Printf("[error] %s", err.Error())
		os.Exit(1)
	}
}
