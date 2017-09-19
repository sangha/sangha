package payments

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/muesli/smolder"
)

// PaymentResponse is the common response to 'payment' requests
type PaymentResponse struct {
	smolder.Response

	Payments []paymentInfoResponse `json:"payments,omitempty"`
	payments []db.Payment
}

type paymentInfoResponse struct {
	ID     int64   `json:"id"`
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

// Init a new response
func (r *PaymentResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Payments = []paymentInfoResponse{}
}

// AddPayment adds a payment to the response
func (r *PaymentResponse) AddPayment(payment *db.Payment) {
	r.payments = append(r.payments, *payment)
	r.Payments = append(r.Payments, preparePaymentResponse(r.Context, payment))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *PaymentResponse) EmptyResponse() interface{} {
	if len(r.payments) == 0 {
		var out struct {
			Payments interface{} `json:"payments"`
		}
		out.Payments = []paymentInfoResponse{}
		return out
	}
	return nil
}

func preparePaymentResponse(context smolder.APIContext, payment *db.Payment) paymentInfoResponse {
	resp := paymentInfoResponse{
		ID:     payment.ID,
		UserID: payment.UserID,
		Amount: payment.Amount,
	}

	return resp
}
