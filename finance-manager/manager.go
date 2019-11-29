package finance_manager

import (
	cache "enlabs-test/cache/user"
	"enlabs-test/repo/transaction"
	"enlabs-test/repo/user"
	"time"
)

type FinanceManagerI interface {
	ProcessTransaction(*user.User, *transaction.Transaction) error
	ScheduleCancellation(duration time.Duration) error
}

type FinanceManager struct {
	userCache       cache.UserCacheI
	transactionRepo transaction.TransactionRepoI
}

func (fm *FinanceManager) ProcessTransaction(u *user.User, t *transaction.Transaction) error {
	/*
		check transaction doesn't exist
		get user from cache
		(check user has enough money to subtract)
		in one DB transaction
		- insert transaction with transaction repo
		- update user balance with transaction sum
		or (if not enough money):
		- insert transaction as failed

	*/

	return nil
}

func (fm *FinanceManager) ScheduleCancellation(d time.Duration) error {
	return nil
}

func NewFinanceManager(uc cache.UserCacheI, tr transaction.TransactionRepoI) FinanceManagerI {
	return &FinanceManager{
		userCache:       uc,
		transactionRepo: tr,
	}
}
