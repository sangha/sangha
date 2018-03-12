package transactions

import (
	"errors"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// TransactionResource is the resource responsible for /transactions
type TransactionResource struct {
	smolder.Resource
}

var (
	_ smolder.GetSupported  = &TransactionResource{}
	_ smolder.PostSupported = &TransactionResource{}
)

// Register this resource with the container to setup all the routes
func (r *TransactionResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "TransactionResource"
	r.TypeName = "transaction"
	r.Endpoint = "transactions"
	r.Doc = "Manage transactions"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// Reads returns the model that will be read by POST, PUT & PATCH operations
func (r *TransactionResource) Reads() interface{} {
	return &TransactionPostStruct{}
}

// Returns returns the model that will be returned
func (r *TransactionResource) Returns() interface{} {
	return TransactionResponse{}
}

// Validate checks an incoming request for data errors
func (r *TransactionResource) Validate(context smolder.APIContext, data interface{}, request *restful.Request) error {
	ups := data.(*TransactionPostStruct)

	if ups.Transaction.Amount <= 0 {
		return errors.New("Invalid transaction amount")
	}

	return nil
}
