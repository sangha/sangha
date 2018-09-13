package payments

import (
	"errors"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// PaymentResource is the resource responsible for /payments
type PaymentResource struct {
	smolder.Resource
}

var (
	_ smolder.GetSupported  = &PaymentResource{}
	_ smolder.PostSupported = &PaymentResource{}
	_ smolder.PutSupported  = &PaymentResource{}
)

// Register this resource with the container to setup all the routes
func (r *PaymentResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "PaymentResource"
	r.TypeName = "payment"
	r.Endpoint = "payments"
	r.Doc = "Manage payments"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// Reads returns the model that will be read by POST, PUT & PATCH operations
func (r *PaymentResource) Reads() interface{} {
	return &PaymentPostStruct{}
}

// Returns returns the model that will be returned
func (r *PaymentResource) Returns() interface{} {
	return PaymentResponse{}
}

// Validate checks an incoming request for data errors
func (r *PaymentResource) Validate(context smolder.APIContext, data interface{}, request *restful.Request) error {
	ups := data.(*PaymentPostStruct)

	if ups.Payment.Amount == 0 {
		return errors.New("Invalid payment amount")
	}

	return nil
}
