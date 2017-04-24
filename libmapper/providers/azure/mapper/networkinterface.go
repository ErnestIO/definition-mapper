/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	"github.com/ernestio/ernestprovider/providers/azure/networkinterface"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapNetworkInterfaces ...
func MapNetworkInterfaces(d *definition.Definition) (groups []*components.NetworkInterface) {
	for _, ni := range d.NetworkInterfaces {
		cv := components.NetworkInterface{}
		cv.Name = ni.Name
		// cv.ResourceGroupName = group.ResourceGroupName
		cv.NetworkSecurityGroup = ni.SecurityGroup
		cv.DNSServers = ni.DNSServers
		cv.InternalDNSNameLabel = ni.InternalDNSNameLabel
		cv.Tags = mapTags(ni.Name, d.Name)
		for _, ip := range ni.IPConfigurations {
			nIP := networkinterface.IPConfiguration{
				Name:                       ip.Name,
				Subnet:                     ip.Subnet,
				PrivateIPAddress:           ip.PrivateIPAddress,
				PrivateIPAddressAllocation: ip.PrivateIPAddressAllocation,
				PublicIPAddress:            ip.PublicIPAddressID,
			}
			cv.IPConfigurations = append(cv.IPConfigurations, nIP)
		}
		if ni.ID != "" {
			cv.SetAction("none")
		}

		cv.SetDefaultVariables()

		groups = append(groups, &cv)
	}

	return
}

// MapDefinitionNetworkInterfaces : ...
func MapDefinitionNetworkInterfaces(g *graph.Graph) (nis []definition.NetworkInterface) {
	for _, c := range g.GetComponents().ByType("network_interface") {
		ni := c.(*components.NetworkInterface)

		nNi := definition.NetworkInterface{
			ID:                   ni.GetID(),
			Name:                 ni.Name,
			SecurityGroup:        ni.NetworkSecurityGroup,
			DNSServers:           ni.DNSServers,
			InternalDNSNameLabel: ni.InternalDNSNameLabel,
		}

		for _, ip := range ni.IPConfigurations {
			nIP := definition.IPConfiguration{
				Name:                       ip.Name,
				Subnet:                     ip.Subnet,
				PrivateIPAddress:           ip.PrivateIPAddress,
				PrivateIPAddressAllocation: ip.PrivateIPAddressAllocation,
				PublicIPAddressID:          ip.PublicIPAddress,
			}
			nNi.IPConfigurations = append(nNi.IPConfigurations, nIP)
		}

		nis = append(nis, nNi)
	}

	return
}
