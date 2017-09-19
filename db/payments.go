package db

import "time"

// Payment represents the db schema of a payment
type Payment struct {
	ID                  int64
	UserID              int64
	Amount              float64
	Currency            string
	Code                string
	Description         string
	Source              string
	SourceID            string
	SourcePayerID       string
	SourceTransactionID string
	CreatedAt           time.Time
}

var ()

// LoadPaymentByID loads a payment by ID from the database
func (context *APIContext) LoadPaymentByID(id int64) (Payment, error) {
	payment := Payment{}
	if id < 1 {
		return payment, ErrInvalidID
	}

	err := context.QueryRow("SELECT id, user_id, amount, currency, code, description, source, source_id, source_payer_id, source_transaction_id, created_at FROM payments WHERE id = $1", id).
		Scan(&payment.ID, &payment.UserID, &payment.Amount, &payment.Currency, &payment.Code, &payment.Description, &payment.Source, &payment.SourceID, &payment.SourcePayerID, &payment.SourceTransactionID, &payment.CreatedAt)
	return payment, err
}

// Save a payment to the database
func (payment *Payment) Save(context *APIContext) error {
	err := context.QueryRow("INSERT INTO payments (user_id, amount, currency, code, description, source, source_id, source_payer_id, source_transaction_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		payment.UserID, payment.Amount, payment.Currency, payment.Code, payment.Description, payment.Source, payment.SourceID, payment.SourcePayerID, payment.SourceTransactionID, payment.CreatedAt).Scan(&payment.ID)
	return err
}
