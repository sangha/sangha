package db

import (
	"time"
)

// Transaction represents the db schema of a transaction
type Transaction struct {
	ID           int64
	BudgetID     int64
	FromBudgetID *int64
	ToBudgetID   *int64
	Amount       int64
	CreatedAt    time.Time
	Purpose      string
	PaymentID    *int64
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

	err := context.QueryRow("SELECT id, budget_id, from_budget_id, to_budget_id, amount, created_at, purpose, payment_id "+
		"FROM transactions "+
		"WHERE id = $1", id).
		Scan(&transaction.ID, &transaction.BudgetID, &transaction.FromBudgetID, &transaction.ToBudgetID, &transaction.Amount,
			&transaction.CreatedAt, &transaction.Purpose, &transaction.PaymentID)

	return transaction, err
}

// LoadTransactions loads all transactions for a project
func (project *Project) LoadTransactions(context *APIContext) ([]Transaction, error) {
	b, err := context.LoadRootBudgetForProject(project)
	if err != nil {
		return []Transaction{}, err
	}

	return b.LoadTransactions(context)
	/*	transactions := []Transaction{}

		rows, err := context.Query("SELECT transactions.id, transactions.budget_id, transactions.from_budget_id, transactions.to_budget_id, transactions.amount, transactions.created_at, transactions.purpose, transactions.payment_id "+
			"FROM transactions, budgets "+
			"WHERE transactions.budget_id = budgets.id AND budgets.project_id = $1"+
			"ORDER BY created_at, id ASC", project.ID)

		if err != nil {
			return transactions, err
		}

		defer rows.Close()
		for rows.Next() {
			transaction := Transaction{}
			err = rows.Scan(&transaction.ID, &transaction.BudgetID, &transaction.FromBudgetID, &transaction.ToBudgetID, &transaction.Amount,
				&transaction.CreatedAt, &transaction.Purpose, &transaction.PaymentID)
			if err != nil {
				return transactions, err
			}

			transactions = append(transactions, transaction)
		}

		return transactions, err*/
}

// LoadTransactions loads all transactions for a budget
func (budget *Budget) LoadTransactions(context *APIContext) ([]Transaction, error) {
	transactions := []Transaction{}

	rows, err := context.Query("SELECT id, budget_id, from_budget_id, to_budget_id, amount, created_at, purpose, payment_id "+
		"FROM transactions "+
		"WHERE budget_id = $1 "+
		"ORDER BY created_at, id ASC", budget.ID)

	if err != nil {
		return transactions, err
	}

	defer rows.Close()
	for rows.Next() {
		transaction := Transaction{}
		err = rows.Scan(&transaction.ID, &transaction.BudgetID, &transaction.FromBudgetID, &transaction.ToBudgetID, &transaction.Amount,
			&transaction.CreatedAt, &transaction.Purpose, &transaction.PaymentID)
		if err != nil {
			return transactions, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, err
}

// Save a transaction to the database
func (transaction *Transaction) Save(context *APIContext) error {
	err := context.QueryRow("INSERT INTO transactions (budget_id, from_budget_id, to_budget_id, amount, created_at, purpose, payment_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		transaction.BudgetID, transaction.FromBudgetID, transaction.ToBudgetID, transaction.Amount, transaction.CreatedAt, transaction.Purpose, transaction.PaymentID).Scan(&transaction.ID)
	return err
}

func (context *APIContext) Transfer(fromBudget, toBudget int64, amount int64, purpose string, paymentID int64, ts time.Time) (Transaction, error) {
	if amount < 0 {
		fromBudget, toBudget = toBudget, fromBudget
		amount *= -1
	}

	// FIXME: move to transaction
	torig := Transaction{
		BudgetID:   fromBudget,
		ToBudgetID: &toBudget,
		Amount:     -amount,
		CreatedAt:  ts, // FIXME: time.Now().UTC(),
		Purpose:    purpose,
	}
	if paymentID > 0 {
		torig.PaymentID = &paymentID
	}
	if err := torig.Save(context); err != nil {
		return Transaction{}, err
	}

	t := Transaction{
		BudgetID:     toBudget,
		FromBudgetID: &fromBudget,
		Amount:       amount,
		CreatedAt:    ts, // FIXME: time.Now().UTC(),
		Purpose:      purpose,
	}
	if paymentID > 0 {
		t.PaymentID = &paymentID
	}
	return torig, t.Save(context)
}
