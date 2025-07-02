package currency

import (
	"context"
	proto "exchange_service/gen/proto"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Exchanger interface {
	GetExchangeRate(ctx context.Context, base_currency string, target_currency string) (float64, time.Time, error)
	Exchange(ctx context.Context, base_currency string, target_currency string, amount float64) (float64, time.Time, error)
}

type Server struct {
	proto.UnimplementedCurrencyServiceServer
	exchanger Exchanger
}

func Register(gRPC *grpc.Server, exchanger Exchanger) {
	proto.RegisterCurrencyServiceServer(gRPC, &Server{exchanger: exchanger})
}

func (s *Server) GetExchangeRate(ctx context.Context, req *proto.ExchangeRateRequest) (*proto.ExchangeRateResponse, error) {
	rate, ts, err := s.exchanger.GetExchangeRate(ctx, req.BaseCurrency, req.TargetCurrency)
	if rate == 0 {
		return nil, status.Error(codes.NotFound, "currency not found")
	}
	if err != nil {
		fmt.Println("Error getting exchange rate: ", err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.ExchangeRateResponse{
		CurrencyName: req.TargetCurrency,
		Rate:         rate,
		Timestamp:    ts.Format(time.DateTime),
	}, nil
}

func (s *Server) Exchange(ctx context.Context, req *proto.ExchangeRequest) (*proto.ExchangeResponse, error) {
	if req.Amount == 0 {
		return nil, status.Error(codes.InvalidArgument, "amount must be greater than 0")
	}
	exchange_amount, ts, err := s.exchanger.Exchange(ctx, req.BaseCurrency, req.TargetCurrency, req.Amount)
	if exchange_amount == 0 {
		return nil, status.Error(codes.NotFound, "currency not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.ExchangeResponse{
		Currency:  req.TargetCurrency,
		Amount:    exchange_amount,
		Timestamp: ts.Format(time.DateTime),
	}, nil
}
