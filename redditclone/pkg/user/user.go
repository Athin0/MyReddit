package user

type User struct {
	ID       int64  `json:"id" bson:"id"`
	Login    string `json:"username" bson:"login"`
	Password string
}

//go:generate mockgen -source=user.go -destination=repo_mock.go -package=user UserRepo
type UserRepo interface {
	Authorize(login, pass string) (*User, error)
	AddUserInRepo(login, pass string) (*User, error)
}
