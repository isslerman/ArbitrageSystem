package marketdata

import (
	"fmt"
	"log/slog"
	"looper/internal/pb/orders"
	"net"

	"google.golang.org/grpc"
)

func GRPCListen(md *MarketData, port string) (*grpc.Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen for gRPC: %w", err)
	}

	s := grpc.NewServer()
	orders.RegisterOrderServiceServer(s, md)

	go func() {
		slog.Info("Market Data gRPC server started", "port", port)
		if err := s.Serve(lis); err != nil {
			slog.Error("gRPC server error", "err", err)
		}
	}()

	return s, nil
}
