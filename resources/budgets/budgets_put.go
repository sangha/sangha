package budgets

import (
	"net/http"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// BudgetPutStruct holds all values of an incoming PUT request
type BudgetPutStruct struct {
	BudgetPostStruct
}

// PutAuthRequired returns true because all requests need authentication
func (r *BudgetResource) PutAuthRequired() bool {
	return true
}

// PutDoc returns the description of this API endpoint
func (r *BudgetResource) PutDoc() string {
	return "update an existing budget"
}

// PutParams returns the parameters supported by this API endpoint
func (r *BudgetResource) PutParams() []*restful.Parameter {
	return nil
}

// Put processes an incoming PUT (update) request
func (r *BudgetResource) Put(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	resp := BudgetResponse{}
	resp.Init(context)

	budget, err := context.(*db.APIContext).GetBudgetByUUID(request.PathParameter("budget-id"))
	if err != nil {
		r.NotFound(request, response)
		return
	}

	/*	auth, err := context.Authentication(request)
		if err != nil || (auth.(db.User).ID != 1 && auth.(db.User).ID != budget.UserID) {
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusUnauthorized,
				false,
				"Admin permission required for this operation",
				"BudgetResource PUT"))
			return
		} */

	pps := data.(*BudgetPutStruct)
	budget.ProjectID = &pps.Budget.ProjectID
	budget.Name = pps.Budget.Name
	budget.Private = pps.Budget.Private
	budget.PrivateBalance = pps.Budget.PrivateBalance

	err = budget.Update(context.(*db.APIContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't update budget",
			"BudgetResource PUT"))
		return
	}

	resp.AddBudget(&budget)
	resp.Send(response)
}
