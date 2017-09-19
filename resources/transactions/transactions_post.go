package transactions

import (
	"log"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// TransactionPostStruct holds all values of an incoming POST request
type TransactionPostStruct struct {
	Transaction struct {
		Source          string  `json:"type"`
		SourceID        string  `json:"source_id"`
		Amount          float64 `json:"amount"`
		TransactionCode string  `json:"transaction_code"`
	} `json:"transaction"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *TransactionResource) PostAuthRequired() bool {
	return true
}

// PostDoc returns the description of this API endpoint
func (r *TransactionResource) PostDoc() string {
	return "start a new transaction"
}

// PostParams returns the parameters supported by this API endpoint
func (r *TransactionResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *TransactionResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	ups := data.(*TransactionPostStruct)
	log.Printf("Got transaction request: %+v\n", ups)

	transaction := db.Transaction{
		Amount: ups.Transaction.Amount,
	}

	resp := TransactionResponse{}
	resp.Init(context)
	resp.AddTransaction(&transaction)
	resp.Send(response)
}
