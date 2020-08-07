package server

import (
	"context"
	"log"

	envoy_ratelimit_v3 "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
)

type RateLimitServer struct {
}

func NewRateLimitServer() *RateLimitServer {
	return &RateLimitServer{}
}

func (server *RateLimitServer) ShouldRateLimit(ctx context.Context, req *envoy_ratelimit_v3.RateLimitRequest) (*envoy_ratelimit_v3.RateLimitResponse, error) {
	log.Println(req.String())
	return &envoy_ratelimit_v3.RateLimitResponse{
		OverallCode: envoy_ratelimit_v3.RateLimitResponse_OK,
	}, nil
}
