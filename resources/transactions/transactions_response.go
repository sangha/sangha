package transactions

import (
	"time"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/muesli/smolder"
)

// TransactionResponse is the common response to 'transaction' requests
type TransactionResponse struct {
	smolder.Response

	Transactions []transactionInfoResponse `json:"transactions,omitempty"`
	transactions []db.Transaction
}

type transactionInfoResponse struct {
	ID           int64     `json:"id"`
	BudgetID     string    `json:"budget_id"`
	FromBudgetID *string   `json:"from_budget_id"`
	ToBudgetID   *string   `json:"to_budget_id"`
	Amount       int64     `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
	Purpose      string    `json:"purpose"`
	PaymentID    *int64    `json:"payment_id"`
}

// Init a new response
func (r *TransactionResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Transactions = []transactionInfoResponse{}
}

// AddTransaction adds a transaction to the response
func (r *TransactionResponse) AddTransaction(transaction db.Transaction) {
	r.transactions = append(r.transactions, transaction)
	r.Transactions = append(r.Transactions, prepareTransactionResponse(r.Context, transaction))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *TransactionResponse) EmptyResponse() interface{} {
	if len(r.transactions) == 0 {
		var out struct {
			Transactions interface{} `json:"transactions"`
		}
		out.Transactions = []transactionInfoResponse{}
		return out
	}
	return nil
}

func prepareTransactionResponse(context smolder.APIContext, transaction db.Transaction) transactionInfoResponse {
	ctx := context.(*db.APIContext)
	resp := transactionInfoResponse{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
		Purpose:   transaction.Purpose,
	}

	if ctx.Auth != nil && ctx.Auth.ID == 1 {
		resp.PaymentID = transaction.PaymentID
	}

	budget, _ := ctx.LoadBudgetByID(transaction.BudgetID)
	resp.BudgetID = budget.UUID
	if transaction.FromBudgetID != nil {
		fromBudget, _ := ctx.LoadBudgetByID(*transaction.FromBudgetID)
		resp.FromBudgetID = &fromBudget.UUID
	}
	if transaction.ToBudgetID != nil {
		toBudget, _ := ctx.LoadBudgetByID(*transaction.ToBudgetID)
		resp.ToBudgetID = &toBudget.UUID
	}

	return resp
}
