package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"ws-test/internal/log"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// BalanceService contains business logic for balance operations
type BalanceService interface {
	GetByUserID(userID int64) decimal.Decimal
	Deposit(userID int64, value decimal.Decimal) error
	Withdraw(userID int64, value decimal.Decimal) error
}

// Server implements http handlers and run the http server
type Server struct {
	port                int
	balanceService      BalanceService
	notificationService NotificationService
}

func New(port int, balanceService BalanceService, notificationService NotificationService) *Server {
	return &Server{port: port, balanceService: balanceService, notificationService: notificationService}
}

// Run runs http server until done signal will be received from the context
func (r *Server) Run(ctx context.Context) {
	router := gin.Default()
	router.GET("/api/wallet/balance/:user_id", r.handleBalance)
	router.POST("/api/wallet/deposit", r.handleDeposit)
	router.POST("/api/wallet/withdraw", r.handleWithdraw)
	router.GET("/connect", r.handleWebsocket)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", r.port),
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("failed to listen http", "err", err, "port", r.port)
		}
	}()
	log.Infow("listening http", "port", r.port)

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	log.Infow("http shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalw("http server forced to shutdown", "err", err)
	}

	log.Infow("http server exiting")
}
