/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapPublicIPs ...
func MapPublicIPs(d *definition.Definition, rg *components.ResourceGroup) (ips []*components.PublicIP) {
	for _, ip := range d.PublicIPs {
		n := components.PublicIP{}
		n.Name = ip.Name
		n.Location = ip.Location
		n.ResourceGroupName = rg.Name
		n.PublicIPAddressAllocation = ip.PublicIPAddressAllocation
		n.Tags = mapTags(ip.Name, d.Name)

		if ip.ID != "" {
			n.SetAction("none")
		}

		n.SetDefaultVariables()

		ips = append(ips, &n)
	}

	return
}

// MapDefinitionPublicIPs : ...
func MapDefinitionPublicIPs(g *graph.Graph) (ips []definition.PublicIP) {
	for _, c := range g.GetComponents().ByType("public_ip") {
		ip := c.(*components.PublicIP)

		nIP := definition.PublicIP{
			ID:       ip.GetID(),
			Name:     ip.Name,
			Location: ip.Location,
		}

		ips = append(ips, nIP)
	}

	return
}
