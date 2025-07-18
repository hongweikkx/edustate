// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             v5.29.3
// source: edustate/v1/analysis.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationAnalysisAnalyze = "/edustate.v1.Analysis/Analyze"

type AnalysisHTTPServer interface {
	Analyze(context.Context, *AnalyzeRequest) (*AnalyzeReply, error)
}

func RegisterAnalysisHTTPServer(s *http.Server, srv AnalysisHTTPServer) {
	r := s.Route("/")
	r.POST("/edustate/api/v1/analysis/analyze", _Analysis_Analyze0_HTTP_Handler(srv))
}

func _Analysis_Analyze0_HTTP_Handler(srv AnalysisHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AnalyzeRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAnalysisAnalyze)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Analyze(ctx, req.(*AnalyzeRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AnalyzeReply)
		return ctx.Result(200, reply)
	}
}

type AnalysisHTTPClient interface {
	Analyze(ctx context.Context, req *AnalyzeRequest, opts ...http.CallOption) (rsp *AnalyzeReply, err error)
}

type AnalysisHTTPClientImpl struct {
	cc *http.Client
}

func NewAnalysisHTTPClient(client *http.Client) AnalysisHTTPClient {
	return &AnalysisHTTPClientImpl{client}
}

func (c *AnalysisHTTPClientImpl) Analyze(ctx context.Context, in *AnalyzeRequest, opts ...http.CallOption) (*AnalyzeReply, error) {
	var out AnalyzeReply
	pattern := "/edustate/api/v1/analysis/analyze"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationAnalysisAnalyze))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
