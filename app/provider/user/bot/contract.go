package bot

import "time"

const BotKey = "bot"

type Service interface {
	Add(map[string]interface{}) map[string]string
	GetList(int) []Bot
	Remove(int) map[string]string
	Update(map[string]interface{}) map[string]string
}

type Bot struct {
	Id          int       `gorm:"column:id;primaryKey" json:"id"`
	UserId      int       `gorm:"column:user_id" json:"userId"`
	Title       string    `gorm:"column:title" json:"title"`
	Description string    `gorm:"column:description" json:"description"`
	Content     string    `gorm:"column:content" json:"content"`
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`
	ModifyTime  time.Time `gorm:"column:modify_time" json:"modifyTime"`
}
