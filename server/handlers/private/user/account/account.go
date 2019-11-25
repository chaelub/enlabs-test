package account

import (
	"enlabs-test/app/repo/transaction"
	"enlabs-test/cache/user"
	"github.com/valyala/fasthttp"
)

// todo: pass finance worker instead of user cache and transaction repo
func UserAccountService(uc user.UserCacheI, tr transaction.TransactionRepoI) *UserAccount {
	return &UserAccount{
		userCache:       uc,
		transactionRepo: tr,
	}
}

type UserAccount struct {
	userCache       user.UserCacheI
	transactionRepo transaction.TransactionRepoI
}

func (s *UserAccount) UpdateAccount(ctx *fasthttp.RequestCtx) {
	/*
		get user id
		parse request
		validate request
		get user from cache
		pass user and transaction to finance worker
	*/
}
