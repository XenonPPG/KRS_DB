package models

type Note struct {
	ID      int    `gorm:"primary_key"`
	Name    string `gorm:"not null"`
	Content string

	UserID int `gorm:"not null"`
}
