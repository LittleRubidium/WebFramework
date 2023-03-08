package user

import (
    "encoding/json"
)

const UserKey = "user"

// Service 用户相关的服务
type Service interface {

    // Register 注册用户,注意这里只是将用户注册, 并没有激活, 需要调用
    // 参数：user必填，username，password, email
    // 返回值： user 带上token
    Register(string, string, string) map[string]string

    // Login 登录相关，使用用户名密码登录，获取完成User信息
    Login(string, string) map[string]string
    // Logout 登出
    //Logout(ctx context.Context, user *User) error
    // VerifyLogin 登录验证
    //VerifyLogin(ctx context.Context, token string) (*User, error)

    // GetUser 获取用户信息
    //GetUser(ctx context.Context, userID int64) (*User, error)
}

type User struct {
    Id int `json:"user_id" gorm:"column:id;primaryKey"`
    Username string `json:"username" gorm:"column:username"`
    Password string `json:"password" gorm:"column:password"`
    Photo string `json:"photo" gorm:"column:photo"`
    Rating int `gorm:"column:rating"`
}

func (user *User) MarshalBinary() ([]byte, error) {
    return json.Marshal(user)
}

func (user *User) UnmarshalBinary(b []byte) error {
    return json.Unmarshal(b,user)
}