package account

import (
	"fmt"
	"github.com/gohade/hade/app/utils/jwt"
	"github.com/gohade/hade/framework"
	"github.com/gohade/hade/framework/contract"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type UserService struct {
	container framework.Container
	logger    contract.Log
	configer  contract.Config
}

func (u *UserService) Register(username, password, confirmPassword string) map[string]string {
	resp := map[string]string{}

	username = strings.TrimSpace(username)
	if username == "" || len(username) == 0 {
		resp["error_message"] = "用户名不能为空"
		return resp
	}
	if len(password) > 100 {
		resp["error_message"] = "密码长度不能超过100"
		return resp
	}
	if password != confirmPassword {
		resp["error_message"] = "两次输入的密码不一致"
		return resp
	}
	ormService := u.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		fmt.Println(err)
		resp["error_message"] = "注册失败，请稍后尝试"
		return resp
	}
	userDB := &User{}
	if err := db.Where(&User{Username: username}).First(userDB).Error; err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		resp["error_message"] = "用户名已存在，请换一个用户名"
		return resp
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
		resp["error_message"] = "注册失败，请稍后尝试"
		return resp
	}
	password = string(hash)
	user := &User{Username: username, Password: password, Rating: 1500, Photo: "https://cdn.acwing.com/media/user/profile/photo/114747_lg_242d90760d.jpg"}
	if err := db.Create(user).Error; err != nil {
		fmt.Println(err)
		resp["error_message"] = "注册失败，请稍后尝试"
		return resp
	}
	resp["error_message"] = "success"
	return resp
}

func (u *UserService) Login(username, password string) map[string]string {
	resp := map[string]string{}
	ormService := u.container.MustMake(contract.ORMKey).(contract.ORMService)
	db, err := ormService.GetDB()
	if err != nil {
		resp["error_message"] = "登录失败，请稍后尝试"
	}
	userDB := &User{}
	if err := db.Where("username=?", username).First(userDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			resp["error_message"] = "用户名不存在"
		} else {
			resp["error_message"] = "登录失败，请稍后尝试"
		}
		return resp
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(password)); err != nil {
		resp["error_message"] = "密码错误"
		return resp
	}
	token := jwt.GetToken(userDB.Id)

	resp["token"] = token
	resp["error_message"] = "success"
	return resp
}

func NewUserService(params []interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	logger := container.MustMake(contract.LogKey).(contract.Log)
	configer := container.MustMake(contract.ConfigKey).(contract.Config)
	return &UserService{container: container, logger: logger, configer: configer}, nil
}
