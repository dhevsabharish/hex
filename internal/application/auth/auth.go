package auth

type AuthService interface {
	Authenticate(token string) (userID string, role string, err error)
}
