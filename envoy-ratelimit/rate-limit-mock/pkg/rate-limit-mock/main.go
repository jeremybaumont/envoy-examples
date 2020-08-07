package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	envoy_service_ratelimit_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
	"google.golang.org/grpc"

	"github.com/jeremybaumont/envoy-examples/envoy-ratelimit/rate-limit-mock/pkg/rate-limit-mock/server"
)

const (
	RateLimitPortEnvVar = "RATE_LIMIT_MOCK_PORT"
)

func startRateLimitMock(ctx context.Context) (<-chan struct{}, error) {
	port, found := os.LookupEnv(RateLimitPortEnvVar)

	if !found || port == "" {
		return nil, fmt.Errorf("%v environment variable must be set", RateLimitPortEnvVar)
	}

	log.Printf("Starting rate-limit grpcServer on %v\n", port)

	grpcServer := grpc.NewServer()
	envoy_service_ratelimit_v3.RegisterRateLimitServiceServer(grpcServer, server.NewRateLimitServer())
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("binding tcp on port %v: %w", port, err)
	}

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Failed to start grpc server: %v", err)
		}
	}()

	cancelComplete := make(chan struct{})

	go func() {
		<-ctx.Done()
		log.Println("Shutting down grpc server")
		grpcServer.Stop()
		close(cancelComplete)
	}()

	return cancelComplete, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancelComplete, err := startRateLimitMock(ctx)
	if err != nil {
		log.Fatalf("Failed to start rate-limit mock: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	cancel()
	<-cancelComplete
}
