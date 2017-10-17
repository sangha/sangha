package payments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// PaymentPostStruct holds all values of an incoming POST request
type PaymentPostStruct struct {
	Payment struct {
		Source          string  `json:"source"`
		SourceID        string  `json:"source_id"`
		Amount          float64 `json:"amount"`
		TransactionCode string  `json:"transaction_code"`
	} `json:"payment"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *PaymentResource) PostAuthRequired() bool {
	return false
}

// PostDoc returns the description of this API endpoint
func (r *PaymentResource) PostDoc() string {
	return "start a new payment"
}

// PostParams returns the parameters supported by this API endpoint
func (r *PaymentResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *PaymentResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	ups := data.(*PaymentPostStruct)
	log.Printf("Got payment request: %+v\n", ups)

	ctx := context.(*db.APIContext)
	payments := struct {
		Payments []struct {
			UserID              int64     `json:"user_id"`
			Amount              float64   `json:"amount"`
			Currency            string    `json:"currency"`
			Code                string    `json:"code"`
			Description         string    `json:"description"`
			Source              string    `json:"source"`
			SourceID            string    `json:"source_id"`
			SourcePayerID       string    `json:"source_payer_id"`
			SourceTransactionID string    `json:"source_transaction_id"`
			CreatedAt           time.Time `json:"created_at"`
		} `json:"payments"`
	}{}
	payment := db.Payment{}

	switch ups.Payment.Source {
	case "bank_transfer":
		payment.Amount = ups.Payment.Amount
		payment.Code = ups.Payment.TransactionCode
		payment.Source = ups.Payment.Source
	case "paypal":
		resp, err := http.Get(ctx.Config.Connections.PayPal + "/" + ups.Payment.SourceID)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))

		if resp.StatusCode != http.StatusOK {
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusBadRequest,
				false,
				"Unknown payment ID",
				"PaymentResource POST"))
			return
		}

		err = json.Unmarshal(body, &payments)
		if err != nil {
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusInternalServerError,
				true,
				"Error decoding payment response",
				"PaymentResource POST"))
			return
		}
		fmt.Printf("Payment: %+v\n", payments.Payments[0])
		payment.UserID = payments.Payments[0].UserID
		payment.Amount = payments.Payments[0].Amount
		payment.Currency = payments.Payments[0].Currency
		payment.Code = payments.Payments[0].Code
		payment.Description = payments.Payments[0].Description
		payment.Source = payments.Payments[0].Source
		payment.SourceID = payments.Payments[0].SourceID
		payment.SourcePayerID = payments.Payments[0].SourcePayerID
		payment.SourceTransactionID = payments.Payments[0].SourceTransactionID
		payment.CreatedAt = payments.Payments[0].CreatedAt
	default:
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Unknown payment source",
			"PaymentResource POST"))
		return
	}

	resp := PaymentResponse{}
	resp.Init(context)
	resp.AddPayment(&payment)
	resp.Send(response)
}
