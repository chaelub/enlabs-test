package account

import (
	"encoding/json"
	"enlabs-test/cache/user"
	finance_manager "enlabs-test/finance-manager"
	transaction_repo "enlabs-test/repo/transaction"
	user_repo "enlabs-test/repo/user"
	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"strconv"
	"time"
)

// todo: pass finance worker instead of user cache and transaction repo
func UserAccountService(log *log.Logger, uc user.UserCacheI, fm finance_manager.FinanceManagerI) *UserAccount {
	return &UserAccount{
		log:            log,
		financeManager: fm,
		userCache:      uc,
	}
}

type UserAccount struct {
	log            *log.Logger
	userCache      user.UserCacheI
	financeManager finance_manager.FinanceManagerI
}

type UpdateAccountReq struct {
	TransactionId string          `json:"transactionId"`
	State         string          `json:"state"`
	Amount        decimal.Decimal `json:"amount"`
}

type UpdateAccountResp struct {
	User  *user_repo.User `json:"user"`
	Error string          `json:"error,omitempty"`
}

func (s *UserAccount) UpdateAccount(ctx *fasthttp.RequestCtx) {
	/*
		get user id
		parse request
		validate request
		get user from cache
		pass user and transaction to finance worker
		todo: get request type from Source-Type header -> save with transaction record;
	*/

	updAccResp := new(UpdateAccountResp)
	defer func() {
		resp, err := json.Marshal(updAccResp)
		if err != nil {
			s.log.Printf("error: can't marshal response: %v", err)
			return
		}

		ctx.Write(resp)
	}()

	uId, ok := ctx.UserValue("userid").(string)
	if !ok {
		s.log.Printf("error: can't cast userId [%+v] to string\n", ctx.UserValue("userid"))
		ctx.SetStatusCode(http.StatusBadRequest)
		updAccResp.Error = "bad request"
		return
	}
	userId, err := strconv.Atoi(uId)
	if err != nil {
		s.log.Printf("error: can't convert userId [%s] to int\n", uId)
		updAccResp.Error = err.Error()
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}
	user, err := s.userCache.User(int64(userId))
	if err != nil {
		// todo handle error
		s.log.Printf("error: can't get user [%d] from cache: %+v\n", userId, err)
		updAccResp.Error = err.Error()
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	// todo: parse body and create transaction instance
	b := ctx.PostBody()
	tr := new(UpdateAccountReq)
	if err = json.Unmarshal(b, tr); err != nil {
		s.log.Printf("error: can't parse request body: %+v\n", err)
		updAccResp.Error = err.Error()
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	trState := transaction_repo.TxStateUnknown
	if tr.State == "lose" {
		trState = transaction_repo.TxStateLose
	} else {
		trState = transaction_repo.TxStateWin
	}

	transaction := &transaction_repo.Transaction{
		ExtId:  tr.TransactionId,
		State:  trState,
		Status: transaction_repo.TxSuccess,
		Amount: tr.Amount,
		Tms:    time.Now().Unix(),
	}

	if user, err = s.financeManager.ProcessTransaction(user, transaction); err != nil {
		// todo handle error
		s.log.Printf("error: transaction[%s] processing error: %+v\n", transaction.ExtId, err)
		ctx.SetStatusCode(http.StatusInternalServerError)
		updAccResp.Error = err.Error()
		return
	}

	// todo send response
	updAccResp.User = user
	ctx.SetStatusCode(http.StatusOK)
}
