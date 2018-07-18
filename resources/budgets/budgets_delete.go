package budgets

import (
	"net/http"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// DeleteAuthRequired returns true because all requests need authentication
func (r *BudgetResource) DeleteAuthRequired() bool {
	return true
}

// DeleteDoc returns the description of this API endpoint
func (r *BudgetResource) DeleteDoc() string {
	return "delete a budget"
}

// DeleteParams returns the parameters supported by this API endpoint
func (r *BudgetResource) DeleteParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *BudgetResource) Delete(context smolder.APIContext, request *restful.Request, response *restful.Response) {
	auth, err := context.Authentication(request)
	if err != nil || auth.(db.Budget).ID != 1 {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Admin permission required for this operation",
			"BudgetResource DELETE"))
		return
	}

	ctx := context.(*db.APIContext)
	budget, err := ctx.GetBudgetByUUID(request.PathParameter("budget-id"))
	if err != nil {
		r.NotFound(request, response)
		return
	}

	err = budget.Delete(ctx)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't delete budget",
			"BudgetResource DELETE"))
		return
	}

	resp := BudgetResponse{}
	resp.Init(context)
	resp.Send(response)
}
