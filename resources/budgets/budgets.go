package budgets

import (
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
