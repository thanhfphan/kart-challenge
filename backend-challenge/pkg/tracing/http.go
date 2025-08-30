package tracing

import (
	"crypto/tls"
	"net/http/httptrace"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func GetClientTrace(span trace.Span) *httptrace.ClientTrace {
	if !span.IsRecording() {
		return &httptrace.ClientTrace{}
	}
	// fmta(t1.Sub(t0)), // dns lookup
	// fmta(t2.Sub(t1)), // tcp connection
	// fmta(t6.Sub(t5)), // tls handshake
	// fmta(t4.Sub(t3)), // server processing
	// fmta(t7.Sub(t4)), // content transfer
	// fmtb(t1.Sub(t0)), // namelookup
	// fmtb(t2.Sub(t0)), // connect
	// fmtb(t3.Sub(t0)), // pretransfer
	// fmtb(t4.Sub(t0)), // starttransfer
	// fmtb(t7.Sub(t0)), // total
	var t0, t1, t2, t3, t4, t5, t6 time.Time
	return &httptrace.ClientTrace{
		GotConn: func(_ httptrace.GotConnInfo) {
			t3 = time.Now()
			span.SetAttributes(attribute.Int("req.pre_transfer.duration", int(t3.Sub(t0)/time.Millisecond)))
		},
		GotFirstResponseByte: func() {
			t4 = time.Now()
			span.SetAttributes(attribute.Int("req.server.processing.duration", int(t4.Sub(t3)/time.Millisecond)))
			span.SetAttributes(attribute.Int("req.start_transfer.duration", int(t4.Sub(t0)/time.Millisecond)))
		},
		DNSStart: func(_ httptrace.DNSStartInfo) {
			t0 = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			t1 = time.Now()
			span.SetAttributes(attribute.Int("req.dns.lookup.duration", int(t1.Sub(t0)/time.Millisecond)))
			span.SetAttributes(attribute.Int("req.dns.name.lookup.duration", int(t1.Sub(t0)/time.Millisecond)))
		},
		ConnectStart: func(_, _ string) {
			if t0.IsZero() {
				t0 = time.Now()
				t1 = time.Now()
				span.SetAttributes(attribute.Int("req.dns.lookup.duration", int(t1.Sub(t0)/time.Millisecond)))
				span.SetAttributes(attribute.Int("req.dns.name.lookup.duration", int(t1.Sub(t0)/time.Millisecond)))
			}
		},
		ConnectDone: func(net, addr string, err error) {
			t2 = time.Now()
			span.SetAttributes(attribute.Int("req.tcp.connection.duration", int(t2.Sub(t1)/time.Millisecond)))
			span.SetAttributes(attribute.Int("req.connect.duration", int(t2.Sub(t0)/time.Millisecond)))
		},
		TLSHandshakeStart: func() {
			t5 = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			t6 = time.Now()
			span.SetAttributes(attribute.Int("req.tls.handshake.duration", int(t6.Sub(t5)/time.Millisecond)))
		},
	}
}
