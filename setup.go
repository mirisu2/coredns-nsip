package nsip

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/pkg/log"
	"strings"
)

const pluginName = "Record"

func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	a, err := parse(c)
	if err != nil {
		log.Error(err)
		return plugin.Error(pluginName, err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		a.Next = next
		return a
	})

	return nil
}

func parse(c *caddy.Controller) (Record, error) {
	a := Record{}
	for c.Next() {
		r := rule{}
		args := c.RemainingArgs()
		r.zones = plugin.OriginsFromArgsOrServerBlock(args, c.ServerBlockKeys)

		for c.NextBlock() {
			p := policy{}

			p.ns = strings.ToLower(c.Val())

			remainingTokens := c.RemainingArgs()
			p.ip = remainingTokens[0]

			r.policies = append(r.policies, p)
		}
		a.Rules = append(a.Rules, r)
	}
	return a, nil
}
