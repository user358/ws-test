package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"ws-test/generated/ws_test"
	"ws-test/internal/log"
)

type BalanceService interface {
	GetByUserID(userID int64) decimal.Decimal
}

type Server struct {
	port int
	ws_test.UnimplementedWSTestServer
	balanceService BalanceService
}

func New(port int, balanceService BalanceService) *Server {
	return &Server{port: port, balanceService: balanceService}
}

func (r *Server) Run(ctx context.Context) {
	server := grpc.NewServer(grpc.UnaryInterceptor(logInterceptor))
	ws_test.RegisterWSTestServer(server, r)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", r.port))
	if err != nil {
		log.Fatalw("failed to listen grpc", "err", err, "port", r.port)
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatalw("failed to start grpc", "err", err)
		}
	}()
	log.Infow("listening grpc", "port", r.port)

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	log.Infow("grpc shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	go server.GracefulStop()

	<-ctx.Done()
	server.Stop()

	log.Infow("grpc server exiting")
}

func (r *Server) GetUserBalance(_ context.Context, request *ws_test.GetUserBalanceRequest) (*ws_test.GetUserBalanceResponse, error) {
	v := r.balanceService.GetByUserID(request.UserId)
	output := &ws_test.GetUserBalanceResponse{
		Value: v.String(),
	}
	return output, nil
}
