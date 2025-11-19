package request

import "time"

type CreateLoanRequest struct {
	Principal     float64 `json:"principal" validate:"required,gt=0"`
	BorrowerID    int64   `json:"borrower_id" validate:"required"`
	Rate          float64 `json:"rate" validate:"required,gt=0"`
	ROI           float64 `json:"roi" validate:"required,gt=0"`
	AgreementLink string  `json:"agreement_link" validate:"required"`
}

type ApproveLoanRequest struct {
	ID          int64     `param:"id" validate:"required"`
	ValidatorID int64     `json:"validator_id" validate:"required"`
	ProofURL    string    `json:"proof_url" validate:"required"`
	ApprovedAt  time.Time `json:"approved_at"`
}

type InvestLoanRequest struct {
	ID         int64   `param:"id" validate:"required"`
	InvestorID int64   `json:"investor_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

type DisburseLoanRequest struct {
	ID           int64     `param:"id" validate:"required"`
	OfficerID    int64     `json:"officer_id" validate:"required"`
	AgreementURL string    `json:"agreement_url" validate:"required"`
	DisbursedAt  time.Time `json:"disbursed_at"`
}

type GetLoanRequest struct {
	ID int64 `param:"id" validate:"required"`
}
