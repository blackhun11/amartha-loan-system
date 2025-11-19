package loan

import (
	"context"
	"encoding/json"
	"fmt"

	"loan_system/internal/model"
	"loan_system/internal/repository/loan"
	"loan_system/internal/repository/pubsub"
)

//go:generate mockgen -source=loan.go -destination=mock/loan_mock.go -package=mock
type Usecase interface {
	CreateLoan(ctx context.Context, loan *model.Loan) error
	ApproveLoan(ctx context.Context, loanID int64, approval model.Approval) (loan *model.Loan, err error)
	AddInvestment(ctx context.Context, loanID int64, investment model.Investment) (loan *model.Loan, err error)
	DisburseLoan(ctx context.Context, loanID int64, disbursement model.Disbursement) (loan *model.Loan, err error)
}

type usecase struct {
	repo   loan.Repository
	pubsub pubsub.Mock
}

func NewUsecase(repo loan.Repository, pubsub pubsub.Mock) Usecase {
	return &usecase{repo: repo, pubsub: pubsub}
}

func (uc *usecase) CreateLoan(ctx context.Context, loan *model.Loan) error {
	loan.State = model.StateProposed

	return uc.repo.Save(ctx, loan)
}

func (uc *usecase) ApproveLoan(ctx context.Context, loanID int64, approval model.Approval) (loan *model.Loan, err error) {
	loan, err = uc.repo.FindByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	if err := loan.Approve(approval); err != nil {
		return nil, fmt.Errorf("approval failed: %w", err)
	}

	return loan, uc.repo.Update(ctx, loan)
}

func (uc *usecase) AddInvestment(ctx context.Context, loanID int64, investment model.Investment) (loan *model.Loan, err error) {
	loan, err = uc.repo.FindByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	if err := loan.AddInvestment(investment); err != nil {
		return nil, fmt.Errorf("investment failed: %w", err)
	}

	if loan.State == model.StateInvested {
		agreement := model.LoanAgreement{
			LoanID: loan.ID,
		}

		jsonData, err := json.Marshal(&agreement)
		if err != nil {
			return nil, fmt.Errorf("marshal agreement failed: %w", err)
		}

		// send trigger to pubsub to notify investor regarding agreement link

		if err := uc.pubsub.Publish(ctx, "loan_invested", jsonData); err != nil {
			return nil, fmt.Errorf("publish loan invested failed: %w", err)
		}

	}

	return loan, uc.repo.Update(ctx, loan)
}

func (uc *usecase) DisburseLoan(ctx context.Context, loanID int64, disbursement model.Disbursement) (loan *model.Loan, err error) {
	loan, err = uc.repo.FindByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	if err := loan.Disburse(disbursement); err != nil {
		return nil, fmt.Errorf("disburse failed: %w", err)
	}

	return loan, uc.repo.Update(ctx, loan)
}
