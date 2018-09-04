package db

// Search represents the db schema of a search
type Search struct {
	ID       string
	Projects []Project
	Budgets  []Budget
	Payments []Payment
}

// Search searches the database for projects, budgets & payments
func (context *APIContext) Search(term string) (Search, error) {
	search := Search{
		ID: term,
	}

	var err error
	search.Projects, err = context.SearchProjects(term)
	if err != nil {
		return search, err
	}
	search.Budgets, err = context.SearchBudgets(term)
	if err != nil {
		return search, err
	}
	search.Payments, err = context.SearchPayments(term)
	if err != nil {
		return search, err
	}

	return search, nil
}
