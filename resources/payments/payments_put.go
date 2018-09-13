package payments

import (
	"net/http"
	"strconv"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// PaymentPutStruct holds all values of an incoming PUT request
type PaymentPutStruct struct {
	PaymentPostStruct
}

// PutAuthRequired returns true because all requests need authentication
func (r *PaymentResource) PutAuthRequired() bool {
	return false
}

// PutDoc returns the description of this API endpoint
func (r *PaymentResource) PutDoc() string {
	return "update an existing payment"
}

// PutParams returns the parameters supported by this API endpoint
func (r *PaymentResource) PutParams() []*restful.Parameter {
	return nil
}

// Put processes an incoming PUT (update) request
func (r *PaymentResource) Put(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	ctx := context.(*db.APIContext)
	resp := PaymentResponse{}
	resp.Init(context)

	iid, _ := strconv.ParseInt(request.PathParameter("payment-id"), 10, 0)

	payment, err := ctx.LoadPaymentByID(iid)
	if err != nil {
		r.NotFound(request, response)
		return
	}

	/*	auth, err := context.Authentication(request)
		if err != nil || (auth.(db.User).ID != 1 && auth.(db.User).ID != project.UserID) {
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusUnauthorized,
				false,
				"Admin permission required for this operation",
				"ProjectResource PUT"))
			return
		} */

	pps := data.(*PaymentPostStruct)
	payment.Code = pps.Payment.Code
	payment.Pending = pps.Payment.Pending

	err = payment.Update(ctx)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			"Can't update payment",
			"PaymentResource PUT"))
		return
	}

	resp.AddPayment(payment)
	resp.Send(response)
}
