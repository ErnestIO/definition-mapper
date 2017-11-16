/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/definition"
)

// MapGateways : Maps input edge gateway to an ernest formatted edge gateway
func MapGateways(d *definition.Definition) []*components.Gateway {
	var routers []*components.Gateway

	for _, router := range d.Gateways {
		r := &components.Gateway{
			Name: router.Name,
		}

		// Map firewall rules
		for _, rule := range router.FirewallRules {
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
		for _, rule := range router.NatRules {
			r.NatRules = append(r.NatRules, components.NatRule{
				Type:            rule.Type,
				OriginIP:        rule.Source,
				OriginPort:      rule.FromPort,
				TranslationIP:   rule.Destination,
				TranslationPort: rule.ToPort,
				Protocol:        rule.Protocol,
			})
		}

		r.SetDefaultVariables()

		routers = append(routers, r)
	}

	return routers
}
