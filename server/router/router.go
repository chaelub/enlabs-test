package router

import (
	"enlabs-test/app/repo/transaction"
	"enlabs-test/cache/user"
	"enlabs-test/server/handlers/private/user/account"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func Router(userCache user.UserCacheI, transactionRepo transaction.TransactionRepoI) *fasthttprouter.Router {
	r := fasthttprouter.New()

	r.POST("/user/:userid/account/:accid", authMW(account.UpdateAccount)) // fixme: call UserAccount service

	return r
}

func authMW(next func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		// todo check JWT/something else
		next(ctx)
	}
}
