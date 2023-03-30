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
	account := r.Group("api/user/account")
	{
		account.POST("/token/", api.Login)
		account.POST("/register/", api.Register)
		account.POST("/info/", api.Info)
	}
	return nil
}
