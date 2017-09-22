package budgets

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *BudgetResource) GetAuthRequired() bool {
	return false
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *BudgetResource) GetByIDsAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *BudgetResource) GetDoc() string {
	return "retrieve budgets"
}

// GetParams returns the parameters supported by this API endpoint
func (r *BudgetResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	params = append(params, restful.QueryParameter("name", "name of a budget").DataType("string"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *BudgetResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := BudgetResponse{}
	resp.Init(context)

	for _, id := range ids {
		budget, err := context.(*db.APIContext).GetBudgetByUUID(id)
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddBudget(&budget)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *BudgetResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	resp := BudgetResponse{}
	resp.Init(context)

	budgets, err := context.(*db.APIContext).LoadAllBudgets()
	if err != nil {
		r.NotFound(request, response)
		return
	}

	for _, budget := range budgets {
		resp.AddBudget(&budget)
	}

	resp.Send(response)
}
