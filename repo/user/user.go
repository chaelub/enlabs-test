package user

import (
	"database/sql"
	"fmt"
	"github.com/shopspring/decimal"
)

type UserRepoI interface {
	User(int64) (*User, error)
}

type UserRepo struct {
	store *sql.DB
}

// todo: don't store money on user model, add account model
type User struct {
	Id      int64           `json:"id"`
	Account decimal.Decimal `json:"account"`
}

func (ur *UserRepo) User(id int64) (*User, error) {
	row := ur.store.QueryRow(fmt.Sprintf(userById, id))
	u := new(User)
	account := 0.0
	err := row.Scan(&u.Id, &account)
	u.Account = decimal.NewFromFloat(account)

	return u, err
}

func NewUserRepo(store *sql.DB) (UserRepoI, error) {
	return &UserRepo{store: store}, nil
}
