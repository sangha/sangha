package db

// Budget represents the db schema of a budget
type Budget struct {
	ID        int64
	ProjectID int64
	Name      string
}

// LoadBudgetByID loads a budget by ID from the database
func (context *APIContext) LoadBudgetByID(id int64) (Budget, error) {
	budget := Budget{}
	if id < 1 {
		return budget, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, project_id, name FROM budgets WHERE id = $1", id).
		Scan(&budget.ID, &budget.ProjectID, &budget.Name)
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

	rows, err := context.Query("SELECT id, project_id, name FROM budgets")
	if err != nil {
		return budgets, err
	}

	defer rows.Close()
	for rows.Next() {
		budget := Budget{}
		err = rows.Scan(&budget.ID, &budget.ProjectID, &budget.Name)
		if err != nil {
			return budgets, err
		}

		budgets = append(budgets, budget)
	}

	return budgets, err
}

// Update a budget in the database
func (budget *Budget) Update(context *APIContext) error {
	_, err := context.Exec("UPDATE budgets SET project_id = $1, name = $2 WHERE id = $3",
		budget.ProjectID, budget.Name, budget.ID)
	if err != nil {
		panic(err)
	}

	budgetsCache.Delete(budget.ID)
	return err
}

// Save a budget to the database
func (budget *Budget) Save(context *APIContext) error {
	err := context.QueryRow("INSERT INTO budgets (project_id, name) VALUES ($1, $2) RETURNING id",
		budget.ProjectID, budget.Name).Scan(&budget.ID)
	budgetsCache.Delete(budget.ID)
	return err
}
