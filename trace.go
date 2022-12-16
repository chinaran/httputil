package httputil

import (
	"net/http"
)

// WithTraceHeaders forward trace headers
func ForwardHeaders(headers http.Header) map[string]string {
	forwardHeaders := make(map[string]string)
	incommingHeaders := []string{
		// All applications should propagate x-request-id. This header is
		// included in access log statements and is used for consistent trace
		// sampling and log sampling decisions in Istio.
		"x-request-id",

		// Lightstep tracing header. Propagate this if you use lightstep tracing
		// in Istio (see
		// https://istio.io/latest/docs/tasks/observability/distributed-tracing/lightstep/)
		// Note: this should probably be changed to use B3 or W3C TRACE_CONTEXT.
		// Lightstep recommends using B3 or TRACE_CONTEXT and most application
		// libraries from lightstep do not support x-ot-span-context.
		"x-ot-span-context",

		// Datadog tracing header. Propagate these headers if you use Datadog
		// tracing.
		"x-datadog-trace-id",
		"x-datadog-parent-id",
		"x-datadog-sampling-priority",

		// W3C Trace Context. Compatible with OpenCensusAgent and Stackdriver Istio
		// configurations.
		"traceparent",
		"tracestate",

		// Cloud trace context. Compatible with OpenCensusAgent and Stackdriver Istio
		// configurations.
		"x-cloud-trace-context",

		// Grpc binary trace context. Compatible with OpenCensusAgent nad
		// Stackdriver Istio configurations.
		"grpc-trace-bin",

		// b3 trace headers. Compatible with Zipkin, OpenCensusAgent, and
		// Stackdriver Istio configurations. Commented out since they are
		// propagated by the OpenTracing tracer above.
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-flags",

		// SkyWalking trace headers.
		"sw8",

		// Application-specific headers to forward.
		"user-agent",

		// Context and session specific headers
		"cookie",
		"authorization",
		"jwt",
	}

	for _, key := range incommingHeaders {
		val := headers.Get(key)
		if val != "" {
			forwardHeaders[key] = val
		}
	}

	return forwardHeaders
}
