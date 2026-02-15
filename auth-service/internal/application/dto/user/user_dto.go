package userdto

type User struct {
	ID       string
	Username string
}

type AuthReq struct {
	Username string
	Password string
}

type UpdateUserReq struct {
	NewUsername string
}
