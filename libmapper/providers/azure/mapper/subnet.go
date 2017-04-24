/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapSubnets ...
func MapSubnets(d *definition.Definition) (subnets []*components.Subnet) {
	for _, rg := range d.ResourceGroups {
		for _, vn := range rg.VirtualNetworks {
			for _, subnet := range vn.Subnets {
				cs := components.Subnet{}
				cs.Name = subnet.Name
				cs.AddressPrefix = subnet.AddressPrefix
				cs.NetworkSecurityGroup = subnet.SecurityGroup
				cs.ResourceGroupName = rg.Name
				cs.VirtualNetworkName = vn.Name

				cs.SetDefaultVariables()

				subnets = append(subnets, &cs)
			}
		}
	}

	return subnets
}

// MapDefinitionSubnets : ...
func MapDefinitionSubnets(g *graph.Graph, rg *definition.ResourceGroup, vn *definition.VirtualNetwork) (subnets []definition.Subnet) {
	for _, c := range g.GetComponents().ByType("subnet") {
		s := c.(*components.Subnet)

		if s.ResourceGroupName != rg.Name && s.VirtualNetworkName != vn.Name {
			continue
		}

		ds := definition.Subnet{
			Name:          s.Name,
			SecurityGroup: s.NetworkSecurityGroup,
			AddressPrefix: s.AddressPrefix,
		}

		subnets = append(subnets, ds)
	}

	return subnets
}
