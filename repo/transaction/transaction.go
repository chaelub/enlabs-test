package transaction

import (
	"database/sql"
	"enlabs-test/store/pg"
	"fmt"
	"github.com/shopspring/decimal"
)

var TransactionStates = [2]string{"lose", "win"}

type TransactionStatus uint8

const (
	TxUnknown TransactionStatus = iota
	TXFail
	TxSuccess
	TXCanceled
)

type TransactionRepoI interface {
	LastNSuccessful(n int) ([]int64, error)
}

// for future needs
type Transaction struct {
	Id     int64           `json:"id"`
	ExtId  string          `json:"transactionId"`
	State  string          `json:"state"`
	Amount decimal.Decimal `json:"amount"`
}

type TransactionRepo struct {
	store *sql.DB
}

func (t *TransactionRepo) LastNSuccessful(n int) ([]int64, error) {
	return t.lastNByStatus(n, TxSuccess)
}

func (t *TransactionRepo) lastNByStatus(n int, status TransactionStatus) ([]int64, error) {
	rows, err := pg.PGStore.Query(fmt.Sprintf(lastNByStatus, status))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ids := make([]int64, 0, 0)

	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func NewTransactionRepo(store *sql.DB) (TransactionRepoI, error) {
	return &TransactionRepo{store: store}, nil
}
