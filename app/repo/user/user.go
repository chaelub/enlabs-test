package user

import "database/sql"

type UserRepoI interface {
	User(int64) (*User, error)
}

type UserRepo struct {
	store *sql.DB
}

type User struct {
	Id int64
}

func (ur *UserRepo) User(id int64) (*User, error) {
	// todo get user from store
	return nil, nil
}

func NewUserRepo(store *sql.DB) (UserRepoI, error) {
	return &UserRepo{store: store}, nil
}
