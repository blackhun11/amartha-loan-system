package http

import (
	"net/http"

	"loan_system/internal/model"
	"loan_system/internal/model/request"
	"loan_system/internal/usecase/loan"

	"github.com/labstack/echo/v4"
)

type LoanHandler struct {
	uc loan.Usecase
}

func NewLoanHandler(uc loan.Usecase) *LoanHandler {
	return &LoanHandler{uc: uc}
}

func (h *LoanHandler) CreateLoan(c echo.Context) error {
	req := new(request.CreateLoanRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	loan := &model.Loan{
		BorrowerID:    req.BorrowerID,
		Principal:     req.Principal,
		Rate:          req.Rate,
		ROI:           req.ROI,
		AgreementLink: req.AgreementLink,
	}

	err := h.uc.CreateLoan(c.Request().Context(), loan)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return Success(c, http.StatusOK, map[string]interface{}{
		"loan": loan,
	})
}

func (h *LoanHandler) ApproveLoan(c echo.Context) error {
	req := new(request.ApproveLoanRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	approveReq := model.Approval{
		ValidatorID: req.ValidatorID,
		ProofURL:    req.ProofURL,
		ApprovedAt:  req.ApprovedAt,
	}

	loan, err := h.uc.ApproveLoan(c.Request().Context(), req.ID, approveReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return Success(c, http.StatusOK, map[string]interface{}{
		"loan": loan,
	})
}

func (h *LoanHandler) AddInvestment(c echo.Context) error {
	req := new(request.InvestLoanRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	investment := model.Investment{
		InvestorID: req.InvestorID,
		Amount:     req.Amount,
	}

	loan, err := h.uc.AddInvestment(c.Request().Context(), req.ID, investment)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return Success(c, http.StatusOK, map[string]interface{}{
		"loan": loan,
	})
}

func (h *LoanHandler) DisburseLoan(c echo.Context) error {
	req := new(request.DisburseLoanRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	disbursement := model.Disbursement{
		OfficerID:    req.OfficerID,
		AgreementURL: req.AgreementURL,
		DisbursedAt:  req.DisbursedAt,
	}

	loan, err := h.uc.DisburseLoan(c.Request().Context(), req.ID, disbursement)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return Success(c, http.StatusOK, map[string]interface{}{
		"loan": loan,
	})
}
