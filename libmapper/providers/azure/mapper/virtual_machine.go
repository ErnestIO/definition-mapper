/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"strconv"
	"strings"

	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	"github.com/ernestio/ernestprovider/providers/azure/virtualmachine"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapVirtualMachines ...
func MapVirtualMachines(d *definition.Definition) (vms []*components.VirtualMachine) {
	for _, rg := range d.ResourceGroups {
		for _, vm := range rg.VirtualMachines {
			image := getImageParts(vm.Image)

			for i := 1; i < vm.Count+1; i++ {
				cvm := components.VirtualMachine{}
				cvm.Name = vm.Name + "-" + strconv.Itoa(i)
				cvm.VMSize = vm.Size

				if len(image) == 4 {
					vm.StorageImageReference.Publisher = image[0]
					vm.StorageImageReference.Offer = image[1]
					vm.StorageImageReference.Sku = image[2]
					vm.StorageImageReference.Version = image[3]
				}

				cvm.NetworkInterfaces = vm.NetworkInterfaces

				cvm.StorageOSDisk.Name = vm.StorageOSDisk.Name
				cvm.StorageOSDisk.Caching = vm.StorageOSDisk.Caching
				cvm.StorageOSDisk.OSType = vm.StorageOSDisk.OSType
				cvm.StorageOSDisk.CreateOption = vm.StorageOSDisk.CreateOption
				cvm.StorageOSDisk.ImageURI = vm.StorageOSDisk.ImageURI
				cvm.StorageOSDisk.VhdURI = vm.StorageOSDisk.VHDURI

				cvm.StorageDataDisk.Name = vm.StorageDataDisk.Name
				cvm.StorageDataDisk.Size = vm.StorageDataDisk.DiskSizeGB
				cvm.StorageDataDisk.VhdURI = vm.StorageDataDisk.VhdURI
				cvm.StorageDataDisk.CreateOption = vm.StorageDataDisk.CreateOption

				cvm.DeleteDataDisksOnTermination = vm.DeleteDataDisksOnTermination
				cvm.DeleteOSDiskOnTermination = vm.DeleteOSDiskOnTermination

				cvm.BootDiagnostics = []virtualmachine.BootDiagnostic{
					virtualmachine.BootDiagnostic{
						Enabled: vm.BootDiagnostics.Enabled,
						URI:     vm.BootDiagnostics.StorageURI,
					},
				}

				cvm.Plan.Name = vm.Plan.Name
				cvm.Plan.Product = vm.Plan.Product
				cvm.Plan.Publisher = vm.Plan.Publisher

				cvm.OSProfileLinuxConfig.SSHKeys = mapSSHKeys(vm.Authentication.SSHKeys)
				cvm.OSProfileLinuxConfig.DisablePasswordAuthentication = vm.Authentication.DisablePasswordAuthentication
				cvm.OSProfileWindowsConfig.ProvisionVMAgent = vm.OSProfileWindowsConfig.ProvisionVMAgent
				cvm.OSProfileWindowsConfig.EnableAutomaticUpgrades = vm.OSProfileWindowsConfig.EnableAutomaticUpgrades
				cvm.OSProfile.AdminPassword = vm.Authentication.AdminUsername
				cvm.OSProfile.AdminPassword = vm.Authentication.AdminPassword

				for _, winrm := range vm.OSProfileWindowsConfig.WinRM {
					cvm.OSProfileWindowsConfig.WinRm = append(cvm.OSProfileWindowsConfig.WinRm, virtualmachine.WinRM{
						Protocol:       winrm.Protocol,
						CertificateURL: winrm.CertificateURL,
					})
				}

				cvm.OSProfileWindowsConfig.AdditionalUnattendConfig = append(cvm.OSProfileWindowsConfig.AdditionalUnattendConfig, virtualmachine.UnattendedConfig{
					Pass:        vm.OSProfileWindowsConfig.AdditionalUnattendConfig.Pass,
					Component:   vm.OSProfileWindowsConfig.AdditionalUnattendConfig.Component,
					SettingName: vm.OSProfileWindowsConfig.AdditionalUnattendConfig.SettingName,
					Content:     vm.OSProfileWindowsConfig.AdditionalUnattendConfig.Content,
				})

				cvm.Tags = vm.Tags
				cvm.LicenseType = vm.LicenseType
				cvm.ResourceGroupName = rg.Name

				cvm.SetDefaultVariables()

				vms = append(vms, &cvm)
			}
		}
	}

	return vms
}

// MapDefinitionVirtualMachines : ...
func MapDefinitionVirtualMachines(g *graph.Graph, rg *definition.ResourceGroup) (vms []definition.VirtualMachine) {
	for _, c := range g.GetComponents().ByType("virtual_machine") {
		vm := c.(*components.VirtualMachine)

		if vm.ResourceGroupName != rg.Name {
			continue
		}

		image := vm.StorageImageReference

		dvm := definition.VirtualMachine{
			Name:              vm.Name,
			Size:              vm.VMSize,
			Image:             strings.Join([]string{image.Publisher, image.Offer, image.Sku, image.Version}, ":"),
			NetworkInterfaces: vm.NetworkInterfaces,
			Tags:              vm.Tags,
			LicenseType:       vm.LicenseType,
		}

		dvm.StorageOSDisk.Name = vm.StorageOSDisk.Name
		dvm.StorageOSDisk.Caching = vm.StorageOSDisk.Caching
		dvm.StorageOSDisk.OSType = vm.StorageOSDisk.OSType
		dvm.StorageOSDisk.CreateOption = vm.StorageOSDisk.CreateOption
		dvm.StorageOSDisk.ImageURI = vm.StorageOSDisk.ImageURI
		dvm.StorageOSDisk.VHDURI = vm.StorageOSDisk.VhdURI

		dvm.StorageDataDisk.Name = vm.StorageDataDisk.Name
		dvm.StorageDataDisk.DiskSizeGB = vm.StorageDataDisk.Size
		dvm.StorageDataDisk.VhdURI = vm.StorageDataDisk.VhdURI
		dvm.StorageDataDisk.CreateOption = vm.StorageDataDisk.CreateOption

		if len(vm.BootDiagnostics) > 0 {
			dvm.BootDiagnostics.Enabled = vm.BootDiagnostics[0].Enabled
			dvm.BootDiagnostics.StorageURI = vm.BootDiagnostics[0].URI
		}

		dvm.Plan.Name = vm.Plan.Name
		dvm.Plan.Product = vm.Plan.Product
		dvm.Plan.Publisher = vm.Plan.Publisher

		dvm.Authentication.SSHKeys = mapDefinitionSSHKeys(vm.OSProfileLinuxConfig.SSHKeys)
		dvm.Authentication.DisablePasswordAuthentication = vm.OSProfileLinuxConfig.DisablePasswordAuthentication
		dvm.OSProfileWindowsConfig.ProvisionVMAgent = vm.OSProfileWindowsConfig.ProvisionVMAgent
		dvm.OSProfileWindowsConfig.EnableAutomaticUpgrades = vm.OSProfileWindowsConfig.EnableAutomaticUpgrades
		dvm.Authentication.AdminUsername = vm.OSProfile.AdminPassword
		dvm.Authentication.AdminPassword = vm.OSProfile.AdminPassword

		vms = append(vms, dvm)
	}

	return vms
}

func mapSSHKeys(keyList map[string]string) (keys []virtualmachine.SSHKey) {
	for path, key := range keyList {
		keys = append(keys, virtualmachine.SSHKey{
			Path:    path,
			KeyData: key,
		})
	}

	return
}

func mapDefinitionSSHKeys(keyList []virtualmachine.SSHKey) map[string]string {
	keys := make(map[string]string)
	for _, key := range keyList {
		keys[key.Path] = key.KeyData
	}

	return keys
}

func getImageParts(image string) []string {
	return strings.Split(image, ":")
}
