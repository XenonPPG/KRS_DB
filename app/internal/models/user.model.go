package models

type User struct {
	ID        int64  `gorm:"primaryKey;not null;autoIncrement"`
	Name      string `gorm:"not null"`
	Password  string `gorm:"not null"`
	DarkTheme bool   `gorm:"default:false"`

	Notes []Note `gorm:"foreignkey:UserID"`
}
