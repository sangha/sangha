package payments

import (
	"time"

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
	ID                  int64     `json:"id"`
	UserID              int64     `json:"user_id"`
	Amount              int64     `json:"amount"`
	Currency            string    `json:"currency"`
	Code                string    `json:"code"`
	Description         string    `json:"description"`
	Source              string    `json:"source"`
	SourceID            string    `json:"source_id"`
	SourcePayerID       string    `json:"source_payer_id"`
	SourceTransactionID string    `json:"source_transaction_id"`
	CreatedAt           time.Time `json:"created_at"`
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
		ID:                  payment.ID,
		UserID:              payment.UserID,
		Amount:              payment.Amount,
		Currency:            payment.Currency,
		Description:         payment.Description,
		Source:              payment.Source,
		SourceID:            payment.SourceID,
		SourcePayerID:       payment.SourcePayerID,
		SourceTransactionID: payment.SourceTransactionID,
		CreatedAt:           payment.CreatedAt,
	}

	return resp
}
