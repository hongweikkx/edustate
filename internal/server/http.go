package server

import (
	v1 "edustate/api/edustate/v1"
	"edustate/internal/conf"
	"edustate/internal/service"
	prometheus_metrics "edustate/pkg/prometheus"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, analysis *service.AnalysisService, tp trace.TracerProvider, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			tracing.Server(tracing.WithTracerProvider(tp)),
			logging.Server(logger),
			validate.Validator(),
			metrics.Server(
				metrics.WithSeconds(prometheus_metrics.MetricSeconds),
				metrics.WithRequests(prometheus_metrics.MetricRequests),
			),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1.RegisterAnalysisHTTPServer(srv, analysis)
	srv.Handle("/metrics", promhttp.Handler())
	return srv
}
