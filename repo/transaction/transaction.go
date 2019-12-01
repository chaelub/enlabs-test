package transaction

import (
	"database/sql"
	"enlabs-test/store/pg"
	"fmt"
	"github.com/shopspring/decimal"
)

var TransactionStates = [2]string{"lose", "win"}

type TransactionStatus uint8
type TransactionState uint8

const (
	TxUnknown TransactionStatus = iota
	TXFail
	TxSuccess
	TXCanceled
)

const (
	TxStateUnknown TransactionState = iota
	TxStateLose
	TxStateWin
)

type TransactionRepoI interface {
	LastNSuccessful(n int) ([]int64, error)
	TransactionByExtId(string) (*Transaction, error)
	CheckExistsByExtId(string) (bool, error)
	Insert(*Transaction) error
	Update(*Transaction) error
}

// for future needs
type Transaction struct {
	Id     int64
	ExtId  string
	UserId int64
	State  TransactionState
	Amount decimal.Decimal
	Tms    int64
	Status TransactionStatus
}

type TransactionRepo struct {
	store *sql.DB
}

func (t *TransactionRepo) TransactionByExtId(extId string) (*Transaction, error) {
	row := pg.PGStore.QueryRow(fmt.Sprintf(transactionByExtId, extId))
	transaction := new(Transaction)
	err := row.Scan(
		&transaction.Id,
		&transaction.ExtId,
		&transaction.Amount,
		&transaction.State,
		&transaction.Status,
		&transaction.UserId,
		&transaction.Tms,
	)

	return transaction, err
}

func (t *TransactionRepo) CheckExistsByExtId(extId string) (bool, error) {
	row := pg.PGStore.QueryRow(fmt.Sprintf(checkExistsByExtId, extId))
	count := 0
	err := row.Scan(&count)

	return count > 0, err
}

func (t *TransactionRepo) Insert(transaction *Transaction) error {
	return nil
}

func (t *TransactionRepo) Update(transaction *Transaction) error {
	return nil
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
