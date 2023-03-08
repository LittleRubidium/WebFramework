package account

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (login *loginParams) name()  {
	
}
