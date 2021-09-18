package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	envoyauth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"

	"github.com/jeremybaumont/auth-mock/server"
)

type Config struct {
	Port int `required:"true"`
}

type RegisterGrpcServiceFunc = func(*grpc.Server)

func StartGrpcServer(ctx context.Context, port int, registerServiceFunc RegisterGrpcServiceFunc) (<-chan struct{}, error) {
	grpcServer := grpc.NewServer()
	registerServiceFunc(grpcServer)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, fmt.Errorf("binding tcp on port %v: %w", port, err)
	}

	log.Println("Starting grpc server on port: ", port)
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

	var cfg Config
	if err := envconfig.Process("AUTH_MOCK", &cfg); err != nil {
		log.Fatalf("Failed to get auth mock config: %v", err)
	}

	cancelComplete, err := StartGrpcServer(ctx, cfg.Port, func(grpcServer *grpc.Server) {
		envoyauth.RegisterAuthorizationServer(grpcServer, server.NewAuthorizationServer())
	})

	if err != nil {
		log.Fatalf("Failed to start auth mock: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	cancel()
	<-cancelComplete
}
