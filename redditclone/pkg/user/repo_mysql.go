package user

import (
	"database/sql"
)

type UserMysqlRepository struct {
	DB *sql.DB
}

func NewMysqlRepo(db *sql.DB) *UserMysqlRepository {
	return &UserMysqlRepository{DB: db}
}

func (repo *UserMysqlRepository) Authorize(login, pass string) (*User, error) {
	u := &User{}
	err := repo.DB.
		QueryRow("SELECT id, login, password FROM users WHERE login = ?", login).
		Scan(&u.ID, &u.Login, &u.Password)
	if err != nil {
		return nil, ErrNoUser
	}
	//pass = CodingPass(u.Password)
	if u.Password != pass {
		return nil, ErrBadPass
	}
	return u, nil
}

func (repo *UserMysqlRepository) AddUserInRepo(login, pass string) (*User, error) {
	//pass = CodingPass(pass)
	result, err := repo.DB.Exec(
		"INSERT INTO users (`login`, `password`) VALUES (?, ?)",
		login,
		pass,
	)
	if err != nil {
		return nil, err
	}
	us := &User{
		Login:    login,
		Password: pass,
	}
	us.ID, _ = result.LastInsertId()
	return us, nil
}
