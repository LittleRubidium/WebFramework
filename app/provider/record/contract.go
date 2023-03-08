package record

import "time"

const RecordKey = "record"

type Record struct {
	Id         int       `gorm:"column:id;primaryKey" json:"id"`
	AId        int       `gorm:"column:a_id" json:"aid"`
	ASx        int       `gorm:"column:a_sx" json:"asx"`
	ASy        int       `gorm:"column:a_sy" json:"asy"`
	BId        int       `gorm:"column:b_id" json:"bid"`
	BSx        int       `gorm:"column:b_sx" json:"bsx"`
	BSy        int       `gorm:"column:b_sy" json:"bsy"`
	ASteps     string    `gorm:"column:a_steps" json:"asteps"`
	BSteps     string    `gorm:"column:b_steps" json:"bsteps"`
	Map        string    `gorm:"column:map" json:"map"`
	Loser      string    `gorm:"column:loser" json:"loser"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
}

type Service interface {
	GetList(int) map[string]interface{}
}
