/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
)

// MapAvailabilitySets ...
func MapAvailabilitySets(d *definition.Definition) (sets []*components.AvailabilitySet) {
	for _, rg := range d.ResourceGroups {
		for _, set := range rg.AvailabilitySets {
			s := &components.AvailabilitySet{}
			s.Name = set.Name
			s.PlatformFaultDomainCount = set.FaultDomainCount
			s.PlatformUpdateDomainCount = set.UpdateDomainCount
			s.Managed = set.Managed
			s.ResourceGroupName = rg.Name
			s.Location = rg.Location

			s.SetDefaultVariables()

			sets = append(sets, s)
		}
	}

	return
}
