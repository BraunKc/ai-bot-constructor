package userusecase

type TokenManager interface {
	GenerateAccessToken(userID string) (string, error)
	VerifyToken(accessToken string) (string, error)
}
