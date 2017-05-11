/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"strconv"

	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
)

// MapPublicIPs ...
func MapPublicIPs(d *definition.Definition) (ips []*components.PublicIP) {
	for _, rg := range d.ResourceGroups {
		for _, vm := range rg.VirtualMachines {
			for _, iface := range vm.NetworkInterfaces {
				for x, config := range iface.IPConfigurations {
					if config.PublicIPAddressAllocation == "" {
						continue
					}

					for i := 1; i < vm.Count+1; i++ {
						n := &components.PublicIP{}
						n.Name = iface.Name + "-" + strconv.Itoa(i) + "-" + strconv.Itoa(x+1)
						n.Location = rg.Location
						n.ResourceGroupName = rg.Name
						n.PublicIPAddressAllocation = config.PublicIPAddressAllocation
						n.Tags = mapTags(n.Name, d.Name)

						n.SetDefaultVariables()

						ips = append(ips, n)
					}
				}
			}
		}

		for _, lb := range rg.LBs {
			for _, config := range lb.FrontendIPConfigurations {
				if config.PublicIPAddressAllocation == "" {
					continue
				}

				n := &components.PublicIP{}
				n.Name = lb.Name
				n.Location = rg.Location
				n.ResourceGroupName = rg.Name
				n.PublicIPAddressAllocation = config.PublicIPAddressAllocation
				n.Tags = mapTags(config.Name, d.Name)

				n.SetDefaultVariables()

				ips = append(ips, n)
			}
		}
	}

	return
}
