package ranklist

import (
	"github.com/gohade/hade/app/provider/user/account"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
)

type RankListService struct {
	container framework.Container
	logger    contract.Log
	configer  contract.Config
}

func (rank *RankListService) GetRankList(page int) map[string]interface{} {
	ormService := rank.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		return nil
	}
	resp := map[string]interface{}{}
	var users []account.User
	if err := db.Order("rating desc").Limit(10).Offset((page - 1) * 10).Find(&users).Error; err != nil {
		return nil
	}
	for _, user := range users {
		user.Password = ""
	}
	var total int64
	db.Model(&account.User{}).Count(&total)
	resp["users"] = users
	resp["users_count"] = total
	return resp
}

func NewRankListService(params []interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	logger := container.MustMake(contract.LogKey).(contract.Log)
	configer := container.MustMake(contract.ConfigKey).(contract.Config)
	return &RankListService{container: container, logger: logger, configer: configer}, nil
}
