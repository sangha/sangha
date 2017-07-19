package budgets

import (
	"errors"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// BudgetResource is the resource responsible for /budgets
type BudgetResource struct {
	smolder.Resource
}

var (
	_ smolder.GetIDSupported = &BudgetResource{}
	_ smolder.GetSupported   = &BudgetResource{}
	_ smolder.PostSupported  = &BudgetResource{}
	_ smolder.PutSupported   = &BudgetResource{}
)

// Register this resource with the container to setup all the routes
func (r *BudgetResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "BudgetResource"
	r.TypeName = "budget"
	r.Endpoint = "budgets"
	r.Doc = "Manage budgets"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// Reads returns the model that will be read by POST, PUT & PATCH operations
func (r *BudgetResource) Reads() interface{} {
	return BudgetPostStruct{}
}

// Returns returns the model that will be returned
func (r *BudgetResource) Returns() interface{} {
	return BudgetResponse{}
}

// Validate checks an incoming request for data errors
func (r *BudgetResource) Validate(context smolder.APIContext, data interface{}, request *restful.Request) error {
	ups := data.(BudgetPostStruct)

	if strings.TrimSpace(ups.Budget.Name) == "" {
		return errors.New("Invalid budget name")
	}

	return nil
}
