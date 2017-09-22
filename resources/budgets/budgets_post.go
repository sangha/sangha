package budgets

import (
	"net/http"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// BudgetPostStruct holds all values of an incoming POST request
type BudgetPostStruct struct {
	Budget struct {
		ProjectID      int64  `json:"project_id"`
		ParentID       int64  `json:"parent_id"`
		Name           string `json:"name"`
		Private        bool   `json:"private"`
		PrivateBalance bool   `json:"private_balance"`
	} `json:"budget"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *BudgetResource) PostAuthRequired() bool {
	return true
}

// PostDoc returns the description of this API endpoint
func (r *BudgetResource) PostDoc() string {
	return "create a new budget"
}

// PostParams returns the parameters supported by this API endpoint
func (r *BudgetResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *BudgetResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	/*auth, err := context.Authentication(request)
		if err != nil || auth.(db.Budget).ID != 1 {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Admin permission required for this operation",
			"BudgetResource POST"))
		return
	}*/

	ups := data.(*BudgetPostStruct)

	budget := db.Budget{
		ProjectID:      &ups.Budget.ProjectID,
		ParentID:       ups.Budget.ParentID,
		Name:           ups.Budget.Name,
		Private:        ups.Budget.Private,
		PrivateBalance: ups.Budget.PrivateBalance,
	}
	err := budget.Save(context.(*db.APIContext))
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
