package domain

import "gorm.io/gorm"

// Continent represents a continent in the system.
type Continent struct {
	// Embed gorm.Model for ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model

	Name string `gorm:"column:name;type:varchar(255);index;not null;unique" json:"name"`
}
