package router

import (
	finance_manager "enlabs-test/app/finance-manager"
	"enlabs-test/cache/user"
	"enlabs-test/server/handlers/private/user/account"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func Router(userCache user.UserCacheI, financeManager finance_manager.FinanceManagerI) *fasthttprouter.Router {
	r := fasthttprouter.New()

	ua := account.UserAccountService(userCache, financeManager)
	r.POST("/user/:userid/account", authMW(ua.UpdateAccount))

	return r
}

func authMW(next func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		// todo check JWT/something else
		next(ctx)
	}
}
