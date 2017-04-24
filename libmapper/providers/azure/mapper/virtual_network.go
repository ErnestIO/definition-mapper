/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapVirtualNetworks ...
func MapVirtualNetworks(d *definition.Definition) (networks []*components.VirtualNetwork) {
	for _, rg := range d.ResourceGroups {
		for _, network := range rg.VirtualNetworks {
			cs := components.VirtualNetwork{}
			cs.Name = network.Name
			cs.AddressSpace = network.AddressSpaces
			cs.DNSServerNames = network.DNSServers
			cs.ResourceGroupName = rg.Name

			cs.SetDefaultVariables()

			networks = append(networks, &cs)
		}
	}

	return networks
}

// MapDefinitionVirtualNetworks : ...
func MapDefinitionVirtualNetworks(g *graph.Graph, rg *definition.ResourceGroup) (networks []definition.VirtualNetwork) {
	for _, c := range g.GetComponents().ByType("virtual_network") {
		n := c.(*components.VirtualNetwork)

		if n.ResourceGroupName != rg.Name {
			continue
		}

		dn := definition.VirtualNetwork{
			Name:          n.Name,
			AddressSpaces: n.AddressSpace,
			DNSServers:    n.DNSServerNames,
		}

		networks = append(networks, dn)
	}

	return networks
}
