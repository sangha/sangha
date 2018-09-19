package db

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	money "github.com/muesli/go-money"
	"github.com/muesli/toktok"
	"gitlab.techcultivation.org/sangha/mq"
)

// Payment represents the db schema of a payment
type Payment struct {
	ID                  int64
	BudgetID            int64
	CreatedAt           time.Time
	Amount              int64
	Currency            string
	Code                string
	Purpose             string
	RemoteAccount       string
	RemoteName          string
	RemoteTransactionID string
	RemoteBankID        string
	Source              string
	Pending             bool
}

// LoadPaymentByID loads a payment by ID from the database
func (context *APIContext) LoadPaymentByID(id int64) (Payment, error) {
	payment := Payment{}
	if id < 1 {
		return payment, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, budget_id, created_at, amount, currency, code, purpose, remote_account, "+
		"remote_name, remote_transaction_id, remote_bank_id, source, pending "+
		"FROM payments "+
		"WHERE id = $1", id).
		Scan(&payment.ID, &payment.BudgetID, &payment.CreatedAt, &payment.Amount, &payment.Currency, &payment.Code,
			&payment.Purpose, &payment.RemoteAccount, &payment.RemoteName, &payment.RemoteTransactionID, &payment.RemoteBankID,
			&payment.Source, &payment.Pending)

	return payment, err
}

// LoadPayments loads all payments for a budget
func (budget *Budget) LoadPayments(context *APIContext) ([]Payment, error) {
	payments := []Payment{}

	rows, err := context.Query("SELECT id, budget_id, created_at, amount, currency, code, purpose, remote_account, "+
		"remote_name, remote_transaction_id, remote_bank_id, source, pending "+
		"FROM payments "+
		"WHERE budget_id = $1 "+
		"ORDER BY created_at ASC", budget.ID)
	if err != nil {
		return payments, err
	}

	defer rows.Close()
	for rows.Next() {
		payment := Payment{}
		err = rows.Scan(&payment.ID, &payment.BudgetID, &payment.CreatedAt, &payment.Amount, &payment.Currency, &payment.Code,
			&payment.Purpose, &payment.RemoteAccount, &payment.RemoteName, &payment.RemoteTransactionID, &payment.RemoteBankID,
			&payment.Source, &payment.Pending)

		if err != nil {
			return payments, err
		}

		payments = append(payments, payment)
	}

	return payments, err
}

// LoadPaymentsForDonor loads all payments for a specific donor
func (context *APIContext) LoadPaymentsForDonor(donor string) ([]Payment, error) {
	payments := []Payment{}

	rows, err := context.Query("SELECT id, budget_id, created_at, amount, currency, code, purpose, remote_account, "+
		"remote_name, remote_transaction_id, remote_bank_id, source, pending "+
		"FROM payments "+
		"WHERE remote_account = $1 "+
		"ORDER BY created_at ASC", donor)
	if err != nil {
		return payments, err
	}

	defer rows.Close()
	for rows.Next() {
		payment := Payment{}
		err = rows.Scan(&payment.ID, &payment.BudgetID, &payment.CreatedAt, &payment.Amount, &payment.Currency, &payment.Code,
			&payment.Purpose, &payment.RemoteAccount, &payment.RemoteName, &payment.RemoteTransactionID, &payment.RemoteBankID,
			&payment.Source, &payment.Pending)

		if err != nil {
			return payments, err
		}

		payments = append(payments, payment)
	}

	return payments, err
}

// LoadPendingPayment loads all pending payments
func (context *APIContext) LoadPendingPayments(direction int) ([]Payment, error) {
	payments := []Payment{}

	var filter string
	switch direction {
	case TRANSACTION_INCOMING:
		filter = "AND amount > 0"
	case TRANSACTION_OUTGOING:
		filter = "AND amount < 0"
	}

	rows, err := context.Query(fmt.Sprintf("SELECT id, budget_id, created_at, amount, currency, code, purpose, remote_account, "+
		"remote_name, remote_transaction_id, remote_bank_id, source, pending "+
		"FROM payments "+
		"WHERE pending = true %s "+
		"ORDER BY created_at ASC", filter))

	if err != nil {
		return payments, err
	}

	defer rows.Close()
	for rows.Next() {
		payment := Payment{}
		err = rows.Scan(&payment.ID, &payment.BudgetID, &payment.CreatedAt, &payment.Amount, &payment.Currency, &payment.Code,
			&payment.Purpose, &payment.RemoteAccount, &payment.RemoteName, &payment.RemoteTransactionID, &payment.RemoteBankID,
			&payment.Source, &payment.Pending)

		if err != nil {
			return payments, err
		}

		payments = append(payments, payment)
	}

	return payments, err
}

// Process turns a payment into various budget transactions
func (payment *Payment) Process(context *APIContext) error {
	code, err := context.LoadCodeByCode(payment.Code)
	if err != nil {
		return err
	}

	var ratios []int
	for _, r := range code.Ratios {
		ratio, _ := strconv.ParseInt(r, 10, 64)
		ratios = append(ratios, int(ratio))
	}

	var budgets []Budget
	var cuts []int
	for _, b := range code.BudgetIDs {
		bid, _ := strconv.ParseInt(b, 10, 64)
		budget, err := context.LoadBudgetByID(bid)
		if err != nil {
			return err
		}

		p, err := context.GetProjectByID(*budget.ProjectID)
		if err != nil {
			return err
		}

		if payment.Amount > 0 {
			cuts = append(cuts, int(p.ProcessingCut))
		} else {
			cuts = append(cuts, int(0))
		}
		budgets = append(budgets, budget)
	}

	cutBudget := int64(49)

	// transaction to cct account
	t := Transaction{
		BudgetID:  payment.BudgetID,
		Amount:    payment.Amount,
		CreatedAt: payment.CreatedAt, // FIXME: time.Now().UTC(),
		Purpose:   payment.Purpose,
		PaymentID: &payment.ID,
	}
	if err = t.Save(context); err != nil {
		return err
	}

	eur := money.New(payment.Amount, "EUR")
	parties, err := eur.Allocate(ratios...)
	if err != nil {
		return err
	}

	for idx, b := range budgets {
		peur := money.New(parties[idx].Amount(), "EUR")
		fees, err := peur.Allocate(cuts[idx], 100-cuts[idx])
		if err != nil {
			return err
		}

		if fees[1].Amount() != 0 && payment.BudgetID != b.ID {
			_, err = context.Transfer(payment.BudgetID, b.ID, fees[1].Amount(), payment.Purpose, payment.ID, payment.CreatedAt)
			if err != nil {
				return err
			}
		}

		if fees[0].Amount() != 0 {
			_, err = context.Transfer(payment.BudgetID, cutBudget, fees[0].Amount(), payment.Purpose, payment.ID, payment.CreatedAt)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Update a payment in the database
func (payment *Payment) Update(context *APIContext) error {
	_, err := context.LoadCodeByCode(payment.Code)
	if err != nil {
		return err
	}

	_, err = context.Exec("UPDATE payments SET code = $1, pending = $2 WHERE id = $3",
		payment.Code, payment.Pending, payment.ID)
	if err != nil {
		return err
	}

	if !payment.Pending {
		p := mq.Payment{
			Name:    payment.RemoteName,
			Address: []string{},

			DateTime: payment.CreatedAt,
			Amount:   payment.Amount,
			Currency: payment.Currency,

			TransactionCode: payment.Code,
			Description:     payment.Purpose,

			Source:              payment.Source,
			SourceID:            payment.RemoteBankID,
			SourcePayerID:       payment.RemoteAccount,
			SourceTransactionID: payment.RemoteTransactionID,

			BudgetID:  payment.BudgetID,
			PaymentID: payment.ID,
		}
		err = p.Process()
	}

	return err
}

// Save a payment to the database
func (payment *Payment) Save(context *APIContext) error {
	if payment.Code == "" {
		r := regexp.MustCompile("[^\\W]+")
		p := r.FindAllString(payment.Purpose, -1)

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

		ldist := 32768
		var lmatch string
		for _, c := range p {
			fmt.Println("Resolving", c)
			match, dist := bucket.Resolve(c)
			if dist < ldist && dist < 4 {
				ldist = dist
				lmatch = match
			}
			fmt.Println("Match", match, dist)
		}
		fmt.Println("Lowest Match", lmatch, ldist)
		if ldist < 2 {
			payment.Code = lmatch
		}
	}

	err := context.QueryRow("INSERT INTO payments (budget_id, created_at, amount, currency, code, purpose, remote_account, "+
		"remote_name, remote_transaction_id, remote_bank_id, source) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) "+
		"RETURNING id",
		payment.BudgetID, payment.CreatedAt, payment.Amount, payment.Currency, payment.Code, payment.Purpose, payment.RemoteAccount,
		payment.RemoteName, payment.RemoteTransactionID, payment.RemoteBankID, payment.Source).Scan(&payment.ID)
	return err
}

func (context *APIContext) LatestPaymentFromSource(source string) (Payment, error) {
	payment := Payment{}

	err := context.QueryRow("SELECT id, budget_id, created_at, amount, currency, code, purpose, remote_account, "+
		"remote_name, remote_transaction_id, remote_bank_id, source, pending "+
		"FROM payments "+
		"WHERE source = $1 "+
		"ORDER BY created_at DESC LIMIT 1", source).
		Scan(&payment.ID, &payment.BudgetID, &payment.CreatedAt, &payment.Amount, &payment.Currency, &payment.Code,
			&payment.Purpose, &payment.RemoteAccount, &payment.RemoteName, &payment.RemoteTransactionID, &payment.RemoteBankID,
			&payment.Source, &payment.Pending)

	return payment, err
}

// SearchPayments searches database for payments
func (context *APIContext) SearchPayments(term string) ([]Payment, error) {
	payments := []Payment{}

	rows, err := context.Query("SELECT DISTINCT id FROM payments WHERE "+
		"(LOWER(purpose) LIKE LOWER('%' || $1 || '%') OR "+
		"LOWER(remote_account) LIKE LOWER('%' || $1 || '%') OR "+
		"LOWER(remote_name) LIKE LOWER('%' || $1 || '%'))", term)
	if err != nil {
		return payments, err
	}

	defer rows.Close()
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			return payments, err
		}

		p, err := context.LoadPaymentByID(id)
		if err != nil {
			return payments, err
		}

		payments = append(payments, p)
	}

	return payments, nil
}
