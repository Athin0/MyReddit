package user

import (
	"crypto/md5"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrNoUser  = errors.New("no user found")
	ErrBadPass = errors.New("invalid Password")
)

type UserMemoryRepository struct {
	LastIndex int
	data      map[string]*User
	mutex     sync.Mutex
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data: make(map[string]*User),
	}
}

func (repo *UserMemoryRepository) Authorize(login, pass string) (*User, error) {
	u, ok := repo.data[login]
	pass = CodingPass(pass)
	if !ok {
		return nil, ErrNoUser
	}
	if u.Password != pass {
		return nil, ErrBadPass
	}
	return u, nil
}

func (repo *UserMemoryRepository) AddUserInRepo(login, pass string) (*User, error) {
	pass = CodingPass(pass)
	repo.mutex.Lock()
	repo.data[login] = &User{
		Login:    login,
		Password: pass,
		ID:       int64(repo.LastIndex),
	}
	repo.LastIndex++
	repo.mutex.Unlock()
	u, ok := repo.data[login]
	if !ok {
		return nil, ErrNoUser
	}
	if u.Password != pass {
		return nil, ErrBadPass
	}
	return u, nil
}

func CodingPass(data string) string {
	DataSignerSalt := ""
	data += DataSignerSalt
	dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	return dataHash
}
