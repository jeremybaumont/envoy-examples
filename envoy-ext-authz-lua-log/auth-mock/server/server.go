package server

import (
	"context"
	"fmt"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	AuthScenarioHeaderKey = "x-auth-scenario-name"

	AuthAllowScenario                        = "allow"
	AuthDenyScenario                         = "deny"
	AuthErrorScenario                        = "error"
	AuthAllowAndAddStatusHeaderDenyScenario  = "transition-deny"
	AuthAllowAndAddStatusHeaderAllowScenario = "transition-allow"

	AuthDenyAndAddHeaderKey    = "x-extauthz-status-code"
	AuthDenyAndAddHeaderValue  = "403"
	AuthAllowAndAddHeaderKey   = "x-extauthz-status-code"
	AuthAllowAndAddHeaderValue = "200"
)

type ScenarioGenerator func(req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error)

type AuthorizationServer struct {
	scenarios map[string]ScenarioGenerator
}

func NewAuthorizationServer() *AuthorizationServer {
	return &AuthorizationServer{
		scenarios: map[string]ScenarioGenerator{
			AuthAllowScenario: func(_ *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
				return &envoy_service_auth_v3.CheckResponse{
					Status: &status.Status{Code: int32(codes.OK)},
				}, nil
			},
			AuthAllowAndAddStatusHeaderDenyScenario: func(_ *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
				return &envoy_service_auth_v3.CheckResponse{
					Status: &status.Status{Code: int32(codes.OK)},
					HttpResponse: &envoy_service_auth_v3.CheckResponse_OkResponse{
						OkResponse: &envoy_service_auth_v3.OkHttpResponse{
							Headers: []*envoy_config_core_v3.HeaderValueOption{
								{
									Append: &wrapperspb.BoolValue{Value: false},
									Header: &envoy_config_core_v3.HeaderValue{
										Key:   AuthDenyAndAddHeaderKey,
										Value: AuthDenyAndAddHeaderValue,
									},
								},
							},
						},
					},
				}, nil
			},
			AuthAllowAndAddStatusHeaderAllowScenario: func(_ *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
				return &envoy_service_auth_v3.CheckResponse{
					Status: &status.Status{Code: int32(codes.OK)},
					HttpResponse: &envoy_service_auth_v3.CheckResponse_OkResponse{
						OkResponse: &envoy_service_auth_v3.OkHttpResponse{
							Headers: []*envoy_config_core_v3.HeaderValueOption{
								{
									Append: &wrapperspb.BoolValue{Value: false},
									Header: &envoy_config_core_v3.HeaderValue{
										Key:   AuthAllowAndAddHeaderKey,
										Value: AuthAllowAndAddHeaderValue,
									},
								},
							},
						},
					},
				}, nil
			},
			AuthDenyScenario: func(_ *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
				return &envoy_service_auth_v3.CheckResponse{
					Status: &status.Status{Code: int32(codes.Unauthenticated)},
					HttpResponse: &envoy_service_auth_v3.CheckResponse_DeniedResponse{
						DeniedResponse: &envoy_service_auth_v3.DeniedHttpResponse{
							Status: &envoy_type_v3.HttpStatus{
								Code: envoy_type_v3.StatusCode_Unauthorized,
							},
						},
					},
				}, nil
			},
			AuthErrorScenario: func(_ *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
				return nil, fmt.Errorf("error from the auth mock")
			},
		},
	}
}

func (server *AuthorizationServer) Check(_ context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
	scenarioName := getScenarioName(req)
	scenarioGenerator, found := server.scenarios[scenarioName]
	if !found {
		return nil, fmt.Errorf("unknown scenario name %q", scenarioName)
	}

	return scenarioGenerator(req)
}

func getScenarioName(req *envoy_service_auth_v3.CheckRequest) string {
	authScenario, found := req.GetAttributes().GetRequest().GetHttp().GetHeaders()[AuthScenarioHeaderKey]
	if !found {
		return AuthAllowScenario
	}

	return authScenario
}