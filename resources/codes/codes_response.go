package codes

import (
	"strconv"

	"gitlab.techcultivation.org/sangha/sangha/db"
	"gitlab.techcultivation.org/sangha/sangha/resources/budgets"
	"gitlab.techcultivation.org/sangha/sangha/resources/projects"

	"github.com/muesli/smolder"
)

// CodeResponse is the common response to 'code' requests
type CodeResponse struct {
	smolder.Response

	Codes    []codeInfoResponse             `json:"codes,omitempty"`
	Budgets  []budgets.BudgetInfoResponse   `json:"budgets,omitempty"`
	Projects []projects.ProjectInfoResponse `json:"projects,omitempty"`

	codes    []db.Code
	budgets  []db.Budget
	projects []db.Project
}

type codeInfoResponse struct {
	ID      string   `json:"id"`
	Code    string   `json:"token"`
	Budgets []string `json:"budgets"`
	Ratios  []string `json:"ratios"`
}

// Init a new response
func (r *CodeResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Codes = []codeInfoResponse{}
	r.Budgets = []budgets.BudgetInfoResponse{}
}

// AddCode adds a code to the response
func (r *CodeResponse) AddCode(code *db.Code) {
	r.codes = append(r.codes, *code)
	r.Codes = append(r.Codes, prepareCodeResponse(r.Context, code))

	for _, b := range code.BudgetIDs {
		bid, _ := strconv.ParseInt(b, 10, 64)
		budget, err := r.Context.(*db.APIContext).LoadBudgetByID(bid)
		if err != nil {
			panic(err)
		}

		r.budgets = append(r.budgets, budget)
		r.Budgets = append(r.Budgets, budgets.PrepareBudgetResponse(r.Context, &budget))

		project, err := r.Context.(*db.APIContext).GetProjectByID(*budget.ProjectID)
		if err != nil {
			panic(err)
		}

		r.projects = append(r.projects, project)
		r.Projects = append(r.Projects, projects.PrepareProjectResponse(r.Context, &project))
	}
}

// AddBudget adds a budget to the response
func (r *CodeResponse) AddBudget(budget *db.Budget) {
	r.budgets = append(r.budgets, *budget)
	r.Budgets = append(r.Budgets, budgets.PrepareBudgetResponse(r.Context, budget))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *CodeResponse) EmptyResponse() interface{} {
	if len(r.codes) == 0 {
		var out struct {
			Codes interface{} `json:"codes"`
		}
		out.Codes = []codeInfoResponse{}
		return out
	}
	return nil
}

func prepareCodeResponse(context smolder.APIContext, code *db.Code) codeInfoResponse {
	ctx := context.(*db.APIContext)
	var budgets []string
	for _, b := range code.BudgetIDs {
		bid, _ := strconv.ParseInt(b, 10, 64)
		budget, err := ctx.LoadBudgetByID(bid)
		if err != nil {
			panic(err)
		}

		budgets = append(budgets, budget.UUID)
	}

	resp := codeInfoResponse{
		ID:      code.Code,
		Code:    code.Code,
		Budgets: budgets,
		Ratios:  code.Ratios,
	}

	return resp
}
