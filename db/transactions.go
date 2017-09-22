package db

import "time"

// Transaction represents the db schema of a transaction
type Transaction struct {
	ID           int64
	BudgetID     int64
	FromBudgetID *int64
	ToBudgetID   *int64
	Amount       float64
	CreatedAt    time.Time
}

var ()

// LoadTransactionByID loads a transaction by ID from the database
func (context *APIContext) LoadTransactionByID(id int64) (Transaction, error) {
	transaction := Transaction{}
	if id < 1 {
		return transaction, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, budget_id, from_budget_id, to_budget_id, amount, created_at FROM transactions WHERE id = $1", id).
		Scan(&transaction.ID, &transaction.BudgetID, &transaction.FromBudgetID, &transaction.ToBudgetID, &transaction.Amount, &transaction.CreatedAt)
	return transaction, err
}

// Save a transaction to the database
func (transaction *Transaction) Save(context *APIContext) error {
	err := context.QueryRow("INSERT INTO transactions (budget_id, from_budget_id, to_budget_id, amount, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		transaction.BudgetID, transaction.FromBudgetID, transaction.ToBudgetID, transaction.Amount, time.Now()).Scan(&transaction.ID)
	return err
}
