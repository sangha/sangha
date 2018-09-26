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
		Project        string `json:"project"`
		ParentID       int64  `json:"parent_id"`
		Name           string `json:"name"`
		Description    string `json:"description"`
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
	auth, err := context.Authentication(request)
	if err != nil || auth.(db.Budget).ID != 1 {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			"Admin permission required for this operation",
			"BudgetResource POST"))
		return
	}

	ups := data.(*BudgetPostStruct)

	project, err := context.(*db.APIContext).LoadProjectByUUID(ups.Budget.Project)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusBadRequest,
			"No such project",
			"BudgetResource POST"))
		return
	}

	budget := db.Budget{
		ProjectID:      &project.ID,
		ParentID:       ups.Budget.ParentID,
		Name:           ups.Budget.Name,
		Description:    ups.Budget.Description,
		Private:        ups.Budget.Private,
		PrivateBalance: ups.Budget.PrivateBalance,
	}
	err = budget.Save(context.(*db.APIContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			"Can't create budget",
			"BudgetResource POST"))
		return
	}

	resp := BudgetResponse{}
	resp.Init(context)
	resp.AddBudget(&budget)
	resp.Send(response)
}
