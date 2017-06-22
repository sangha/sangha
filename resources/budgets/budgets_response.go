package budgets

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/muesli/smolder"
)

// BudgetResponse is the common response to 'budget' requests
type BudgetResponse struct {
	smolder.Response

	Budgets []budgetInfoResponse `json:"budgets,omitempty"`
	budgets []db.Budget
}

type budgetInfoResponse struct {
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	Name      string `json:"name"`
}

// Init a new response
func (r *BudgetResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Budgets = []budgetInfoResponse{}
}

// AddBudget adds a budget to the response
func (r *BudgetResponse) AddBudget(budget *db.Budget) {
	r.budgets = append(r.budgets, *budget)
	r.Budgets = append(r.Budgets, prepareBudgetResponse(r.Context, budget))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *BudgetResponse) EmptyResponse() interface{} {
	if len(r.budgets) == 0 {
		var out struct {
			Budgets interface{} `json:"budgets"`
		}
		out.Budgets = []budgetInfoResponse{}
		return out
	}
	return nil
}

func prepareBudgetResponse(context smolder.APIContext, budget *db.Budget) budgetInfoResponse {
	resp := budgetInfoResponse{
		ID:        budget.ID,
		ProjectID: budget.ProjectID,
		Name:      budget.Name,
	}

	return resp
}
