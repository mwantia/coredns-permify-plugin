package template

import (
	"github.com/coredns/caddy"
	"github.com/rschone/corefile2struct/pkg/corefile"
)

type TemplateConfig struct {
	Foo string `cf:"foo" check:"nonempty"`
	Bar string `cf:"bar" default:"Hello World"`
}

func CreatePlugin(c *caddy.Controller) (*TemplatePlugin, error) {
	cfg, err := ParseConfig(c)
	if err != nil {
		return nil, err
	}

	return &TemplatePlugin{
		Cfg: *cfg,
	}, nil
}

func ParseConfig(c *caddy.Controller) (*TemplateConfig, error) {
	var cfg TemplateConfig
	if err := corefile.Parse(c, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
