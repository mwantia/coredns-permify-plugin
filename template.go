package template

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"

	otelsetup "github.com/mwantia/coredns-opentelemetry-plugin/otel"

	"go.opentelemetry.io/otel/codes"
)

func (p TemplatePlugin) ServeDNS(ctx context.Context, writer dns.ResponseWriter, msg *dns.Msg) (int, error) {
	state := request.Request{W: writer, Req: msg.Copy()}

	ctx, span := otelsetup.StartDnsServerTracerSpan(ctx, state,
		"github.com/mwantia/coredns-plugin-template", p.Name())
	defer span.End()

	status, err := p.ServeDnsRequest(ctx, state)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return status, err
}

func (p TemplatePlugin) ServeDnsRequest(ctx context.Context, state request.Request) (int, error) {
	return plugin.NextOrFailure(p.Name(), p.Next, ctx, state.W, state.Req)
}
