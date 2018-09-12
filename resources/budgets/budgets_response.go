package budgets

import (
	"gitlab.techcultivation.org/sangha/sangha/db"
	"gitlab.techcultivation.org/sangha/sangha/resources/projects"

	"github.com/muesli/smolder"
)

// BudgetResponse is the common response to 'budget' requests
type BudgetResponse struct {
	smolder.Response

	Budgets []BudgetInfoResponse `json:"budgets,omitempty"`
	budgets []db.Budget

	Projects []projects.ProjectInfoResponse `json:"projects,omitempty"`
	projects []db.Project
}

type BudgetInfoResponse struct {
	ID          string `json:"id"`
	Project     string `json:"project"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Balance     int64  `json:"balance"`
	Code        string `json:"code"`
}

// Init a new response
func (r *BudgetResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Budgets = []BudgetInfoResponse{}
	r.Projects = []projects.ProjectInfoResponse{}
}

// AddBudget adds a budget to the response
func (r *BudgetResponse) AddBudget(budget *db.Budget) {
	r.budgets = append(r.budgets, *budget)
	r.Budgets = append(r.Budgets, PrepareBudgetResponse(r.Context, budget))

	project, err := r.Context.(*db.APIContext).GetProjectByID(*budget.ProjectID)
	if err != nil {
		panic(err)
	}

	r.projects = append(r.projects, project)
	r.Projects = append(r.Projects, projects.PrepareProjectResponse(r.Context, &project))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *BudgetResponse) EmptyResponse() interface{} {
	if len(r.budgets) == 0 {
		var out struct {
			Budgets interface{} `json:"budgets"`
		}
		out.Budgets = []BudgetInfoResponse{}
		return out
	}
	return nil
}

func PrepareBudgetResponse(context smolder.APIContext, budget *db.Budget) BudgetInfoResponse {
	ctx := context.(*db.APIContext)
	project, err := ctx.GetProjectByID(*budget.ProjectID)
	if err != nil {
		panic(err)
	}

	resp := BudgetInfoResponse{
		ID:          budget.UUID,
		Project:     project.UUID,
		Name:        budget.Name,
		Description: budget.Description,
	}

	resp.Balance, _ = budget.Balance(ctx)
	code, err := ctx.LoadCodeByBudgetUUID(budget.UUID)
	if err == nil {
		resp.Code = code.Code
	}

	return resp
}
