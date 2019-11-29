package router

import (
	cache "enlabs-test/cache/user"
	finance_manager "enlabs-test/finance-manager"
	"enlabs-test/server/handlers/private/user/account"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func Router(userCache cache.UserCacheI, financeManager finance_manager.FinanceManagerI) *fasthttprouter.Router {
	r := fasthttprouter.New()

	ua := account.UserAccountService(userCache, financeManager)
	r.POST("/user/:userid/account", authMW(ua.UpdateAccount))
	r.GET("/health", func(ctx *fasthttp.RequestCtx) {
		ctx.Write([]byte("ok: I'm alive!"))
	})

	return r
}

func authMW(next func(*fasthttp.RequestCtx)) func(*fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		// todo check JWT/something else
		next(ctx)
	}
}
