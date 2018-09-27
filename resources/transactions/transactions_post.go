package transactions

import (
	"log"
	"net/http"
	"time"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// TransactionPostStruct holds all values of an incoming POST request
type TransactionPostStruct struct {
	Transaction struct {
		BudgetID   string `json:"budget_id"`
		ToBudgetID string `json:"to_budget_id"`
		Amount     int64  `json:"amount"`
		Purpose    string `json:"purpose"`
	} `json:"transaction"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *TransactionResource) PostAuthRequired() bool {
	return true
}

// PostDoc returns the description of this API endpoint
func (r *TransactionResource) PostDoc() string {
	return "start a new transaction"
}

// PostParams returns the parameters supported by this API endpoint
func (r *TransactionResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *TransactionResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	auth, err := context.Authentication(request)
	if err != nil || auth.(db.User).ID != 1 {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			"Admin permission required for this operation",
			"TransactionResource POST"))
		return
	}

	ctx := context.(*db.APIContext)
	ups := data.(*TransactionPostStruct)
	log.Printf("Got transaction request: %+v\n", ups)

	from, err := ctx.LoadBudgetByUUID(ups.Transaction.BudgetID)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusBadRequest,
			"A budget with this ID does not exist",
			"TransactionResource POST"))
		return
	}
	to, err := ctx.LoadBudgetByUUID(ups.Transaction.ToBudgetID)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusBadRequest,
			"A budget with this ID does not exist",
			"TransactionResource POST"))
		return
	}

	bal, err := from.Balance(ctx)
	if err != nil {
		panic(err)
	}
	if bal < ups.Transaction.Amount {
		smolder.ErrorResponseHandler(request, response, nil, smolder.NewErrorResponse(
			http.StatusBadRequest,
			"This budget does not have the necessary funds",
			"TransactionResource POST"))
		return
	}

	t, err := ctx.Transfer(from.ID, to.ID, ups.Transaction.Amount, ups.Transaction.Purpose, 0, time.Now().UTC())
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			"Could not process transaction",
			"TransactionResource POST"))
		return
	}

	resp := TransactionResponse{}
	resp.Init(context)
	resp.AddTransaction(t)
	resp.Send(response)
}
