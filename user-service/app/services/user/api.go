package user

type Service interface {
	CreateUser(email string) (User, error)
	GetUserBalance(email string) (float64, error)
}
