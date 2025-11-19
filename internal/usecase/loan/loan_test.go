package loan_test

import (
	"context"
	"errors"
	"loan_system/internal/model"
	loanrepo "loan_system/internal/repository/loan/mock"
	pubsubrepo "loan_system/internal/repository/pubsub"
	"loan_system/internal/usecase/loan"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoanUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoMock := loanrepo.NewMockRepository(ctrl)
	pubsubMock := pubsubrepo.NewMock()
	uc := loan.NewUsecase(repoMock, pubsubMock)

	t.Run("FindAll", func(t *testing.T) {
		repoMock.EXPECT().FindAll(gomock.Any()).Return([]*model.Loan{{ID: 1}, {ID: 2}}, nil)

		loans, err := uc.FindAll(context.Background())
		assert.NoError(t, err)
		assert.Len(t, loans, 2)
	})

	t.Run("FindByID", func(t *testing.T) {
		repoMock.EXPECT().FindByID(gomock.Any(), int64(1)).Return(&model.Loan{ID: 1}, nil)

		loan, err := uc.FindByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), loan.ID)
	})

	t.Run("CreateLoan", func(t *testing.T) {
		repoMock.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
		err := uc.CreateLoan(context.Background(), &model.Loan{})
		assert.NoError(t, err)
	})

	t.Run("ApproveLoan Success", func(t *testing.T) {
		mockLoan := &model.Loan{ID: 1, State: model.StateProposed}
		repoMock.EXPECT().FindByID(gomock.Any(), int64(1)).Return(mockLoan, nil)
		repoMock.EXPECT().Update(gomock.Any(), mockLoan).Return(nil)

		_, err := uc.ApproveLoan(context.Background(), 1, model.Approval{})
		assert.NoError(t, err)
	})

	t.Run("ApproveLoan InvalidState", func(t *testing.T) {
		mockLoan := &model.Loan{ID: 1, State: model.StateApproved}
		repoMock.EXPECT().FindByID(gomock.Any(), int64(1)).Return(mockLoan, nil)

		_, err := uc.ApproveLoan(context.Background(), 1, model.Approval{})
		assert.ErrorContains(t, err, "can only approve")
	})

	t.Run("ApproveLoan loan not exist", func(t *testing.T) {
		repoMock.EXPECT().FindByID(gomock.Any(), int64(1)).Return(nil, errors.New("loan not found"))

		_, err := uc.ApproveLoan(context.Background(), 1, model.Approval{})
		assert.ErrorContains(t, err, "loan not found")
	})

	t.Run("AddInvestment FullFunding", func(t *testing.T) {
		loan := &model.Loan{
			ID:          2,
			Principal:   1000,
			Investments: []model.Investment{{Amount: 900}},
			State:       model.StateApproved,
		}
		repoMock.EXPECT().FindByID(gomock.Any(), int64(2)).Return(loan, nil)

		repoMock.EXPECT().Update(gomock.Any(), loan).Return(nil)

		_, err := uc.AddInvestment(context.Background(), 2, model.Investment{Amount: 100})
		assert.Equal(t, loan.State, model.StateInvested)
		assert.NoError(t, err)
	})

	t.Run("AddInvestment InvalidState", func(t *testing.T) {
		loan := &model.Loan{ID: 2, State: model.StateInvested}
		repoMock.EXPECT().FindByID(gomock.Any(), int64(2)).Return(loan, nil)

		_, err := uc.AddInvestment(context.Background(), 2, model.Investment{Amount: 100})
		assert.ErrorContains(t, err, "investments exceed principal")
	})

	t.Run("AddInvestment loan not found", func(t *testing.T) {
		repoMock.EXPECT().FindByID(gomock.Any(), int64(2)).Return(nil, errors.New("loan not found"))

		_, err := uc.AddInvestment(context.Background(), 2, model.Investment{Amount: 100})
		assert.ErrorContains(t, err, "loan not found")
	})

	t.Run("DisburseLoan Success", func(t *testing.T) {
		loan := &model.Loan{ID: 3, State: model.StateInvested}
		repoMock.EXPECT().FindByID(gomock.Any(), int64(3)).Return(loan, nil)
		repoMock.EXPECT().Update(gomock.Any(), loan).Return(nil)

		_, err := uc.DisburseLoan(context.Background(), 3, model.Disbursement{})
		assert.NoError(t, err)
	})

	t.Run("DisburseLoan InvalidState", func(t *testing.T) {
		repoMock.EXPECT().FindByID(gomock.Any(), int64(3)).Return(&model.Loan{State: model.StateApproved}, nil)
		_, err := uc.DisburseLoan(context.Background(), 3, model.Disbursement{})
		assert.ErrorContains(t, err, "can only disburse")
	})

	t.Run("DisburseLoan loan not found", func(t *testing.T) {
		repoMock.EXPECT().FindByID(gomock.Any(), int64(3)).Return(nil, errors.New("loan not found"))

		_, err := uc.DisburseLoan(context.Background(), 3, model.Disbursement{})
		assert.ErrorContains(t, err, "loan not found")
	})
}
