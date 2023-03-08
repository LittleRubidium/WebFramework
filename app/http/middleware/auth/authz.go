package auth

import (
	"fmt"
	"github.com/casbin/casbin/v3"
	"github.com/gohade/hade/app/provider/user/account"
	"github.com/gohade/hade/app/utils/jwt"
	"github.com/gohade/hade/framework/contract"
	"github.com/gohade/hade/framework/gin"
	"net/http"
	"strconv"
	"strings"
)

type BasicAuthorizer struct {
	enforce *casbin.Enforcer
}

func AuthMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	a := &BasicAuthorizer{enforce: e}
	return func(c *gin.Context) {
		if !a.CheckURLPermission(c.Request) && !a.CheckTokenValid(c) {
			a.RequirePermission(c)
		}
	}
}

func (a *BasicAuthorizer) CheckURLPermission(r *http.Request) bool {
	path := r.URL.Path
	allowed, err := a.enforce.Enforce("", path, "")
	if err != nil {
		fmt.Println(err)
	}
	if allowed {
		return true
	}
	allowed, err = a.enforce.Enforce("127.0.0.1", path, "")
	if err != nil {
		fmt.Println(err)
	}
	return allowed
}

func (a *BasicAuthorizer) CheckTokenValid(c *gin.Context) bool {
	token := c.GetHeader("Authorization")

	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return false
	}
	token = token[7:]

	userId := jwt.GetUserIdFromToken(token)
	if _, ok := c.Get(strconv.Itoa(userId)); ok {
		return true
	}
	ormService := c.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		return false
	}
	userDB := &account.User{}
	if err := db.Where("id=?", userId).First(userDB).Error; err != nil {
		return false
	}
	userDB.Password = ""
	c.Set(strconv.Itoa(userDB.Id), userDB)
	return true
}

// RequirePermission returns the 403 Forbidden to the client
func (a *BasicAuthorizer) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
}
