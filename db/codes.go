package db

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"

	"github.com/xrash/smetrics"
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

// LoadCodeByID loads a code by ID from the database
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

func (context *APIContext) LoadCodeByBudgetsAndRatios(budgetIDs, ratios StringSlice, userID int64) (Code, error) {
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

	// sort budgets & ratios
	sort.Sort(BudgetSorter(BudgetRatioPair{budgetIDs, ratios}))
	fmt.Println(budgetIDs)
	fmt.Println(ratios)

	code = Code{
		BudgetIDs: budgetIDs,
		Ratios:    ratios,
		UserID:    &userID,
	}

	codes, err := context.LoadAllCodes()
	if err != nil {
		panic(err)
	}
	for {
		code.Code = GenerateToken(8)

		dupe := false
		for _, c := range codes {
			if hd := smetrics.WagnerFischer(c.Code, code.Code, 1, 1, 2); hd <= 8 {
				fmt.Printf("%s is too similar to existing code %s (distance %d)\n", code.Code, c.Code, hd)
				dupe = true
				break
			}
		}
		if !dupe {
			break
		}
	}

	if userID > 0 {
		err = context.QueryRow("SELECT id, code, budget_ids, ratios, user_id FROM codes WHERE budget_ids = $1 AND ratios = $2 AND user_id = $3", budgetIDs, ratios, userID).
			Scan(&code.ID, &code.Code, &code.BudgetIDs, &code.Ratios, &code.UserID)
	} else {
		err = context.QueryRow("SELECT id, code, budget_ids, ratios FROM codes WHERE budget_ids = $1 AND ratios = $2 AND user_id IS NULL", budgetIDs, ratios).
			Scan(&code.ID, &code.Code, &code.BudgetIDs, &code.Ratios)
	}
	if err != nil {
		if userID > 0 {
			err = context.QueryRow("INSERT INTO codes (code, budget_ids, ratios, user_id) VALUES ($1, $2, $3, $4) RETURNING id",
				code.Code, budgetIDs, ratios, userID).Scan(&code.ID)
		} else {
			err = context.QueryRow("INSERT INTO codes (code, budget_ids, ratios, user_id) VALUES ($1, $2, $3, null) RETURNING id",
				code.Code, budgetIDs, ratios).Scan(&code.ID)
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

func GenerateToken(n int) string {
	var letterRunes = []rune("ABCDEFGHKLMNRSTWX3459")
	b := make([]rune, n)
	for i := range b {
		var lastrune rune
		if i > 0 {
			lastrune = b[i-1]
		}
		b[i] = lastrune
		for lastrune == b[i] {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
	}

	return string(b)
}
