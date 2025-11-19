package loan

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpHandler "loan_system/internal/delivery/http"
	"loan_system/internal/pkg/config"
	loanRepository "loan_system/internal/repository/loan"
	"loan_system/internal/repository/pubsub"

	loanUsecase "loan_system/internal/usecase/loan"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type application struct {
	httpHandler.LoanHandler
}

func newApplication() application {
	return application{}
}

func (a application) config() application {
	config.Load()
	return a
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func (a application) serveHTTP() {
	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{},
		AllowOriginFunc: func(origin string) (bool, error) {
			return true, nil
		},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	e.GET("/healthcheck", func(c echo.Context) error {
		response := map[string]string{
			"status": "healthy",
		}
		return c.JSON(http.StatusOK, response)
	})

	loanGroup := e.Group("/loans")

	loanGroup.POST("", a.CreateLoan)
	loanGroup.PUT("/:id/approve", a.ApproveLoan)
	loanGroup.POST("/:id/invest", a.AddInvestment)
	loanGroup.PUT("/:id/disburse", a.DisburseLoan)

	loanGroup.GET("/:id", a.GetLoan)
	loanGroup.GET("", a.GetLoans)

	h2s := &http2.Server{}
	h1s := &http.Server{
		Addr:    ":" + config.Instance().App.ServerPort,
		Handler: h2c.NewHandler(e, h2s),
	}

	// Start server
	go func() {
		fmt.Println("Server Started at:", config.Instance().App.ServerPort)
		if err := h1s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	// Accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGTERM (Ctrl+/) is emitted by Docker stop command
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Shutting down server...")
	// Attempt to gracefully shut down the server
	if err := h1s.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("Server gracefully stopped")
}

func (a application) init() application {
	// init repo
	loanRepository := loanRepository.NewRepository()
	// init pubsub mock
	pubsubMock := pubsub.NewMock()

	loanUsecase := loanUsecase.NewUsecase(loanRepository, pubsubMock)

	a.LoanHandler = *httpHandler.NewLoanHandler(loanUsecase)
	return a
}

func Execute() {
	newApplication().config().init().serveHTTP()
}
