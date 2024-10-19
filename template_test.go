package template

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/plugin/pkg/dnstest"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	otelsetup "github.com/mwantia/coredns-opentelemetry-plugin/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type TestRecord struct {
	Name  string
	Type  uint16
	Match string
}

func TestDNS(tst *testing.T) {
	if err := OverwriteStdOut("netboxdns"); err != nil {
		tst.Errorf("Unable to overwrite stdout: %v", err)
	}

	c := caddy.NewTestController("dns", `
		template {
		  foo not_empty
		}
	`)

	p, err := CreatePlugin(c)
	if err != nil {
		tst.Errorf("Unable to create plugin: %v", err)
	}

	ctx := context.TODO()
	shutdown, err := otelsetup.SetupOpentelemetry(ctx, otelsetup.OpenTelemtryConfig{
		Endpoint:     "jaeger:4318",
		ServiceName:  "coredns",
		Hostname:     "localhost",
		BatchTimeout: 3 * time.Second,
		BatchSize:    10,
		SamplingRate: 1,
	})

	if err != nil {
		tst.Errorf("Unable to create opentelemetry provider: %v", err)
	}

	RunTests(ctx, tst, p, []TestRecord{
		{
			Name:  dns.Fqdn("bm1-storage1.node.ts"),
			Type:  dns.TypeA,
			Match: "100.68.166.82",
		},
	})

	if shutdown != nil {
		err := shutdown(ctx)
		if err != nil {
			tst.Errorf("Unable to gracefully shutdown plugin: %v", err)
		}
	}
}

func RunTests(ctx context.Context, tst *testing.T, p *TemplatePlugin, tests []TestRecord) {
	for _, ts := range tests {
		tst.Run("Name: "+ts.Name, func(t *testing.T) {
			ctx, span := otelsetup.StartDnsTestTracerSpan(ctx, ts.Name, "A",
				"github.com/mwantia/coredns-netbox-dns-plugin")
			defer span.End()

			req := new(dns.Msg)
			req.SetQuestion(ts.Name, ts.Type)
			rec := dnstest.NewRecorder(&test.ResponseWriter{})

			rcode, err := p.ServeDNS(ctx, rec, req)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())

				tst.Errorf("Expected no error, but got: %v", err)
				return
			}

			span.SetAttributes(
				attribute.String("dns.rcode", dns.RcodeToString[rcode]),
			)

			if rec.Msg == nil || len(rec.Msg.Answer) == 0 {
				tst.Errorf("Expected an answer, but got none")
				return
			}

			answer := rec.Msg.Answer[0]
			address := answer.(*dns.A).A.String()

			if address != ts.Match {
				tst.Errorf("Expected '%s', but received '%s'", ts.Match, address)
				return
			}
		})
	}
}

func OverwriteStdOut(pluginname string) error {
	tempFile, err := os.CreateTemp("", pluginname)
	if err != nil {
		return err
	}

	defer os.Remove(tempFile.Name())
	log.SetOutput(os.Stdout)

	clog.D.Set()
	return nil
}
