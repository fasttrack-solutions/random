package main

import (
	"context"
	"fmt"
	"github.com/fasttrack-solutions/random"
	"github.com/fasttrack-solutions/random/internal/config"
	"github.com/fasttrack-solutions/random/pkg/pb"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
	"runtime/debug"
)

func main() {
	seed := *config.SEEDHEX
	if len(seed) != 64 {
		panic("seed must be 64 hex characters")
	} else if seed == "0000000000000000000000000000000000000000000000000000000000000000" {
		panic("a unique seed value is required")
	}

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			slog.Error("[PANIC] recovered panic", "error", p, "stacktrace", string(debug.Stack()))
			return fmt.Errorf("recovered panic: %v", p)
		}),
	}

	s := grpc.NewServer(
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_recovery.StreamServerInterceptor(recoveryOpts...),
			),
		),
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
	)

	reflection.Register(s)

	randomServer := NewRandomGRPCServer(seed)
	pb.RegisterRandomServer(s, randomServer)

	lis, errListen := net.Listen("tcp", fmt.Sprintf(":%v", *config.GRPCPort))
	if errListen != nil {
		slog.Error("failed to start grpc api", "error", errListen.Error())
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("gRPC server listening on %v", lis.Addr()))

	errServe := s.Serve(lis)
	if errServe != nil {
		slog.Error("failed to serve", "error", errServe.Error())
		os.Exit(1)
	}
}

type RandomGRPCServer struct {
	pb.UnimplementedRandomServer
	seed string
}

func NewRandomGRPCServer(seed string) *RandomGRPCServer {
	return &RandomGRPCServer{
		seed: seed,
	}
}

func (rs *RandomGRPCServer) GetRandomInt64(ctx context.Context, req *pb.GetRandomInt64Request) (*pb.GetRandomInt64Response, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	number, err := random.UniformInt64(req.Min, req.Max)
	if err != nil {
		return nil, err
	}

	return &pb.GetRandomInt64Response{
		Number: number,
	}, nil
}

func (rs *RandomGRPCServer) GetRandomFloat64(ctx context.Context, req *pb.GetRandomFloat64Request) (*pb.GetRandomFloat64Response, error) {
	number, err := random.UniformFloat64()
	if err != nil {
		return nil, err
	}

	return &pb.GetRandomFloat64Response{
		Number: number,
	}, nil
}

func (rs *RandomGRPCServer) GetDeterministicRandom(ctx context.Context, req *pb.GetDeterministicRandomRequest) (*pb.GetDeterministicRandomResponse, error) {
	number, err := random.DeterministicRandom(rs.seed, req.Sequence, req.Probabilities)
	if err != nil {
		return nil, err
	}

	return &pb.GetDeterministicRandomResponse{
		Number: number,
	}, nil
}
