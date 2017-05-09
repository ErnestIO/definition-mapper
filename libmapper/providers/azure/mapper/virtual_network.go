/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	"github.com/ernestio/ernestprovider/providers/azure/virtualnetwork"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapVirtualNetworks ...
func MapVirtualNetworks(d *definition.Definition) (networks []*components.VirtualNetwork) {
	for _, rg := range d.ResourceGroups {
		for _, network := range rg.VirtualNetworks {
			cs := &components.VirtualNetwork{}
			cs.Name = network.Name
			cs.AddressSpace = network.AddressSpaces
			cs.DNSServerNames = network.DNSServers
			cs.ResourceGroupName = rg.Name
			cs.Location = rg.Location

			for _, subnet := range network.Subnets {
				sn := virtualnetwork.Subnet{}
				sn.Name = subnet.Name
				sn.AddressPrefix = subnet.AddressPrefix
				sn.SecurityGroupName = subnet.SecurityGroup
				cs.Subnets = append(cs.Subnets, sn)
			}

			cs.SetDefaultVariables()

			networks = append(networks, cs)
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

		for _, c := range g.GetComponents().ByType("subnet") {
			s := c.(*components.Subnet)

			if s.ResourceGroupName != rg.Name && s.VirtualNetworkName != dn.Name {
				continue
			}

			ds := definition.Subnet{
				Name:          s.Name,
				SecurityGroup: s.NetworkSecurityGroup,
				AddressPrefix: s.AddressPrefix,
			}

			dn.Subnets = append(dn.Subnets, ds)
		}

		networks = append(networks, dn)
	}

	return networks
}
