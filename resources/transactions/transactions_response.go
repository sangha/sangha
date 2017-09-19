package transactions

import (
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
	ID       int64   `json:"id"`
	BudgetID int64   `json:"budget_id"`
	Amount   float64 `json:"amount"`
}

// Init a new response
func (r *TransactionResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Transactions = []transactionInfoResponse{}
}

// AddTransaction adds a transaction to the response
func (r *TransactionResponse) AddTransaction(transaction *db.Transaction) {
	r.transactions = append(r.transactions, *transaction)
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

func prepareTransactionResponse(context smolder.APIContext, transaction *db.Transaction) transactionInfoResponse {
	resp := transactionInfoResponse{
		ID:       transaction.ID,
		BudgetID: transaction.BudgetID,
		Amount:   transaction.Amount,
	}

	return resp
}
