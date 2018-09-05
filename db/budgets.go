package db

import (
	"math/rand"
	"strconv"
	"time"
)

// Budget represents the db schema of a budget
type Budget struct {
	ID             int64
	UUID           string
	ProjectID      *int64
	UserID         *int64
	ParentID       int64
	Name           string
	Description    string
	Private        bool
	PrivateBalance bool
}

// LoadBudgetByID loads a budget by UUID from the database
func (context *APIContext) LoadBudgetByID(id int64) (Budget, error) {
	budget := Budget{}

	err := context.QueryRow("SELECT id, uuid, project_id, user_id, parent, name, description, private, private_balance FROM budgets WHERE id = $1", id).
		Scan(&budget.ID, &budget.UUID, &budget.ProjectID, &budget.UserID, &budget.ParentID, &budget.Name, &budget.Description, &budget.Private, &budget.PrivateBalance)
	return budget, err
}

// LoadBudgetByUUID loads a budget by UUID from the database
func (context *APIContext) LoadBudgetByUUID(uuid string) (Budget, error) {
	budget := Budget{}
	if len(uuid) == 0 {
		return budget, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, uuid, project_id, user_id, parent, name, description, private, private_balance FROM budgets WHERE uuid = $1", uuid).
		Scan(&budget.ID, &budget.UUID, &budget.ProjectID, &budget.UserID, &budget.ParentID, &budget.Name, &budget.Description, &budget.Private, &budget.PrivateBalance)
	return budget, err
}

// LoadRootBudgetForProject loads the root budget for a project from the database
func (context *APIContext) LoadRootBudgetForProject(project *Project) (Budget, error) {
	budget := Budget{}
	if project == nil {
		return budget, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, uuid, project_id, user_id, parent, name, description, private, private_balance FROM budgets WHERE project_id = $1 AND parent = 0 ORDER BY id ASC", project.ID).
		Scan(&budget.ID, &budget.UUID, &budget.ProjectID, &budget.UserID, &budget.ParentID, &budget.Name, &budget.Description, &budget.Private, &budget.PrivateBalance)
	return budget, err
}

// LoadBudgets loads all budgets for a project
func (context *APIContext) LoadBudgets(project *Project) ([]Budget, error) {
	budgets := []Budget{}

	rows, err := context.Query("SELECT id, uuid, project_id, user_id, parent, name, description, private, private_balance FROM budgets WHERE project_id = $1 AND parent = 0 ORDER BY id ASC", project.ID)
	if err != nil {
		return budgets, err
	}

	defer rows.Close()
	for rows.Next() {
		budget := Budget{}
		err = rows.Scan(&budget.ID, &budget.UUID, &budget.ProjectID, &budget.UserID, &budget.ParentID, &budget.Name, &budget.Description, &budget.Private, &budget.PrivateBalance)
		if err != nil {
			return budgets, err
		}

		budgets = append(budgets, budget)
	}

	return budgets, err
}

// GetBudgetByUUID returns a budget by UUID from the cache
func (context *APIContext) GetBudgetByUUID(uuid string) (Budget, error) {
	budget := Budget{}
	budgetsCache, err := budgetsCache.Value(uuid, context)
	if err != nil {
		return budget, err
	}

	budget = *budgetsCache.Data().(*Budget)
	return budget, nil
}

// LoadAllBudgets loads all budgets from the database
func (context *APIContext) LoadAllBudgets() ([]Budget, error) {
	budgets := []Budget{}

	rows, err := context.Query("SELECT id, uuid, project_id, user_id, parent, name, description, private, private_balance FROM budgets")
	if err != nil {
		return budgets, err
	}

	defer rows.Close()
	for rows.Next() {
		budget := Budget{}
		err = rows.Scan(&budget.ID, &budget.UUID, &budget.ProjectID, &budget.UserID, &budget.ParentID, &budget.Name, &budget.Description, &budget.Private, &budget.PrivateBalance)
		if err != nil {
			return budgets, err
		}

		budgets = append(budgets, budget)
	}

	return budgets, err
}

// Update a budget in the database
func (budget *Budget) Update(context *APIContext) error {
	_, err := context.Exec("UPDATE budgets SET project_id = $1, user_id = $2, parent = $3, name = $4, description = $5, private = $6, private_balance = $7 WHERE id = $8",
		budget.ProjectID, budget.UserID, budget.ParentID, budget.Name, budget.Description, budget.Private, budget.PrivateBalance, budget.ID)
	budgetsCache.Delete(budget.UUID)
	return err
}

// Save a budget to the database
func (budget *Budget) Save(context *APIContext) error {
	budget.UUID, _ = UUID()

	err := context.QueryRow("INSERT INTO budgets (uuid, project_id, user_id, parent, name, description, private, private_balance) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		budget.UUID, budget.ProjectID, budget.UserID, budget.ParentID, budget.Name, budget.Description, budget.Private, budget.PrivateBalance).Scan(&budget.ID)
	budgetsCache.Delete(budget.UUID)
	return err
}

// Delete a budget from the database
func (budget *Budget) Delete(context *APIContext) error {
	_, err := context.Exec("DELETE FROM budgets WHERE id = $1", budget.ID)
	budgetsCache.Delete(budget.UUID)
	return err
}

// Balance returns this budget's total balance
func (budget *Budget) Balance(context *APIContext) (int64, error) {
	var val int64
	err := context.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE budget_id = $1", budget.ID).
		Scan(&val)
	return val, err
}

// BalanceStats returns this budget's total balance for the past months
func (budget *Budget) BalanceStats(context *APIContext) ([]int64, error) {
	var val []int64

	lom := time.Now().UTC()

	b := int64(-1)
	for b != 0 {
		err := context.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE budget_id = $1 AND created_at <= $2", budget.ID, lom).
			Scan(&b)
		if err != nil {
			return val, err
		}

		val = append(val, b)
		fom := time.Date(lom.Year(), lom.Month(), 1, 0, 0, 0, 0, time.UTC)
		lom = fom.AddDate(0, 0, 0).Add(time.Nanosecond * -1)
	}

	return val, nil
}

// SearchBudgets searches database for budgets
func (context *APIContext) SearchBudgets(term string) ([]Budget, error) {
	budgets := []Budget{}

	rows, err := context.Query("SELECT DISTINCT budgets.id FROM budgets, projects "+
		"WHERE projects.id = budgets.project_id AND "+
		"(LOWER(budgets.name) LIKE LOWER('%' || $1 || '%') OR "+
		"LOWER(projects.name) LIKE LOWER('%' || $1 || '%') OR "+
		"LOWER(budgets.description) LIKE LOWER('%' || $1 || '%'))", term)
	if err != nil {
		return budgets, err
	}

	defer rows.Close()
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return budgets, err
		}

		p, err := context.LoadBudgetByID(id)
		if err != nil {
			return budgets, err
		}

		budgets = append(budgets, p)
	}

	return budgets, nil
}

// BudgetRatioPair represents a pair of budgets and ratios
type BudgetRatioPair struct {
	budgetIDs []string
	ratios    []string
}

// BudgetSorter is used to sort the pair by ratio
type BudgetSorter BudgetRatioPair

func (a BudgetSorter) Len() int {
	return len(a.budgetIDs)
}

func (a BudgetSorter) Swap(i, j int) {
	a.budgetIDs[i], a.budgetIDs[j] = a.budgetIDs[j], a.budgetIDs[i]
	a.ratios[i], a.ratios[j] = a.ratios[j], a.ratios[i]
}

func (a BudgetSorter) Less(i, j int) bool {
	in, _ := strconv.Atoi(a.budgetIDs[i])
	jn, _ := strconv.Atoi(a.budgetIDs[j])
	return in < jn
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
