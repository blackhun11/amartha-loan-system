package model

import (
	"errors"
	"time"
)

type LoanState string

const (
	StateProposed  LoanState = "PROPOSED"
	StateApproved  LoanState = "APPROVED"
	StateInvested  LoanState = "INVESTED"
	StateDisbursed LoanState = "DISBURSED"
)

type Loan struct {
	ID            int64         `json:"id,omitempty"`
	BorrowerID    int64         `json:"borrower_id,omitempty"`
	Principal     float64       `json:"principal,omitempty"`
	Rate          float64       `json:"rate,omitempty"`
	ROI           float64       `json:"roi,omitempty"`
	State         LoanState     `json:"state,omitempty"`
	Approval      *Approval     `json:"approval,omitempty"`
	Investments   []Investment  `json:"investments,omitempty"`
	Disbursement  *Disbursement `json:"disbursement,omitempty"`
	AgreementLink string        `json:"agreement_link,omitempty"`
}

type Approval struct {
	ValidatorID int64     `json:"validator_id,omitempty"`
	ProofURL    string    `json:"proof_url,omitempty"`
	ApprovedAt  time.Time `json:"approved_at,omitempty"`
}

type Investment struct {
	InvestorID int64   `json:"investor_id,omitempty"`
	Amount     float64 `json:"amount,omitempty"`
}

type Disbursement struct {
	OfficerID    int64     `json:"officer_id,omitempty"`
	AgreementURL string    `json:"agreement_url,omitempty"`
	DisbursedAt  time.Time `json:"disbursed_at,omitempty"`
}

func (l *Loan) CanTransitionTo(newState LoanState) bool {
	currentOrder := stateOrder[l.State]
	newOrder := stateOrder[newState]
	return newOrder-currentOrder == 1
}

var stateOrder = map[LoanState]int{
	StateProposed:  1,
	StateApproved:  2,
	StateInvested:  3,
	StateDisbursed: 4,
}

func (l *Loan) Approve(approval Approval) error {
	if l.CanTransitionTo(StateApproved) {
		l.State = StateApproved
		l.Approval = &approval
	} else {
		return errors.New("can only approve when loan is proposed")
	}
	return nil
}

func (l *Loan) AddInvestment(investment Investment) error {
	total := investment.Amount
	for _, inv := range l.Investments {
		total += inv.Amount
	}

	if total > l.Principal {
		return errors.New("total investments exceed principal")
	}

	if l.CanTransitionTo(StateInvested) {
		l.Investments = append(l.Investments, investment)
		if total == l.Principal {
			l.State = StateInvested
		}
	} else {
		return errors.New("can only invest when loan is approved")
	}

	return nil
}

func (l *Loan) Disburse(disbursement Disbursement) error {
	if l.CanTransitionTo(StateDisbursed) {
		l.State = StateDisbursed
		l.Disbursement = &disbursement
	} else {
		return errors.New("can only disburse when loan is invested")
	}

	return nil
}
