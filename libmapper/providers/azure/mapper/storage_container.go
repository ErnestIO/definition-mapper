/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
)

// MapStorageContainers ...
func MapStorageContainers(d *definition.Definition) (ips []*components.StorageContainer) {
	for _, rg := range d.ResourceGroups {
		for _, ss := range rg.StorageAccounts {
			for _, sd := range ss.Containers {
				n := &components.StorageContainer{}
				n.Name = sd.Name
				n.ResourceGroupName = rg.Name
				n.StorageAccountName = ss.Name
				n.ContainerAccessType = sd.AccessType

				if n.ID != "" {
					n.SetAction("none")
				}

				n.SetDefaultVariables()

				ips = append(ips, n)
			}
		}
	}

	return
}
