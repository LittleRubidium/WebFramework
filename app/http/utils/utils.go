package utils

import (
	"github.com/gohade/hade/app/provider/user/account"
	"github.com/gohade/hade/app/utils/jwt"
	"github.com/gohade/hade/framework/gin"
	"strconv"
)

func GetUser(c *gin.Context) *account.User {
	token := c.GetHeader("Authorization")
	token = token[7:]

	userId := jwt.GetUserIdFromToken(token)
	ctxUser, _ := c.Get(strconv.Itoa(userId))
	return ctxUser.(*account.User)
}
