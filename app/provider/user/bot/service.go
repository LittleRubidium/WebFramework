package bot

import (
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"gorm.io/gorm"
	"time"
)

type BotService struct {
	container framework.Container
	logger contract.Log
	configer contract.Config
}

func (bot *BotService) Add(data map[string]interface{}) map[string]string {
	resp := map[string]string{}
	title,description,content := data["title"].(string),data["description"].(string),data["content"].(string)
	userId := data["userId"].(int)
	if title == "" || len(title) == 0 {
		resp["error_message"] = "标题不能为空"
		return resp
	}
	if len(title) > 100 {
		resp["error_message"] = "标题长度不能超过100"
		return resp
	}
	if description == "" || len(description) == 0 {
		description = "这个用户很懒，什么也没有留下～～"
	}
	if len(description) > 300 {
		resp["error_message"] = "Bot描述长度不能超过300"
		return resp
	}
	if content == "" || len(content) == 0 {
		resp["error_message"] = "代码不能为空"
		return resp
	}

	if len(content) > 10000 {
		resp["error_message"] = "代码长度不能超过10k"
		return resp
	}
	ormService := bot.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		resp["error_message"] = "创建失败，请稍后尝试"
		return resp
	}
	if users := db.Table("bot").Find(&Bot{UserId: userId}).RowsAffected; users > 10 {
		resp["error_message"] = "每个用户只能创建10个Bot"
		return resp
	}
	now := time.Now()
	addBot := &Bot{UserId: userId,Title: title,Description: description,Content: content,CreateTime: now,ModifyTime: now}
	if err := db.Table("bot").Create(addBot).Error; err != nil {
		resp["error_message"] = "创建失败，请稍后尝试"
		return resp
	}
	resp["error_message"] = "success"
	return resp
}

func (bot *BotService) GetList(userId int) []Bot {
	var bots []Bot
	ormService := bot.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		return bots
	}
	if err := db.Table("bot").Where("user_id=?",userId).Find(&bots).Error; err != nil {
		return nil
	}
	return bots
}

func (bot *BotService) Remove(botId int) map[string]string {
	resp := map[string]string{}
	ormService := bot.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		resp["error_message"] = "删除失败，请稍后尝试"
		return resp
	}
	if err := db.Table("bot").Delete(&Bot{Id: botId}).Error; err != nil {
		resp["error_message"] = "删除失败，请稍后尝试"
		return resp
	}
	resp["error_message"] = "success"
	return resp
}

func (bot *BotService) Update(data map[string]interface{}) map[string]string {
	resp := map[string]string{}
	bot_id,user_id := data["bot_id"].(int),data["user_id"].(int)
	title,description,content := data["title"].(string),data["description"].(string),data["content"].(string)
	if title == "" || len(title) == 0 {
		resp["error_message"] = "标题不能为空"
		return resp
	}
	if len(title) > 100 {
		resp["error_message"] = "标题长度不能超过100"
		return resp
	}
	if description == "" || len(description) == 0 {
		description = "这个用户很懒，什么也没有留下～～"
	}
	if len(description) > 300 {
		resp["error_message"] = "Bot描述长度不能超过300"
		return resp
	}
	if content == "" || len(content) == 0 {
		resp["error_message"] = "代码不能为空"
		return resp
	}

	if len(content) > 10000 {
		resp["error_message"] = "代码长度不能超过10k"
		return resp
	}

	ormService := bot.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		resp["error_message"] = "修改失败，请稍后尝试"
		return resp
	}
	botDB := &Bot{}
	if db.Table("bot").Where("id=?",bot_id).First(botDB).Error == gorm.ErrRecordNotFound {
		resp["error_message"] = "改Bot不存在"
		return resp
	}
	if botDB.UserId != user_id {
		resp["error_message"] = "没有权限修改该Bot"
		return resp
	}

	newBot := &Bot{Id: bot_id, UserId: user_id, Title: title, Description: description, Content: content, CreateTime: botDB.CreateTime, ModifyTime: time.Now()}
	if err := db.Table("bot").Where("id=?",bot_id).Updates(newBot).Error; err != nil {
		resp["error_message"] = "修改失败"
		return resp
	}
	resp["error_message"] = "success"
	return resp
}

func NewBotService(params []interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	logger := container.MustMake(contract.LogKey).(contract.Log)
	configer := container.MustMake(contract.ConfigKey).(contract.Config)
	return &BotService{container: container,logger: logger,configer: configer},nil
}
