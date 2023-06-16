package model

import (
	"fmt"

	"github.com/google/uuid"
)

// Entity models the database entity.
type Entity struct {
	UUID      uuid.UUID
	Name      string
	OtherName string `gorm:"column:other_name"`
	Age       int
	Salary    int
}

// printInfo prints an entity in a simpler format.
func (e *Entity) PrintInfo() {
	fmt.Printf("UUID: %s\tNAME: %s\t OTHER_NAME: %s\t AGE: %d\t SALARY: %d\n", e.UUID, e.Name, e.OtherName, e.Age, e.Salary)
}
