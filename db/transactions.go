package db

import (
	"fmt"
	"time"
)

// Transaction represents the db schema of a transaction
type Transaction struct {
	ID            int64
	BudgetID      int64
	FromBudgetID  *int64
	ToBudgetID    *int64
	Amount        int64
	CreatedAt     time.Time
	RemotePurpose string
	RemoteAccount string
	RemoteBankID  string
	RemoteName    string
}

const (
	TRANSACTION_ALL = iota
	TRANSACTION_INCOMING
	TRANSACTION_OUTGOING
)

// LoadTransactionByID loads a transaction by ID from the database
func (context *APIContext) LoadTransactionByID(id int64) (Transaction, error) {
	transaction := Transaction{}
	if id < 1 {
		return transaction, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, budget_id, from_budget_id, to_budget_id, amount, created_at, remote_purpose, remote_account, remote_bank_id, remote_name FROM transactions WHERE id = $1", id).
		Scan(&transaction.ID, &transaction.BudgetID, &transaction.FromBudgetID, &transaction.ToBudgetID, &transaction.Amount, &transaction.CreatedAt, &transaction.RemotePurpose, &transaction.RemoteAccount, &transaction.RemoteBankID, &transaction.RemoteName)
	return transaction, err
}

// LoadTransactions loads all transactions for a project
func (project *Project) LoadTransactions(context *APIContext) ([]Transaction, error) {
	transactions := []Transaction{}
	budget, _ := context.LoadRootBudgetForProject(project)

	rows, err := context.Query("SELECT id, budget_id, from_budget_id, to_budget_id, amount, created_at, remote_purpose, remote_account, remote_bank_id, remote_name FROM transactions WHERE budget_id = $1 ORDER BY created_at ASC", budget.ID)
	if err != nil {
		return transactions, err
	}

	defer rows.Close()
	for rows.Next() {
		transaction := Transaction{}
		err = rows.Scan(&transaction.ID, &transaction.BudgetID, &transaction.FromBudgetID, &transaction.ToBudgetID, &transaction.Amount, &transaction.CreatedAt, &transaction.RemotePurpose, &transaction.RemoteAccount, &transaction.RemoteBankID, &transaction.RemoteName)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, err
}

// LoadTransactionsForDonor loads all transactions for a specific donor
func (context *APIContext) LoadTransactionsForDonor(donor string) ([]Transaction, error) {
	transactions := []Transaction{}

	rows, err := context.Query("SELECT id, budget_id, from_budget_id, to_budget_id, amount, created_at, remote_purpose, remote_account, remote_bank_id, remote_name FROM transactions WHERE remote_account = $1 ORDER BY created_at ASC", donor)
	if err != nil {
		return transactions, err
	}

	defer rows.Close()
	for rows.Next() {
		transaction := Transaction{}
		err = rows.Scan(&transaction.ID, &transaction.BudgetID, &transaction.FromBudgetID, &transaction.ToBudgetID, &transaction.Amount, &transaction.CreatedAt, &transaction.RemotePurpose, &transaction.RemoteAccount, &transaction.RemoteBankID, &transaction.RemoteName)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, err
}

// LoadPendingTransaction loads all pending transactions
func (context *APIContext) LoadPendingTransactions(direction int) ([]Transaction, error) {
	transactions := []Transaction{}

	var filter string
	switch direction {
	case TRANSACTION_INCOMING:
		filter = "AND amount > 0"
	case TRANSACTION_OUTGOING:
		filter = "AND amount < 0"
	}
	rows, err := context.Query(fmt.Sprintf("SELECT id, budget_id, from_budget_id, to_budget_id, amount, created_at, remote_purpose, remote_account, remote_bank_id, remote_name FROM transactions WHERE pending = true %s ORDER BY created_at ASC", filter))
	if err != nil {
		return transactions, err
	}

	defer rows.Close()
	for rows.Next() {
		transaction := Transaction{}
		err = rows.Scan(&transaction.ID, &transaction.BudgetID, &transaction.FromBudgetID, &transaction.ToBudgetID, &transaction.Amount, &transaction.CreatedAt, &transaction.RemotePurpose, &transaction.RemoteAccount, &transaction.RemoteBankID, &transaction.RemoteName)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, err
}

// Save a transaction to the database
func (transaction *Transaction) Save(context *APIContext) error {
	err := context.QueryRow("INSERT INTO transactions (budget_id, from_budget_id, to_budget_id, amount, created_at, remote_purpose, remote_account, remote_bank_id, remote_name) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		transaction.BudgetID, transaction.FromBudgetID, transaction.ToBudgetID, transaction.Amount, transaction.CreatedAt, transaction.RemotePurpose, transaction.RemoteAccount, transaction.RemoteBankID, transaction.RemoteName).Scan(&transaction.ID)
	return err
}
