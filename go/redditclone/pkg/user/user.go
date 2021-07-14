package user

import "errors"

type User struct {
	ID       uint32
	Login    string
	Password string
}

type UserRepo struct {
	data   map[string]*User
	nextID uint32
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		data: map[string]*User{},
	}
}

var (
	ErrNoUser     = errors.New("user not found")
	ErrBadPass    = errors.New("invalid password")
	ErrUserExists = errors.New("username already exists")
)

func (repo *UserRepo) Authorize(login, pass string) (*User, error) {
	u, ok := repo.data[login]
	if !ok {
		return nil, ErrNoUser
	}

	if u.Password != pass {
		return nil, ErrBadPass
	}

	return u, nil
}

func (repo *UserRepo) Signup(login, pass string) (*User, error) {
	_, ok := repo.data[login]
	if ok {
		return nil, ErrUserExists
	}

	repo.nextID++
	repo.data[login] = &User{
		ID:       repo.nextID,
		Login:    login,
		Password: pass,
	}

	return repo.data[login], nil
}
