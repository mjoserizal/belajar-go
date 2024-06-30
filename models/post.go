package models

type Post struct {
	Id          int64  `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"type:varchar(300)" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	PublishDate string `gorm:"type:date" json:"publish_date"`
	Link        string `gorm:"type:varchar(300)" json:"link"`
}
