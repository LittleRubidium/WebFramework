package record

import (
	"github.com/gohade/hade/app/provider/user/account"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"strings"
	"time"
)

type RecordService struct {
	container framework.Container
	logger    contract.Log
	configer  contract.Config
}

func NewRecord(aId, aSx, aSy, bId, bSx, bSy int, aSteps, bSteps, gamemap, loser string, createTime time.Time) *Record {
	return &Record{
		AId:        aId,
		ASx:        aSx,
		ASy:        aSy,
		BId:        bId,
		BSx:        bSx,
		BSy:        bSy,
		ASteps:     aSteps,
		BSteps:     bSteps,
		Map:        gamemap,
		Loser:      loser,
		CreateTime: createTime,
	}
}

func (record *RecordService) GetList(page int) map[string]interface{} {
	ormService := record.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		return nil
	}
	resp := map[string]interface{}{}
	var records []Record
	if err := db.Table("record").Limit(10).Offset((page - 1) * 10).Find(&records).Error; err != nil {
		return nil
	}
	var items []map[string]interface{}
	for _, record := range records {
		userA, userB := &account.User{}, &account.User{}
		db.Where("id=?", record.AId).First(userA)
		db.Where("id=?", record.BId).First(userB)
		item := map[string]interface{}{
			"a_photo":    userA.Photo,
			"a_username": userA.Username,
			"b_photo":    userB.Photo,
			"b_username": userB.Username,
		}
		result := "平局"
		if strings.Compare("A", record.Loser) == 0 {
			result = "B胜"
		} else if strings.Compare("B", record.Loser) == 0 {
			result = "A胜"
		}
		item["result"] = result
		item["record"] = record
		items = append(items, item)
	}
	var total int64
	db.Table("record").Model(&Record{}).Count(&total)
	resp["records"] = items
	resp["record_count"] = total
	return resp
}

func NewRecordService(params []interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	logger := container.MustMake(contract.LogKey).(contract.Log)
	configer := container.MustMake(contract.ConfigKey).(contract.Config)
	return &RecordService{container: container, logger: logger, configer: configer}, nil
}
