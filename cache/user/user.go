package user

import "enlabs-test/app/repo/user"

type UserCacheI interface {
	User(int64) (*user.User, error)
}
