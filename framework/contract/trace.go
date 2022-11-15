package contract

import (
	"context"
	"net/http"
)

const TraceKey = "hade:trace"

const (
	TraceKeyTraceID  = "trace_id"
	TraceKeySpanID   = "span_id"
	TraceKeyCspanID  = "cspan_id"
	TraceKeyParentID = "parent_id"
	TraceKeyMethod   = "method"
	TraceKeyCaller   = "caller"
	TraceKeyTime     = "time"
)

type TraceContext struct {
	TraceID string
	ParentID string
	SpanID string
	CspanID string

	Annotation map[string]string
}

type Trace interface {
	// // SetTraceIDService set TraceID generator, default hadeIDGenerator
	// SetTraceIDService(IDService)
	// // SetTraceIDService set SpanID generator, default hadeIDGenerator
	// SetSpanIDService(IDService)

	// WithContext register new trace to context
	WithTrace(c context.Context, trace *TraceContext) context.Context
	// GetTrace From trace context
	GetTrace(c context.Context) *TraceContext
	// NewTrace generate a new trace
	NewTrace() *TraceContext
	// StartSpan generate cspan for child call
	StartSpan(trace *TraceContext) *TraceContext

	// traceContext to map for logger
	ToMap(trace *TraceContext) map[string]string

	// GetTrace By Http
	ExtractHTTP(req *http.Request) *TraceContext
	// Set Trace to Http
	InjectHTTP(req *http.Request, trace *TraceContext) *http.Request
}

