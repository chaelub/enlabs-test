package user

import "enlabs-test/app/repo/user"

type UserCacheI interface {
	UserR(int64) (*user.User, func(), error)
	UserW(int64) (*user.User, func(), error)
	Save(*user.User) error
}
