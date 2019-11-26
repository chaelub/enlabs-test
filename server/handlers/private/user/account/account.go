package account

import (
	"encoding/json"
	finance_manager "enlabs-test/app/finance-manager"
	"enlabs-test/app/repo/transaction"
	"enlabs-test/cache/user"
	"github.com/valyala/fasthttp"
	"strconv"
)

// todo: pass finance worker instead of user cache and transaction repo
func UserAccountService(uc user.UserCacheI, fm finance_manager.FinanceManagerI) *UserAccount {
	return &UserAccount{
		financeManager: fm,
		userCache:      uc,
	}
}

type UserAccount struct {
	userCache      user.UserCacheI
	financeManager finance_manager.FinanceManagerI
}

func (s *UserAccount) UpdateAccount(ctx *fasthttp.RequestCtx) {
	/*
		get user id
		parse request
		validate request
		get user from cache
		pass user and transaction to finance worker
	*/
	uId, ok := ctx.UserValue("userid").(string)
	if !ok {
		// todo handle error
	}
	userId, err := strconv.Atoi(uId)
	if err != nil {
		// todo handle error
	}
	user, commit, err := s.userCache.UserW(int64(userId))
	if err != nil {
		// todo handle error
	}

	defer commit()

	// todo: parse body and create transaction instance
	b := ctx.PostBody()
	tr := new(transaction.Transaction)
	if err = json.Unmarshal(b, tr); err != nil {
		// todo handle error
	}

	if err = s.financeManager.ProcessTransaction(user, tr); err != nil {
		// todo handle error
	}

	// todo send response
}
