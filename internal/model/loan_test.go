package model_test

import (
	"loan_system/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoan_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		currentState model.LoanState
		newState     model.LoanState
		want         bool
	}{
		// TODO: Add test cases.
		{
			name:         "valid transition",
			currentState: model.StateProposed,
			newState:     model.StateApproved,
			want:         true,
		},
		{
			name:         "valid transition",
			currentState: model.StateApproved,
			newState:     model.StateInvested,
			want:         true,
		},
		{
			name:         "valid transition",
			currentState: model.StateInvested,
			newState:     model.StateDisbursed,
			want:         true,
		},
		{
			name:         "invalid transition",
			currentState: model.StateDisbursed,
			newState:     model.StateProposed,
			want:         false,
		},
		{
			name:         "invalid transition",
			currentState: model.StateDisbursed,
			newState:     model.StateApproved,
			want:         false,
		},
		{
			name:         "invalid transition",
			currentState: model.StateDisbursed,
			newState:     model.StateInvested,
			want:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: construct the receiver type.
			l := model.Loan{
				State: tt.currentState,
			}
			got := l.CanTransitionTo(tt.newState)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoanApprove(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		initialLoan   *model.Loan
		approval      model.Approval
		expectedError string
	}{
		{
			name:        "happy path approval",
			initialLoan: &model.Loan{State: model.StateProposed},
			approval:    model.Approval{ApprovedAt: now, ValidatorID: 123, ProofURL: "https://proof.com"},
		},
		{
			name:          "already approved",
			initialLoan:   &model.Loan{State: model.StateApproved},
			approval:      model.Approval{},
			expectedError: "can only approve when loan is proposed",
		},
		{
			name:          "invalid invested state",
			initialLoan:   &model.Loan{State: model.StateInvested},
			approval:      model.Approval{},
			expectedError: "can only approve when loan is proposed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.initialLoan.Approve(tt.approval)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, model.StateApproved, tt.initialLoan.State)

		})
	}
}

func TestLoanAddInvestment(t *testing.T) {
	principal := 5000.0
	tests := []struct {
		name          string
		initialLoan   *model.Loan
		investment    model.Investment
		expectedState model.LoanState
		expectedError string
	}{
		{
			name: "valid partial investment",
			initialLoan: &model.Loan{
				State:       model.StateApproved,
				Principal:   principal,
				Investments: []model.Investment{{Amount: 2000}},
			},
			investment:    model.Investment{Amount: 1500},
			expectedState: model.StateApproved,
		},
		{
			name: "full funding transition",
			initialLoan: &model.Loan{
				State:     model.StateApproved,
				Principal: principal,
			},
			investment:    model.Investment{Amount: 5000},
			expectedState: model.StateInvested,
		},
		{
			name: "exceed principal",
			initialLoan: &model.Loan{
				State:       model.StateApproved,
				Principal:   principal,
				Investments: []model.Investment{{Amount: 3000}},
			},
			investment:    model.Investment{Amount: 2500},
			expectedError: "total investments exceed principal",
		},
		{
			name: "invalid state investment",
			initialLoan: &model.Loan{
				State:     model.StateProposed,
				Principal: principal,
			},
			investment:    model.Investment{Amount: 1000},
			expectedError: "can only invest when loan is approved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.initialLoan.AddInvestment(tt.investment)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, tt.initialLoan.Investments, tt.investment)
			assert.Equal(t, tt.expectedState, tt.initialLoan.State)
		})
	}
}

func TestLoan_Disburse(t *testing.T) {
	tests := []struct {
		name          string
		initialState  model.LoanState
		wantErr       bool
		expectedState model.LoanState
	}{
		{
			name:          "valid from invested",
			initialState:  model.StateInvested,
			wantErr:       false,
			expectedState: model.StateDisbursed,
		},
		{
			name:         "invalid from approved",
			initialState: model.StateApproved,
			wantErr:      true,
		},
		{
			name:         "invalid from proposed",
			initialState: model.StateProposed,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &model.Loan{
				State:       tt.initialState,
				Investments: []model.Investment{{Amount: 1000}}, // Simulate fully invested
				Principal:   1000,
			}

			disbursement := model.Disbursement{
				OfficerID:    123,
				AgreementURL: "https://agreement.com",
				DisbursedAt:  time.Now(),
			}
			err := l.Disburse(disbursement)

			if tt.wantErr {
				assert.Error(t, err, "Expected error but got none")
				assert.Equal(t, tt.initialState, l.State, "State should not change on error")
			} else {
				assert.NoError(t, err, "Unexpected error")
				assert.Equal(t, tt.expectedState, l.State, "Incorrect final state")
				assert.Equal(t, &disbursement, l.Disbursement, "Disbursement data not set")
			}
		})
	}
}
