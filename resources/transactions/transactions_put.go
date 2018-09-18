package transactions

import (
	"fmt"
	"strconv"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// TransactionPutStruct holds all values of an incoming PUT request
type TransactionPutStruct struct {
	TransactionPostStruct
}

// PutAuthRequired returns true because all requests need authentication
func (r *TransactionResource) PutAuthRequired() bool {
	return false
}

// PutDoc returns the description of this API endpoint
func (r *TransactionResource) PutDoc() string {
	return "update an existing transaction"
}

// PutParams returns the parameters supported by this API endpoint
func (r *TransactionResource) PutParams() []*restful.Parameter {
	return nil
}

// Put processes an incoming PUT (update) request
func (r *TransactionResource) Put(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	resp := TransactionResponse{}
	resp.Init(context)

	pps := data.(*TransactionPostStruct)
	tid, err := strconv.Atoi(request.PathParameter("transaction-id"))
	if err != nil {
		r.NotFound(request, response)
		return
	}

	transaction, err := context.(*db.APIContext).LoadTransactionByID(int64(tid))
	if err != nil {
		r.NotFound(request, response)
		return
	}

	/*	auth, err := context.Authentication(request)
		if err != nil || (auth.(db.User).ID != 1 && auth.(db.User).ID != transaction.UserID) {
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusUnauthorized,
				false,
				"Admin permission required for this operation",
				"TransactionResource PUT"))
			return
		} */

	fmt.Println("PPS:", pps)
	// transaction.Pending = pps.Transaction.Pending

	/*	err = transaction.Update(context.(*db.APIContext))
		if err != nil {
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusInternalServerError,
				true,
				"Can't update transaction",
				"TransactionResource PUT"))
			return
		} */

	resp.AddTransaction(transaction)
	resp.Send(response)
}
