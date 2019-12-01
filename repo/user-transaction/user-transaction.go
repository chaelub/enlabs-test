package user_transaction

import (
	"database/sql"
	"enlabs-test/repo/transaction"
	"enlabs-test/repo/user"
	"time"
)

type UserTransactionRepoI interface {
	user.UserRepoI
	transaction.TransactionRepoI
	ProcessTransaction(*user.User, *transaction.Transaction) error
	CancelLastNTransactions(int) error
}

type UserTransactionRepo struct {
	user.UserRepoI
	transaction.TransactionRepoI
	store *sql.DB
}

func (ut *UserTransactionRepo) ProcessTransaction(u *user.User, t *transaction.Transaction) error {
	tx, err := ut.store.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(processTransaction)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(t.ExtId, t.Amount, t.State, t.Status, u.Id, time.Now().Unix())
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ut *UserTransactionRepo) CancelLastNTransactions(n int) error {
	tx, err := ut.store.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(cancelLastNTransactions)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = stmt.Exec(n * 2)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func NewUserTransactionRepo(store *sql.DB, ur user.UserRepoI, tr transaction.TransactionRepoI) UserTransactionRepoI {
	return &UserTransactionRepo{
		store:            store,
		UserRepoI:        ur,
		TransactionRepoI: tr,
	}
}
