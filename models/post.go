package models

type Post struct {
	ID          uint   `gorm:"primaryKey"`
	UserID      uint   // Add UserID field
	Title       string `gorm:"type:varchar(300)" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	PublishDate string `gorm:"type:date" json:"publish_date"`
	Link        string `gorm:"type:varchar(300)" json:"link"`
}
