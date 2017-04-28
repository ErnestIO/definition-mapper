/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	"github.com/ernestio/ernestprovider/providers/azure/securitygroup"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapSecurityGroups ...
func MapSecurityGroups(d *definition.Definition) (sgs []*components.SecurityGroup) {
	for _, rg := range d.ResourceGroups {
		for _, sg := range rg.SecurityGroups {
			n := &components.SecurityGroup{}
			n.Name = sg.Name
			n.ResourceGroupName = rg.Name
			n.Tags = mapTags(sg.Name, d.Name)

			for _, rule := range sg.Rules {
				n.SecurityRules = append(n.SecurityRules, securitygroup.SecurityRule{
					Name:                     rule.Name,
					Description:              rule.Description,
					Protocol:                 rule.Protocol,
					SourcePort:               rule.SourcePortRange,
					DestinationPortRange:     rule.DestinationPortRange,
					SourceAddressPrefix:      rule.SourceAddressPrefix,
					DestinationAddressPrefix: rule.DestinationAddressPrefix,
				})
			}

			if sg.ID != "" {
				n.SetAction("none")
			}

			n.SetDefaultVariables()

			sgs = append(sgs, n)
		}
	}

	return
}

// MapDefinitionSecurityGroups : ...
func MapDefinitionSecurityGroups(g *graph.Graph, rg *definition.ResourceGroup) (sgs []definition.SecurityGroup) {
	for _, c := range g.GetComponents().ByType("security_group") {
		sg := c.(*components.SecurityGroup)

		if sg.ResourceGroupName != rg.Name {
			continue
		}

		nSG := definition.SecurityGroup{
			ID:   sg.GetProviderID(),
			Name: sg.Name,
		}

		for _, rule := range sg.SecurityRules {
			nSG.Rules = append(nSG.Rules, definition.SecurityGroupRule{
				Name:                     rule.Name,
				Description:              rule.Description,
				Protocol:                 rule.Protocol,
				SourcePortRange:          rule.SourcePort,
				DestinationPortRange:     rule.DestinationPortRange,
				SourceAddressPrefix:      rule.SourceAddressPrefix,
				DestinationAddressPrefix: rule.DestinationAddressPrefix,
			})
		}

		sgs = append(sgs, nSG)
	}

	return
}
