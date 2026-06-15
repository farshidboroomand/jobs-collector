package domain

import "gorm.io/gorm"

// Bot represents a bot in the system.
type Bot struct {
	// Embed gorm.Model for ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model

	Title       string `gorm:"column:title;type:varchar(255);index;not null"`
	Status      int    `gorm:"column:status;type:int;not null"`
	JobPosition string `gorm:"column:job_position;type:varchar(255)"`
	CountryID   int    `gorm:"column:country_id;type:int;not null"`
	IsActive    bool   `gorm:"default:true"`
}

// TableName specifies the table name for Bot model.
func (Bot) TableName() string {
	return "bots"
}
