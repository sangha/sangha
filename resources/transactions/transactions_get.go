package transactions

import (
	"strconv"
	"strings"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *TransactionResource) GetAuthRequired() bool {
	return false
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *TransactionResource) GetByIDsAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *TransactionResource) GetDoc() string {
	return "retrieve transactions"
}

// GetParams returns the parameters supported by this API endpoint
func (r *TransactionResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	params = append(params, restful.QueryParameter("pending", "returns only pending transactions").DataType("bool"))
	params = append(params, restful.QueryParameter("direction", "returns only 'incoming' or 'outgoing' transactions").DataType("string"))
	params = append(params, restful.QueryParameter("project", "returns transactions for a specific project only").DataType("string"))
	params = append(params, restful.QueryParameter("donor", "returns transactions for a specific donor only").DataType("string"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *TransactionResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := TransactionResponse{}
	resp.Init(context)

	for _, id := range ids {
		iid, _ := strconv.ParseInt(id, 10, 0)
		transaction, err := context.(*db.APIContext).LoadTransactionByID(iid)
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddTransaction(&transaction)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *TransactionResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	resp := TransactionResponse{}
	resp.Init(context)

	var err error
	var project db.Project
	var transactions []db.Transaction

	if len(params["project"]) > 0 {
		project, err = context.(*db.APIContext).LoadProjectByUUID(params["project"][0])
		if err != nil {
			r.NotFound(request, response)
			return
		}

		transactions, err = project.LoadTransactions(context.(*db.APIContext))
		if err != nil {
			r.NotFound(request, response)
			return
		}
	} else if len(params["donor"]) > 0 {
		transactions, err = context.(*db.APIContext).LoadTransactionsForDonor(params["donor"][0])
		if err != nil {
			r.NotFound(request, response)
			return
		}
	} else {
		direction := db.TRANSACTION_ALL
		if len(params["direction"]) > 0 {
			if strings.ToLower(params["direction"][0]) == "incoming" {
				direction = db.TRANSACTION_INCOMING
			}
			if strings.ToLower(params["direction"][0]) == "outgoing" {
				direction = db.TRANSACTION_OUTGOING
			}
		}
		transactions, err = context.(*db.APIContext).LoadPendingTransactions(direction)
		if err != nil {
			r.NotFound(request, response)
			return
		}
	}

	for _, transaction := range transactions {
		resp.AddTransaction(&transaction)
	}

	resp.Send(response)
}
