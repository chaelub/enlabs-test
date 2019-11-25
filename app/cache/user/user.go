package user

import (
	"enlabs-test/app/repo/user"
	cache "enlabs-test/cache/user"
)

func Cache(ur user.UserRepoI) (cache.UserCacheI, error) {
	return &UserCache{
		users:    make(map[int64]*user.User),
		userRepo: ur,
	}, nil
}

type UserCache struct {
	users    map[int64]*user.User
	userRepo user.UserRepoI
}

func (uc *UserCache) User(id int64) (*user.User, error) {
	u := uc.users[id]
	return u, nil
}
