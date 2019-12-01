package user

import (
	"enlabs-test/repo/user"
)

type UserCacheI interface {
	User(int64) (*user.User, error)
	Save(*user.User) error
}

func Cache(ur user.UserRepoI) (UserCacheI, error) {
	return &UserCache{
		userRepo: ur,
	}, nil
}

type UserCache struct {
	userRepo user.UserRepoI
}

func (uc *UserCache) User(id int64) (*user.User, error) {
	return uc.user(id)
}

func (uc *UserCache) user(id int64) (*user.User, error) {
	return uc.userRepo.User(id)
}

func (uc *UserCache) Save(u *user.User) error {
	return nil
}
