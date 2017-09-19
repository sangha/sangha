package payments

import (
	"log"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// PaymentPostStruct holds all values of an incoming POST request
type PaymentPostStruct struct {
	Payment struct {
		Source          string  `json:"type"`
		SourceID        string  `json:"source_id"`
		Amount          float64 `json:"amount"`
		TransactionCode string  `json:"transaction_code"`
	} `json:"payment"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *PaymentResource) PostAuthRequired() bool {
	return true
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

	payment := db.Payment{
		Amount: ups.Payment.Amount,
	}

	resp := PaymentResponse{}
	resp.Init(context)
	resp.AddPayment(&payment)
	resp.Send(response)
}
