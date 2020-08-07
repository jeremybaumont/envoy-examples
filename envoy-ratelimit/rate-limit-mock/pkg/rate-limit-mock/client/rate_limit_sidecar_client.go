package client

import (
	"context"
	"fmt"
	"log"

	envoy_ratelimit_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
	"google.golang.org/grpc"
)

type RateLimitSidecarClient struct {
	address string
}

func NewRateLimitSidecarClient(address string) *RateLimitSidecarClient {
	return &RateLimitSidecarClient{address: address}
}

func (client *RateLimitSidecarClient) SendRequest(rateLimitRequest *envoy_ratelimit_v3.RateLimitRequest) (*envoy_ratelimit_v3.RateLimitResponse, error) {
	grpcConnection, err := grpc.Dial(client.address, grpc.WithInsecure(), grpc.WithUnaryInterceptor(clientLoggingIntercepter))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to the rate-limit sidecar: %w", err)
	}
	defer func() {
		if err := grpcConnection.Close(); err != nil {
			log.Printf("Failed to close grpc connection: %v", err)
		}
	}()
	rateLimitClient := envoy_ratelimit_v3.NewRateLimitServiceClient(grpcConnection)

	response, err := rateLimitClient.ShouldRateLimit(context.Background(), rateLimitRequest)
	if err != nil {
		return nil, fmt.Errorf("sending check request: %w", err)
	}

	return response, nil
}

func clientLoggingIntercepter(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("gRPC client request:%+v\n", req)

	return invoker(ctx, method, req, reply, cc, opts...)
}
