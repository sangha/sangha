package db

import (
	"math/rand"
	"strconv"
	"time"
)

// Budget represents the db schema of a budget
type Budget struct {
	ID             int64
	ProjectID      int64
	UserID         int64
	ParentID       int64
	Name           string
	Private        bool
	PrivateBalance bool
}

// LoadBudgetByID loads a budget by ID from the database
func (context *APIContext) LoadBudgetByID(id int64) (Budget, error) {
	budget := Budget{}
	if id < 1 {
		return budget, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, project_id, user_id, parent, name, private, private_balance FROM budgets WHERE id = $1", id).
		Scan(&budget.ID, &budget.ProjectID, &budget.UserID, &budget.ParentID, &budget.Name, &budget.Private, &budget.PrivateBalance)
	return budget, err
}

// LoadRootBudgetForProject loads the root budget for a project from the database
func (context *APIContext) LoadRootBudgetForProject(project *Project) (Budget, error) {
	budget := Budget{}
	if project == nil {
		return budget, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, project_id, user_id, parent, name, private, private_balance FROM budgets WHERE project_id = $1 AND parent = 0", project.ID).
		Scan(&budget.ID, &budget.ProjectID, &budget.UserID, &budget.ParentID, &budget.Name, &budget.Private, &budget.PrivateBalance)
	return budget, err
}

// GetBudgetByID returns a budget by ID from the cache
func (context *APIContext) GetBudgetByID(id int64) (Budget, error) {
	budget := Budget{}
	budgetsCache, err := budgetsCache.Value(id, context)
	if err != nil {
		return budget, err
	}

	budget = *budgetsCache.Data().(*Budget)
	return budget, nil
}

// LoadAllBudgets loads all budgets from the database
func (context *APIContext) LoadAllBudgets() ([]Budget, error) {
	budgets := []Budget{}

	rows, err := context.Query("SELECT id, project_id, user_id, parent, name, private, private_balance FROM budgets")
	if err != nil {
		return budgets, err
	}

	defer rows.Close()
	for rows.Next() {
		budget := Budget{}
		err = rows.Scan(&budget.ID, &budget.ProjectID, &budget.UserID, &budget.ParentID, &budget.Name, &budget.Private, &budget.PrivateBalance)
		if err != nil {
			return budgets, err
		}

		budgets = append(budgets, budget)
	}

	return budgets, err
}

// Update a budget in the database
func (budget *Budget) Update(context *APIContext) error {
	_, err := context.Exec("UPDATE budgets SET project_id = $1, user_id = $2, parent = $3, name = $4, private = $5, private_balance = $6 WHERE id = $7",
		budget.ProjectID, budget.UserID, budget.ParentID, budget.Name, budget.Private, budget.PrivateBalance, budget.ID)
	if err != nil {
		panic(err)
	}

	budgetsCache.Delete(budget.ID)
	return err
}

// Save a budget to the database
func (budget *Budget) Save(context *APIContext) error {
	err := context.QueryRow("INSERT INTO budgets (project_id, user_id, parent, name, private, private_balance) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		budget.ProjectID, budget.UserID, budget.ParentID, budget.Name, budget.Private, budget.PrivateBalance).Scan(&budget.ID)
	budgetsCache.Delete(budget.ID)
	return err
}

func (budget *Budget) Balance(context *APIContext) (float64, error) {
	var val float64
	err := context.QueryRow("SELECT SUM(value) FROM transactions WHERE budget_id = $1", budget.ID).
		Scan(&val)
	return val, err
}

type BudgetRatioPair struct {
	budget_ids []string
	ratios     []string
}
type BudgetSorter BudgetRatioPair

func (a BudgetSorter) Len() int {
	return len(a.budget_ids)
}

func (a BudgetSorter) Swap(i, j int) {
	a.budget_ids[i], a.budget_ids[j] = a.budget_ids[j], a.budget_ids[i]
	a.ratios[i], a.ratios[j] = a.ratios[j], a.ratios[i]
}

func (a BudgetSorter) Less(i, j int) bool {
	in, _ := strconv.Atoi(a.budget_ids[i])
	jn, _ := strconv.Atoi(a.budget_ids[j])
	return in < jn
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
