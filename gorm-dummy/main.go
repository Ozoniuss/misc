package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Entity struct {
	UUID      uuid.UUID
	Name      string
	OtherName string `gorm:"column:other_name"`
	Age       int
	Salary    int
}

func (e *Entity) printInfo() {
	fmt.Printf("UUID: %s\tNAME: %s\t OTHER_NAME: %s\t AGE: %d\t SALARY: %d\n", e.UUID, e.Name, e.OtherName, e.Age, e.Salary)
}

func main() {

	host := "localhost"
	port := 5432
	user := "root"
	dbname := "defaultdb"
	password := "somefatguy"

	logger := log.Default()

	conn, err := gorm.Open(postgres.Open(
		// fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s", host, port, user, dbname, sslmode),
		fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", user, password, host, port, dbname),
	))

	if err != nil {
		logger.Fatalf("could not connect to db: %s", err.Error())
	}

	var entities = make([]Entity, 0)

	conn.Order("name DESC").Order("salary ASC").Order("uuid ASC").
		Find(&entities)

	for _, e := range entities {
		e.printInfo()
	}
}
