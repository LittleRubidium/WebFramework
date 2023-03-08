package account

import (
	"github.com/gohade/hade/app/provider/user/account"
	"github.com/gohade/hade/framework/gin"
)

type AccountApi struct {
}

func Register(r *gin.Engine) error {
	api := &AccountApi{}
	r.Bind(&account.UserProvider{})
	r.POST("/api/user/account/token/",api.Login)
	r.POST("/api/user/account/register/",api.Register)
	r.GET("/api/user/account/info/",api.Info)
	return nil
}
