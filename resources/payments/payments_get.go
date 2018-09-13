package payments

import (
	"net/http"
	"strconv"
	"strings"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *PaymentResource) GetAuthRequired() bool {
	return true
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *PaymentResource) GetByIDsAuthRequired() bool {
	return true
}

// GetDoc returns the description of this API endpoint
func (r *PaymentResource) GetDoc() string {
	return "retrieve payments"
}

// GetParams returns the parameters supported by this API endpoint
func (r *PaymentResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	params = append(params, restful.QueryParameter("limit", "returns at most n payments").DataType("int"))
	params = append(params, restful.QueryParameter("direction", "returns only 'incoming' or 'outgoing' payments").DataType("string"))
	params = append(params, restful.QueryParameter("donor", "returns payments for a specific donor only").DataType("string"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *PaymentResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := PaymentResponse{}
	resp.Init(context)

	auth, err := context.Authentication(request)
	if err != nil || auth == nil || auth.(db.User).ID != 1 {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			"Admin permission required for this operation",
			"PaymentResource GET"))
		return
	}

	for _, id := range ids {
		iid, _ := strconv.ParseInt(id, 10, 0)
		payment, err := context.(*db.APIContext).LoadPaymentByID(iid)
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddPayment(payment)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *PaymentResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	ctx := context.(*db.APIContext)
	resp := PaymentResponse{}
	resp.Init(context)

	auth, err := context.Authentication(request)
	if err != nil || auth == nil || auth.(db.User).ID != 1 {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			"Admin permission required for this operation",
			"PaymentResource GET"))
		return
	}

	var payments []db.Payment

	if len(params["budget"]) > 0 {
		var budget db.Budget
		budget, err = ctx.LoadBudgetByUUID(params["budget"][0])
		if err != nil {
			r.NotFound(request, response)
			return
		}

		payments, err = budget.LoadPayments(ctx)
		if err != nil {
			smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
				http.StatusInternalServerError,
				"Can't load payments",
				"PaymentsResource GET"))
			return
		}
	} else if len(params["donor"]) > 0 {
		payments, err = ctx.LoadPaymentsForDonor(params["donor"][0])
		if err != nil {
			smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
				http.StatusInternalServerError,
				"Can't load payments",
				"PaymentsResource GET"))
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
		payments, err = ctx.LoadPendingPayments(direction)
		if err != nil {
			smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
				http.StatusInternalServerError,
				"Can't load payments",
				"PaymentsResource GET"))
			return
		}
	}

	n := 0
	if len(params["limit"]) > 0 {
		limit, _ := strconv.ParseInt(params["limit"][0], 10, 0)
		n = len(payments) - int(limit)
		if n < 0 {
			n = 0
		}
	}
	for _, payment := range payments[n:] {
		resp.AddPayment(payment)
	}

	resp.Send(response)
}
