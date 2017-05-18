/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
)

// MapLBProbes ...
func MapLBProbes(d *definition.Definition) (ps []*components.LBProbe) {
	for _, rg := range d.ResourceGroups {
		for _, loadbalancer := range rg.LBs {
			for _, probe := range loadbalancer.Probes {
				n := &components.LBProbe{}
				n.Name = probe.Name
				n.Port = probe.Port
				n.Protocol = probe.Protocol
				n.RequestPath = probe.RequestPath
				n.NumberOfProbes = probe.MaximumFailures
				n.IntervalInSeconds = probe.Interval
				n.Loadbalancer = loadbalancer.Name
				n.ResourceGroupName = rg.Name

				n.SetDefaultVariables()

				ps = append(ps, n)
			}
		}
	}

	return
}
