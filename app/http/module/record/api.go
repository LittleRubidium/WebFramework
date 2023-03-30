package record

import (
	"github.com/gohade/hade/app/provider/record"
	"github.com/gohade/hade/framework/gin"
)

type RecordApi struct {
}

func Register(r *gin.Engine) error {
	recordApi := &RecordApi{}
	r.Bind(&record.RecordProvider{})
	r.GET("/api/record/", recordApi.GetList)
	return nil
}
