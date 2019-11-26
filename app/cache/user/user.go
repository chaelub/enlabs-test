package user

import (
	"enlabs-test/app/repo/user"
	cache "enlabs-test/cache/user"
	"sync"
)

func Cache(ur user.UserRepoI) (cache.UserCacheI, error) {
	return &UserCache{
		users:    make(map[int64]*userWrapper),
		userRepo: ur,
	}, nil
}

type userWrapper struct {
	*sync.RWMutex
	user *user.User
}

type UserCache struct {
	users    map[int64]*userWrapper
	userRepo user.UserRepoI
}

func (uc *UserCache) UserR(id int64) (*user.User, func(), error) {
	uw := uc.users[id]
	uw.RLock()
	return uw.user, uw.RUnlock, nil
}

func (uc *UserCache) UserW(id int64) (*user.User, func(), error) {
	uw := uc.users[id]
	uw.RLock()
	return uw.user, uw.Unlock, nil
}

func (uc *UserCache) user(id int64) (*userWrapper, error) {
	uw, got := uc.users[id]
	if !got {
		// todo get from DB and save to cache
	}
	return uw, nil
}

func (uc *UserCache) Save(u *user.User) error {
	uc.users[u.Id] = &userWrapper{
		RWMutex: new(sync.RWMutex),
		user:    u,
	}
	return nil
}
