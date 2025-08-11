package v1

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/usecase"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	EndpointLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "library_enpoint_latency_s",
			Help:    "Latency of endpoint",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
	EndpointRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "library_endpoint_requests",
			Help: "Endpoints request count",
		},
		[]string{"endpoint"},
	)
)

func init() {
	prometheus.MustRegister(
		EndpointLatency,
		EndpointRequests,
	)
}

var tracer = otel.Tracer("library-service")

var _ generated.LibraryServer = (*ControllerImpl)(nil)

type ControllerImpl struct {
	zap           *zap.Logger
	bookUseCase   usecase.BookUseCase
	authorUseCase usecase.AuthorUseCase
}

func NewControllerImpl(
	logger *zap.Logger,
	bookUseCase usecase.BookUseCase,
	authorUseCase usecase.AuthorUseCase,
) *ControllerImpl {
	return &ControllerImpl{
		zap:           logger,
		bookUseCase:   bookUseCase,
		authorUseCase: authorUseCase,
	}
}
