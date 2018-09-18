package transactions

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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
	params = append(params, restful.QueryParameter("limit", "returns at most n transactions").DataType("int"))
	params = append(params, restful.QueryParameter("pending", "returns only pending transactions").DataType("bool"))
	params = append(params, restful.QueryParameter("direction", "returns only 'incoming' or 'outgoing' transactions").DataType("string"))
	params = append(params, restful.QueryParameter("project", "returns transactions for a specific project only").DataType("string"))
	params = append(params, restful.QueryParameter("budget", "returns transactions for a specific budget only").DataType("string"))
	params = append(params, restful.QueryParameter("donor", "returns transactions for a specific donor only").DataType("string"))
	params = append(params, restful.QueryParameter("from_date", "returns transactions starting with a specific date").DataType("string"))
	params = append(params, restful.QueryParameter("to_date", "returns transactions up to a specific date").DataType("string"))
	params = append(params, restful.QueryParameter("search", "returns transactions matching a search term only").DataType("string"))

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

		resp.AddTransaction(transaction)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *TransactionResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	ctx := context.(*db.APIContext)
	resp := TransactionResponse{}
	resp.Init(context)

	var err error
	var transactions []db.Transaction

	if len(params["project"]) > 0 {
		var project db.Project
		project, err = context.(*db.APIContext).LoadProjectByUUID(params["project"][0])
		if err != nil {
			r.NotFound(request, response)
			return
		}

		transactions, err = project.LoadTransactions(ctx)
		if err != nil {
			smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
				http.StatusInternalServerError,
				"Can't load transactions",
				"TransactionsResource GET"))
			return
		}
	} else if len(params["budget"]) > 0 {
		var budget db.Budget
		budget, err = ctx.LoadBudgetByUUID(params["budget"][0])
		if err != nil {
			r.NotFound(request, response)
			return
		}

		transactions, err = budget.LoadTransactions(ctx)
		if err != nil {
			smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
				http.StatusInternalServerError,
				"Can't load transactions",
				"TransactionsResource GET"))
			return
		}
	}

	var from, to time.Time
	timeLayout := "2006-01-02"

	if len(params["from_date"]) > 0 {
		from, _ = time.Parse(timeLayout, params["from_date"][0])
	}
	if len(params["to_date"]) > 0 {
		to, err = time.Parse(timeLayout, params["to_date"][0])
		if err == nil {
			to = to.Add(time.Hour * 24)
		}
	}
	var ts []db.Transaction
	for _, t := range transactions {
		if !to.IsZero() && t.CreatedAt.After(to) {
			continue
		}
		if !from.IsZero() && t.CreatedAt.Before(from) {
			continue
		}

		ts = append(ts, t)
	}
	transactions = ts

	ts = []db.Transaction{}
	if len(params["search"]) > 0 && params["search"][0] != "" {
		search := strings.ToLower(params["search"][0])
		for _, t := range transactions {
			if !strings.Contains(strings.ToLower(t.Purpose), search) {
				continue
			}

			ts = append(ts, t)
		}
		transactions = ts
	}

	n := 0
	if len(params["limit"]) > 0 {
		limit, _ := strconv.ParseInt(params["limit"][0], 10, 0)
		n = len(transactions) - int(limit)
		if n < 0 {
			n = 0
		}
	}
	for _, transaction := range transactions[n:] {
		resp.AddTransaction(transaction)
	}

	resp.Send(response)
}
