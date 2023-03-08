package bot

import "time"

type Service interface {
	Add(map[string]string) map[string]string
	GetList() []Bot
	Remove(map[string]string) map[string]string
	Update(map[string]string) map[string]string
}

type Bot struct {
	Id int `gorm:"column:id;primaryKey"`
	UserId int `gorm:"column:user_id"`
	Title string `gorm:"column:title"`
	Description string `gorm:"column:description"`
	Content string	`gorm:"column:content"`
	CreateTime time.Time `gorm:"column:create_time"`
	ModifyTime time.Time `gorm:"column:modify_time"`
}
