/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/definition"
)

// MapRouters : Maps input router to an ernest formatted router
func MapRouters(d *definition.Definition) []*components.Router {
	var routers []*components.Router

	for _, router := range d.Routers {
		r := &components.Router{
			Name: router.Name,
		}

		// Map firewall rules
		for _, rule := range router.Rules {
			snw := d.FindNetwork(rule.Source)
			if snw != nil {
				rule.Source = snw.Subnet
			}

			dnw := d.FindNetwork(rule.Destination)
			if dnw != nil {
				rule.Destination = dnw.Subnet
			}

			r.FirewallRules = append(r.FirewallRules, components.FirewallRule{
				Name:            rule.Name,
				SourceIP:        rule.Source,
				SourcePort:      rule.FromPort,
				DestinationIP:   rule.Destination,
				DestinationPort: rule.ToPort,
				Protocol:        rule.Protocol,
			})
		}

		// Map nat rules
		for _, rule := range router.PortForwarding {
			r.NatRules = append(r.NatRules, components.NatRule{
				Type:            "dnat",
				OriginIP:        rule.Source,
				OriginPort:      rule.FromPort,
				TranslationIP:   rule.Destination,
				TranslationPort: rule.ToPort,
				Protocol:        "tcp",
			})
		}

		r.SetDefaultVariables()

		routers = append(routers, r)
	}

	return routers
}
