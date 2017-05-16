package budgets

import (
	"net/http"

	"gitlab.techcultivation.org/techcultivation/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// BudgetPostStruct holds all values of an incoming POST request
type BudgetPostStruct struct {
	Budget struct {
		ProjectID int64  `json:"project_id"`
		Name      string `json:"name"`
	} `json:"budget"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *BudgetResource) PostAuthRequired() bool {
	return false
}

// PostDoc returns the description of this API endpoint
func (r *BudgetResource) PostDoc() string {
	return "create a new budget invitation"
}

// PostParams returns the parameters supported by this API endpoint
func (r *BudgetResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *BudgetResource) Post(context smolder.APIContext, request *restful.Request, response *restful.Response) {
	/*auth, err := context.Authentication(request)
		if err != nil || auth.(db.Budget).ID != 1 {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Admin permission required for this operation",
			"BudgetResource POST"))
		return
	}*/

	ups := BudgetPostStruct{}
	err := request.ReadEntity(&ups)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse POST data",
			"BudgetResource POST"))
		return
	}

	budget := db.Budget{
		ProjectID: ups.Budget.ProjectID,
		Name:      ups.Budget.Name,
	}
	err = budget.Save(context.(*db.APIContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't create budget",
			"BudgetResource POST"))
		return
	}

	resp := BudgetResponse{}
	resp.Init(context)
	resp.AddBudget(&budget)
	resp.Send(response)
}
