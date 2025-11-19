package http_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	httpHandler "loan_system/internal/delivery/http"
	"loan_system/internal/model"
	loanmock "loan_system/internal/usecase/loan/mock"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestCreateLoanHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	mockUsecase := loanmock.NewMockUsecase(ctrl)
	handler := httpHandler.NewLoanHandler(mockUsecase)

	t.Run("successful creation", func(t *testing.T) {

		loanExample := &model.Loan{
			Principal:     10000,
			BorrowerID:    1234,
			Rate:          5.0,
			ROI:           6.0,
			AgreementLink: "https://example.com/agreement.pdf",
		}
		mockUsecase.EXPECT().CreateLoan(gomock.Any(), loanExample).Return(nil)

		reqBody := `{"principal":10000,"borrower_id":1234,"rate":5.0,"roi":6.0,"agreement_link":"https://example.com/agreement.pdf"}`
		req := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(reqBody))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		assert.NoError(t, handler.CreateLoan(e.NewContext(req, rec)))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid param", func(t *testing.T) {
		reqBody := `{"borrower_id":1234,"rate":5.0,"roi":6.0,"agreement_link":"https://example.com/agreement.pdf"}`
		req := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(reqBody))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		err := handler.CreateLoan(e.NewContext(req, rec))

		assert.ErrorContains(t, err, "required")
	})

	t.Run("failure", func(t *testing.T) {
		mockUsecase.EXPECT().CreateLoan(gomock.Any(), gomock.Any()).Return(errors.New("usecase error"))

		reqBody := `{"principal":10000,"borrower_id":1234,"rate":5.0,"roi":6.0,"agreement_link":"https://example.com/agreement.pdf"}`
		req := httptest.NewRequest(http.MethodPost, "/loans", strings.NewReader(reqBody))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		err := handler.CreateLoan(e.NewContext(req, rec))

		assert.ErrorContains(t, err, "usecase error")
	})

}

func TestApproveLoanHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	mockUsecase := loanmock.NewMockUsecase(ctrl)
	handler := httpHandler.NewLoanHandler(mockUsecase)

	t.Run("success approval", func(t *testing.T) {
		mockUsecase.EXPECT().ApproveLoan(gomock.Any(), int64(1), gomock.Any()).Return(&model.Loan{ID: 1}, nil)

		body := bytes.NewBufferString(`{
		  "validator_id": 1234,
		  "proof_url": "https://example.com/proof.pdf",
		  "approver": "investor123"
		}`)
		req := httptest.NewRequest(http.MethodPut, "/loans/1/approve", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans/:id/approve")
		c.SetParamNames("id")
		c.SetParamValues("1")

		assert.NoError(t, handler.ApproveLoan(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid loan ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/loans/invalid/approve", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.ApproveLoan(c)
		assert.ErrorContains(t, err, "invalid")
	})

	t.Run("invalid param", func(t *testing.T) {
		body := bytes.NewBufferString(`{
		  "proof_url": "https://example.com/proof.pdf",
		  "approver": "investor123"
		}`)
		req := httptest.NewRequest(http.MethodPut, "/loans/1/approve", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans/:id/approve")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.ApproveLoan(c)
		assert.ErrorContains(t, err, "required")
	})

	t.Run("failure", func(t *testing.T) {
		mockUsecase.EXPECT().ApproveLoan(gomock.Any(), int64(1), gomock.Any()).Return(nil, errors.New("usecase error"))

		body := bytes.NewBufferString(`{
		  "validator_id": 1234,
		  "proof_url": "https://example.com/proof.pdf",
		  "approver": "investor123"
		}`)
		req := httptest.NewRequest(http.MethodPut, "/loans/1/approve", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans/:id/approve")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.ApproveLoan(c)
		assert.ErrorContains(t, err, "usecase error")
	})
}

func TestAddInvestmentHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	mockUsecase := loanmock.NewMockUsecase(ctrl)
	handler := httpHandler.NewLoanHandler(mockUsecase)

	t.Run("success add investment", func(t *testing.T) {
		mockUsecase.EXPECT().AddInvestment(gomock.Any(), int64(1), gomock.Any()).Return(&model.Loan{ID: 1}, nil)

		body := bytes.NewBufferString(`{
		  "amount": 5000,
		  "investor_id": 1234
		}`)
		req := httptest.NewRequest(http.MethodPut, "/loans/1/invest", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans/:id/invest")
		c.SetParamNames("id")
		c.SetParamValues("1")

		assert.NoError(t, handler.AddInvestment(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid loan ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/loans/invalid/invest", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.AddInvestment(c)
		assert.ErrorContains(t, err, "invalid")
	})

	t.Run("invalid param", func(t *testing.T) {
		body := bytes.NewBufferString(`{
		  "investor_id": 1234
		}`)
		req := httptest.NewRequest(http.MethodPut, "/loans/1/invest", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans/:id/invest")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.AddInvestment(c)
		assert.ErrorContains(t, err, "required")
	})

	t.Run("failure", func(t *testing.T) {
		mockUsecase.EXPECT().AddInvestment(gomock.Any(), int64(1), gomock.Any()).Return(nil, errors.New("usecase error"))

		body := bytes.NewBufferString(`{
		  "amount": 5000,
		  "investor_id": 1234
		}`)
		req := httptest.NewRequest(http.MethodPut, "/loans/1/invest", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans/:id/invest")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.AddInvestment(c)
		assert.ErrorContains(t, err, "usecase error")
	})
}

func TestDisburseLoanHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	mockUsecase := loanmock.NewMockUsecase(ctrl)
	handler := httpHandler.NewLoanHandler(mockUsecase)

	t.Run("success disburse loan", func(t *testing.T) {
		mockUsecase.EXPECT().DisburseLoan(gomock.Any(), int64(1), gomock.Any()).Return(&model.Loan{ID: 1}, nil)
		body := bytes.NewBufferString(`{
		  "id": 1,
		  "officer_id": 1234,
		  "agreement_url": "https://example.com/agreement.pdf",
		  "disbursed_at": "2023-01-01T12:00:00Z"
		}`)

		req := httptest.NewRequest(http.MethodPut, "/loans/1/disburse", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans/:id/disburse")
		c.SetParamNames("id")
		c.SetParamValues("1")

		assert.NoError(t, handler.DisburseLoan(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid loan ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/loans/invalid/disburse", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.DisburseLoan(c)
		assert.ErrorContains(t, err, "invalid")
	})

	t.Run("invalid param", func(t *testing.T) {
		body := bytes.NewBufferString(`{
		  "agreement_url": "https://example.com/agreement.pdf",
		  "disbursed_at": "2023-01-01T12:00:00Z"
		}`)

		req := httptest.NewRequest(http.MethodPut, "/loans/1/disburse", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans/:id/disburse")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.DisburseLoan(c)
		assert.ErrorContains(t, err, "required")
	})

	t.Run("failure", func(t *testing.T) {
		mockUsecase.EXPECT().DisburseLoan(gomock.Any(), int64(1), gomock.Any()).Return(nil, errors.New("usecase error"))

		body := bytes.NewBufferString(`{
		  "id": 1,
		  "officer_id": 1234,
		  "agreement_url": "https://example.com/agreement.pdf",
		  "disbursed_at": "2023-01-01T12:00:00Z"
		}`)

		req := httptest.NewRequest(http.MethodPut, "/loans/1/disburse", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans/:id/disburse")
		c.SetParamNames("id")
		c.SetParamValues("1")

		err := handler.DisburseLoan(c)
		assert.ErrorContains(t, err, "usecase error")
	})
}

func TestGetLoanHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	mockUsecase := loanmock.NewMockUsecase(ctrl)
	handler := httpHandler.NewLoanHandler(mockUsecase)

	t.Run("success get loan", func(t *testing.T) {
		mockUsecase.EXPECT().FindByID(gomock.Any(), int64(1)).Return(&model.Loan{ID: 1}, nil)

		req := httptest.NewRequest(http.MethodGet, "/loans/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		assert.NoError(t, handler.GetLoan(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("invalid loan ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/loans/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		err := handler.GetLoan(c)
		assert.ErrorContains(t, err, "invalid")
	})

}

func TestGetLoansHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	mockUsecase := loanmock.NewMockUsecase(ctrl)
	handler := httpHandler.NewLoanHandler(mockUsecase)

	t.Run("success get loans", func(t *testing.T) {
		mockUsecase.EXPECT().FindAll(gomock.Any()).Return([]*model.Loan{{ID: 1}, {ID: 2}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/loans", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans")

		assert.NoError(t, handler.GetLoans(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("failure", func(t *testing.T) {
		mockUsecase.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("usecase error"))

		req := httptest.NewRequest(http.MethodGet, "/loans", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/loans")

		err := handler.GetLoans(c)
		assert.ErrorContains(t, err, "usecase error")
	})
}
