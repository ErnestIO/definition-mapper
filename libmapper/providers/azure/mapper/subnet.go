/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
)

// MapSubnets ...
func MapSubnets(d *definition.Definition) (subnets []*components.Subnet) {
	for _, rg := range d.ResourceGroups {
		for _, vn := range rg.VirtualNetworks {
			for _, subnet := range vn.Subnets {
				cs := &components.Subnet{}
				cs.Name = subnet.Name
				cs.AddressPrefix = subnet.AddressPrefix
				cs.NetworkSecurityGroup = subnet.SecurityGroup
				cs.ResourceGroupName = rg.Name
				cs.VirtualNetworkName = vn.Name

				cs.SetDefaultVariables()

				subnets = append(subnets, cs)
			}
		}
	}

	return subnets
}
