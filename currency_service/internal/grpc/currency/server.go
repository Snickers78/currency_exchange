package currency

import (
	"context"
	proto "currency-exchange/currency_service/gen/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	GetExchangeRate(ctx context.Context, currency string) (float64, time.Time, error)
}

type Server struct {
	proto.UnimplementedCurrencyServiceServer
	service Service
}

func Register(gRPC *grpc.Server, service Service) {
	proto.RegisterCurrencyServiceServer(gRPC, &Server{service: service})
}

func (s *Server) GetExchangeRate(ctx context.Context, req *proto.ExchangeRateRequest) (*proto.ExchangeRateResponse, error) {
	//TODO: validate currency name

	rate, ts, err := s.service.GetExchangeRate(ctx, req.CurrencyName)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.ExchangeRateResponse{
		CurrencyName: req.CurrencyName,
		Rate:         rate,
		Timestamp:    ts.Format(time.DateTime),
	}, nil
}
