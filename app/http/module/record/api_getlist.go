package record

import (
	"github.com/gohade/hade/app/provider/record"
	"github.com/gohade/hade/framework/gin"
	"strconv"
)

func (r *RecordApi) GetList(c *gin.Context) {
	recordService := c.MustMake(record.RecordKey).(record.Service)
	page, _ := strconv.Atoi(c.Query("page"))
	resp := recordService.GetList(page)
	c.ISetOkStatus().IJson(resp)
}
