/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"strconv"

	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
)

// MapManagedDisks ...
func MapManagedDisks(d *definition.Definition) (mds []*components.ManagedDisk) {
	for _, rg := range d.ResourceGroups {
		for _, vm := range rg.VirtualMachines {
			if vm.Count == 0 {
				vm.Count = 1
			}
			if vm.StorageOSDisk.ManagedDiskType != "" {
				for i := 1; i < vm.Count+1; i++ {
					md := &components.ManagedDisk{}
					md.Name = vm.Name + "-" + strconv.Itoa(i) + "-" + vm.StorageOSDisk.Name
					md.ResourceGroupName = rg.Name
					md.Location = rg.Location
					md.StorageAccountType = vm.StorageOSDisk.ManagedDiskType
					md.CreateOption = vm.StorageOSDisk.CreateOption
					md.SourceURI = vm.StorageOSDisk.ImageURI
					md.OSType = vm.StorageOSDisk.OSType
					md.DiskSizeGB = vm.StorageOSDisk.DiskSizeGB
					md.Tags = mapTags(md.Name, d.Name)
					for k, v := range vm.Tags {
						md.Tags[k] = v
					}

					if md.ID != "" {
						md.SetAction("none")
					}

					md.SetDefaultVariables()

					mds = append(mds, md)
				}
			}

			if vm.StorageDataDisk.ManagedDiskType != "" {
				for i := 1; i < vm.Count+1; i++ {
					md := &components.ManagedDisk{}
					md.Name = vm.Name + "-" + strconv.Itoa(i) + "-" + vm.StorageDataDisk.Name
					md.ResourceGroupName = rg.Name
					md.Location = rg.Location
					md.StorageAccountType = vm.StorageDataDisk.ManagedDiskType
					md.CreateOption = vm.StorageDataDisk.CreateOption
					md.SourceURI = vm.StorageDataDisk.ImageURI
					md.SourceResourceID = vm.StorageDataDisk.StorageSourceResourceID
					md.OSType = vm.StorageDataDisk.OSType
					md.DiskSizeGB = *vm.StorageDataDisk.DiskSizeGB
					md.Tags = mapTags(md.Name, d.Name)
					for k, v := range vm.Tags {
						md.Tags[k] = v
					}

					if md.ID != "" {
						md.SetAction("none")
					}

					md.SetDefaultVariables()

					mds = append(mds, md)
				}
			}

		}
	}

	return
}
