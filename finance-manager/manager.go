package finance_manager

import (
	"database/sql"
	cache "enlabs-test/cache/user"
	"enlabs-test/repo/transaction"
	"enlabs-test/repo/user"
	user_transaction "enlabs-test/repo/user-transaction"
	"errors"
	"log"
	"time"
)

type Config struct {
	CancellationTimeout int `toml:"cancellationTimeout"`
	Num2Cancel          int `toml:"num2Cancel"`
}

type FinanceManagerI interface {
	ProcessTransaction(*user.User, *transaction.Transaction) (*user.User, error)
	ScheduleCancellation() error
}

type FinanceManager struct {
	conf                *Config
	log                 *log.Logger
	userCache           cache.UserCacheI
	userTransactionRepo user_transaction.UserTransactionRepoI
}

func (fm *FinanceManager) ProcessTransaction(u *user.User, t *transaction.Transaction) (*user.User, error) {
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

	exists, err := fm.userTransactionRepo.CheckExistsByExtId(t.ExtId)
	if err != nil && err != sql.ErrNoRows {
		return u, err
	}

	if exists {
		return u, errors.New("transaction exists")
	}

	amount := t.Amount
	if t.State == transaction.TxStateLose {
		amount = amount.Neg()
	}

	updatedAcc := u.Account.Add(amount)
	if updatedAcc.IsNegative() {
		return u, errors.New("not enough money")
	}

	u.Account = updatedAcc
	if err = fm.userTransactionRepo.ProcessTransaction(u, t); err != nil {
		t.Status = transaction.TXFail
		return u, err
	}

	return u, nil
}

func (fm *FinanceManager) ScheduleCancellation() error {
	tick := time.Tick(time.Duration(fm.conf.CancellationTimeout) * time.Minute)
	for range tick {
		fm.log.Println("info: starting cancellation...")
		if err := fm.userTransactionRepo.CancelLastNTransactions(fm.conf.Num2Cancel); err != nil {
			fm.log.Printf("error: can't cancel last %d transactions: %+v\n", fm.conf.Num2Cancel, err)
			continue
		}
		fm.log.Printf("info: last %d transactions successfully canceled\n", fm.conf.Num2Cancel)
	}
	return nil
}

func NewFinanceManager(log *log.Logger, conf *Config, uc cache.UserCacheI, utr user_transaction.UserTransactionRepoI) FinanceManagerI {
	return &FinanceManager{
		conf:                conf,
		log:                 log,
		userCache:           uc,
		userTransactionRepo: utr,
	}
}
