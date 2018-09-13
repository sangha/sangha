package payments

import (
	"strconv"
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
	BudgetID            string    `json:"budget_id"`
	CreatedAt           time.Time `json:"created_at"`
	Amount              int64     `json:"amount"`
	Currency            string    `json:"currency"`
	Code                *string   `json:"code"`
	Purpose             string    `json:"purpose"`
	RemoteAccount       string    `json:"remote_account"`
	RemoteBankID        string    `json:"remote_bank_id"`
	RemoteTransactionID string    `json:"remote_transaction_id"`
	RemoteName          string    `json:"remote_name"`
	Source              string    `json:"source"`
	Pending             bool      `json:"pending"`
}

// Init a new response
func (r *PaymentResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Payments = []paymentInfoResponse{}
}

// AddPayment adds a payment to the response
func (r *PaymentResponse) AddPayment(payment db.Payment) {
	r.payments = append(r.payments, payment)
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

func preparePaymentResponse(context smolder.APIContext, payment db.Payment) paymentInfoResponse {
	resp := paymentInfoResponse{
		ID:                  payment.ID,
		CreatedAt:           payment.CreatedAt,
		Amount:              payment.Amount,
		Currency:            payment.Currency,
		Purpose:             payment.Purpose,
		RemoteAccount:       payment.RemoteAccount,
		RemoteBankID:        payment.RemoteBankID,
		RemoteTransactionID: payment.RemoteTransactionID,
		RemoteName:          payment.RemoteName,
		Source:              payment.Source,
		Pending:             payment.Pending,
	}

	if payment.Code != "" {
		resp.Code = &payment.Code
	}

	c, err := context.(*db.APIContext).LoadCodeByCode(payment.Code)
	if err == nil {
		bid, err := strconv.ParseInt(c.BudgetIDs[0], 10, 64)
		if err != nil {
			panic(err)
		}
		b, err := context.(*db.APIContext).LoadBudgetByID(bid)
		if err != nil {
			panic(err)
		}

		resp.BudgetID = b.UUID
	}

	return resp
}
