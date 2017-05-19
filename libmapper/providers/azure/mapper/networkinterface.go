/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"net"
	"strconv"
	"strings"

	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	"github.com/ernestio/ernestprovider/providers/azure/networkinterface"
)

// MapNetworkInterfaces ...
func MapNetworkInterfaces(d *definition.Definition) (interfaces []*components.NetworkInterface) {
	for _, rg := range d.ResourceGroups {
		for _, vm := range rg.VirtualMachines {
			for _, ni := range vm.NetworkInterfaces {
				var addresses []net.IP

				for _, config := range ni.IPConfigurations {
					addresses = append(addresses, net.ParseIP(config.PrivateIPAddress).To4())
				}

				for i := 1; i < vm.Count+1; i++ {
					cv := &components.NetworkInterface{}
					cv.Name = ni.Name + "-" + strconv.Itoa(i)
					cv.NetworkSecurityGroup = ni.SecurityGroup
					cv.DNSServers = ni.DNSServers
					cv.InternalDNSNameLabel = ni.InternalDNSNameLabel
					cv.ResourceGroupName = rg.Name
					cv.VirtualMachineID = components.TYPEVIRTUALMACHINE + components.TYPEDELIMITER + vm.Name
					cv.Location = rg.Location
					cv.Tags = mapTags(ni.Name, d.Name)
					for k, v := range ni.Tags {
						cv.Tags[k] = v
					}

					for x, ip := range ni.IPConfigurations {
						subnet := strings.Split(ip.Subnet, ":")[1]

						nIP := networkinterface.IPConfiguration{
							Name:   ip.Name,
							Subnet: subnet,
							PrivateIPAddressAllocation:      ip.PrivateIPAddressAllocation,
							LoadbalancerBackendAddressPools: ip.LoadBalancerBackendAddressPools,
						}
						if nIP.PrivateIPAddressAllocation == "" {
							nIP.PrivateIPAddressAllocation = "static"
						}

						if nIP.PrivateIPAddressAllocation == "static" {
							nIP.PrivateIPAddress = addresses[x].String()
							addresses[x][3]++
						}

						if ip.PublicIPAddressAllocation != "" {
							nIP.PublicIPAddress = cv.Name + "-" + ip.Name
						}

						cv.IPConfigurations = append(cv.IPConfigurations, nIP)
					}

					if ni.ID != "" {
						cv.SetAction("none")
					}

					cv.SetDefaultVariables()

					interfaces = append(interfaces, cv)
				}
			}
		}
	}

	return
}
