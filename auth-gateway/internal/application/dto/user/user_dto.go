package userdto

type User struct {
	Username string `json:"username"`
}

type Token struct {
	Token string `json:"token"`
}

type AuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	NewUsername string `json:"new_username"`
}
