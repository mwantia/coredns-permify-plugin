package template

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

type TemplatePlugin struct {
	Next plugin.Handler
	Cfg  TemplateConfig
}

func init() {
	plugin.Register("template", setup)
}

func (p TemplatePlugin) Name() string {
	return "template"
}

func setup(c *caddy.Controller) error {
	p, err := CreatePlugin(c)
	if err != nil {
		return plugin.Error("template", err)
	}

	c.OnStartup(func() error {
		return nil
	})
	c.OnShutdown(func() error {
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		p.Next = next
		return p
	})

	return nil
}
