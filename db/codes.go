package db

import (
	"errors"
	"sort"
	"strconv"

	"github.com/muesli/toktok"
)

// Code represents the db schema of a code
type Code struct {
	ID        int64
	Code      string
	BudgetIDs StringSlice
	Ratios    StringSlice
	UserID    *int64
}

var (
	// ErrInvalidBudgetRatioSet is the error returned when encountering an invalid set of budgets & ratios
	ErrInvalidBudgetRatioSet = errors.New("Budget & ratio sets have different sizes")
	// ErrInvalidRatio is the error returned when encountering an invalid ratio
	ErrInvalidRatio = errors.New("Ratios contain non-numerical characters or don't add up to 100%")
)

// LoadCodeByCode loads a code by Code from the database
func (context *APIContext) LoadCodeByCode(c string) (Code, error) {
	code := Code{}

	err := context.QueryRow("SELECT id, code, budget_ids, ratios, user_id FROM codes WHERE code = $1", c).
		Scan(&code.ID, &code.Code, &code.BudgetIDs, &code.Ratios, &code.UserID)
	return code, err
}

// LoadCodeByID loads a code by ID from the database
func (context *APIContext) LoadCodeByID(id int64) (Code, error) {
	code := Code{}
	if id < 1 {
		return code, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, code, budget_ids, ratios, user_id FROM codes WHERE id = $1", id).
		Scan(&code.ID, &code.Code, &code.BudgetIDs, &code.Ratios, &code.UserID)
	return code, err
}

// LoadCodeByBudgetsAndRatios loads a code by budgetIDs and their ratios from the database
func (context *APIContext) LoadCodeByBudgetsAndRatios(budgetIDs, ratios StringSlice, userID string) (Code, error) {
	code := Code{}
	if len(budgetIDs) != len(ratios) {
		return code, ErrInvalidBudgetRatioSet
	}

	// make sure proper ratios have been submitted
	totalRatio := 0
	for _, ratio := range ratios {
		r, err := strconv.Atoi(ratio)
		if err != nil {
			return code, ErrInvalidRatio
		}

		totalRatio += r
	}
	if totalRatio != 100 {
		return code, ErrInvalidRatio
	}

	var bids StringSlice
	for _, bid := range budgetIDs {
		budget, err := context.GetBudgetByUUID(bid)
		if err != nil {
			panic(err)
		}
		bids = append(bids, strconv.FormatInt(budget.ID, 10))
	}

	// user may be empty
	user, _ := context.GetUserByUUID(userID)

	// sort budgets & ratios
	sort.Sort(BudgetSorter(BudgetRatioPair{bids, ratios}))

	code = Code{
		BudgetIDs: bids,
		Ratios:    ratios,
		UserID:    &user.ID,
	}

	codes, err := context.LoadAllCodes()
	if err != nil {
		panic(err)
	}
	tokens := []string{}
	for _, code := range codes {
		tokens = append(tokens, code.Code)
	}

	// FIXME: we want to populate the bucket at startup
	bucket, _ := toktok.NewBucket(8)
	bucket.LoadTokens(tokens)
	code.Code, err = bucket.NewToken(8)
	if err != nil {
		return code, err
	}

	if user.ID > 0 {
		err = context.QueryRow("SELECT id, code, budget_ids, ratios, user_id FROM codes WHERE budget_ids = $1 AND ratios = $2 AND user_id = $3", bids, ratios, user.ID).
			Scan(&code.ID, &code.Code, &code.BudgetIDs, &code.Ratios, &code.UserID)
	} else {
		err = context.QueryRow("SELECT id, code, budget_ids, ratios FROM codes WHERE budget_ids = $1 AND ratios = $2 AND user_id IS NULL", bids, ratios).
			Scan(&code.ID, &code.Code, &code.BudgetIDs, &code.Ratios)
	}
	if err != nil {
		if user.ID > 0 {
			err = context.QueryRow("INSERT INTO codes (code, budget_ids, ratios, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
				code.Code, code.BudgetIDs, code.Ratios, code.UserID).Scan(&code.ID)
		} else {
			err = context.QueryRow("INSERT INTO codes (code, budget_ids, ratios, user_id) VALUES ($1, $2, $3, null) RETURNING id",
				code.Code, code.BudgetIDs, code.Ratios).Scan(&code.ID)
		}
		codesCache.Delete(code.ID)
	}

	return code, err
}

// GetCodeByID returns a code by ID from the cache
func (context *APIContext) GetCodeByID(id int64) (Code, error) {
	code := Code{}
	codesCache, err := codesCache.Value(id, context)
	if err != nil {
		return code, err
	}

	code = *codesCache.Data().(*Code)
	return code, nil
}

// LoadAllCodes loads all codes from the database
func (context *APIContext) LoadAllCodes() ([]Code, error) {
	codes := []Code{}

	rows, err := context.Query("SELECT id, code, budget_ids, ratios, user_id FROM codes")
	if err != nil {
		return codes, err
	}

	defer rows.Close()
	for rows.Next() {
		code := Code{}
		err = rows.Scan(&code.ID, &code.Code, &code.BudgetIDs, &code.Ratios, &code.UserID)
		if err != nil {
			return codes, err
		}

		codes = append(codes, code)
	}

	return codes, err
}

// Save a code to the database
/*
func (code *Code) Save(context *APIContext) error {
	err := context.QueryRow("INSERT INTO codes (code) VALUES ($1) RETURNING id",
		code.Code).Scan(&code.ID)
	codesCache.Delete(code.ID)
	return err
}
*/
