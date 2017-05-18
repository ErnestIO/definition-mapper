/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
)

// MapLBBackendAddressPools ...
func MapLBBackendAddressPools(d *definition.Definition) (aps []*components.LBBackendAddressPool) {
	for _, rg := range d.ResourceGroups {
		for _, loadbalancer := range rg.LBs {
			for _, ap := range loadbalancer.BackendAddressPools {
				n := &components.LBBackendAddressPool{}
				n.Name = ap
				n.Loadbalancer = loadbalancer.Name
				n.ResourceGroupName = rg.Name

				n.SetDefaultVariables()

				aps = append(aps, n)
			}
		}
	}

	return
}
