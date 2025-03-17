package remote

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	pb "github.com/fflow-tech/fflow/api/foundation/faas"
	"github.com/fflow-tech/fflow/service/pkg/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// DefaultAbilityCallerConfig 默认客户端配置
type DefaultAbilityCallerConfig struct {
	FaasTarget          string `json:"faasTarget,omitempty"`
	FaasAccessToken     string `json:"faasAccessToken,omitempty"`
	LoadBalancingPolicy string `json:"loadBalancingPolicy,omitempty"`
}

type DefaultAbilityCaller struct {
	config     *DefaultAbilityCallerConfig
	faasClient pb.FaasClient
}

func NewDefaultAbilityCaller(config *DefaultAbilityCallerConfig) (*DefaultAbilityCaller, error) {
	conn, err := grpc.Dial(config.FaasTarget,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, config.LoadBalancingPolicy)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &DefaultAbilityCaller{
		config:     config,
		faasClient: pb.NewFaasClient(conn),
	}, nil
}

// CallFAAS 调用 FAAS 能力
func (c *DefaultAbilityCaller) CallFAAS(ctx context.Context, req *CallFAASReqDTO) (map[string]interface{}, error) {
	rsp, err := c.faasClient.Call(ctx, &pb.CallReq{
		BasicReq: &pb.BasicReq{
			Namespace:   req.Namespace,
			AccessToken: c.config.FaasAccessToken,
		},
		Function: req.Function,
		Input:    utils.MapToStr(req.Body),
	})
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal([]byte(rsp.Output), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// CallHTTP 调用 HTTP 能力
func (c *DefaultAbilityCaller) CallHTTP(ctx context.Context, req *CallHTTPReqDTO) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	resp, err := resty.New().R().SetContext(ctx).SetResult(&result).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeaders(req.Header).SetBody(req.Body).SetQueryParams(req.Query).Execute(req.Method, req.URL)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("call %s failed, resp: %d", req.URL, resp.StatusCode())
	}

	return result, nil
}
