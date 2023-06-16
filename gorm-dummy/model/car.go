package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Entity models the database entity.
type Car struct {
	Id        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// printInfo prints an entity in a simpler format.
func (c *Car) PrintInfo() {
	fmt.Printf("Id: %d\tName: %s\t CreatedAt: %v\t UpdatedAt: %v\t DeletedAt: %v\n", c.Id, c.Name, c.CreatedAt, c.UpdatedAt, c.DeletedAt)
}
