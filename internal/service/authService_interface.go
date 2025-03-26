package service

type AuthService interface {
	Authenticate(username, password string) (string, error)
}
