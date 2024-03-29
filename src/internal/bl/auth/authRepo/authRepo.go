package repository

type IAuthRepository interface {
	AddToken(token string) error
}
