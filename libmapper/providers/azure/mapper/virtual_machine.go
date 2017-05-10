/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"fmt"
	"net/url"
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
				cvm := &components.VirtualMachine{}
				cvm.Name = vm.Name + "-" + strconv.Itoa(i)
				cvm.VMSize = vm.Size

				if len(image) == 4 {
					cvm.StorageImageReference.Publisher = image[0]
					cvm.StorageImageReference.Offer = image[1]
					cvm.StorageImageReference.Sku = image[2]
					cvm.StorageImageReference.Version = image[3]
				}

				for _, ni := range vm.NetworkInterfaces {
					cvm.NetworkInterfaces = append(cvm.NetworkInterfaces, ni.Name+"-"+strconv.Itoa(i))
				}

				cvm.StorageOSDisk.Name = vm.StorageOSDisk.Name
				cvm.StorageOSDisk.Caching = vm.StorageOSDisk.Caching
				cvm.StorageOSDisk.OSType = vm.StorageOSDisk.OSType
				cvm.StorageOSDisk.CreateOption = vm.StorageOSDisk.CreateOption
				cvm.StorageOSDisk.ImageURI = vm.StorageOSDisk.ImageURI
				cvm.StorageOSDisk.VhdURI = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s.vhd", vm.StorageOSDisk.StorageAccount, vm.StorageOSDisk.StorageContainer, vm.StorageOSDisk.Name+"-"+strconv.Itoa(i))
				cvm.StorageOSDisk.StorageAccount = vm.StorageOSDisk.StorageAccount
				cvm.StorageOSDisk.StorageContainer = vm.StorageOSDisk.StorageContainer

				cvm.StorageDataDisk.Name = vm.StorageDataDisk.Name
				cvm.StorageDataDisk.Size = vm.StorageDataDisk.DiskSizeGB
				cvm.StorageDataDisk.CreateOption = vm.StorageDataDisk.CreateOption
				cvm.StorageDataDisk.VhdURI = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s.vhd", vm.StorageDataDisk.StorageAccount, vm.StorageDataDisk.StorageContainer, vm.StorageDataDisk.Name+"-"+strconv.Itoa(i))
				cvm.StorageDataDisk.StorageAccount = vm.StorageDataDisk.StorageAccount
				cvm.StorageDataDisk.StorageContainer = vm.StorageDataDisk.StorageContainer

				cvm.DeleteDataDisksOnTermination = vm.DeleteDataDisksOnTermination
				cvm.DeleteOSDiskOnTermination = vm.DeleteOSDiskOnTermination

				if vm.BootDiagnostics.Enabled != false {
					cvm.BootDiagnostics = []virtualmachine.BootDiagnostic{
						virtualmachine.BootDiagnostic{
							Enabled: vm.BootDiagnostics.Enabled,
							URI:     vm.BootDiagnostics.StorageURI,
						},
					}
				}

				cvm.Plan.Name = vm.Plan.Name
				cvm.Plan.Product = vm.Plan.Product
				cvm.Plan.Publisher = vm.Plan.Publisher

				cvm.OSProfile.ComputerName = vm.OSProfile.ComputerName

				cvm.OSProfileLinuxConfig.SSHKeys = mapSSHKeys(vm.Authentication.SSHKeys)
				cvm.OSProfileLinuxConfig.DisablePasswordAuthentication = vm.Authentication.DisablePasswordAuthentication
				cvm.OSProfileWindowsConfig.ProvisionVMAgent = vm.OSProfileWindowsConfig.ProvisionVMAgent
				cvm.OSProfileWindowsConfig.EnableAutomaticUpgrades = vm.OSProfileWindowsConfig.EnableAutomaticUpgrades
				cvm.OSProfile.AdminUsername = vm.Authentication.AdminUsername
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
				cvm.Location = rg.Location

				cvm.SetDefaultVariables()

				vms = append(vms, cvm)
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
			Name:        vm.Name,
			Size:        vm.VMSize,
			Image:       strings.Join([]string{image.Publisher, image.Offer, image.Sku, image.Version}, ":"),
			Tags:        vm.Tags,
			LicenseType: vm.LicenseType,
		}

		_, osaccount, oscontainer := getStorageDetails(vm.StorageOSDisk.VhdURI)
		_, dataaccount, datacontainer := getStorageDetails(vm.StorageDataDisk.VhdURI)

		dvm.StorageOSDisk.Name = vm.StorageOSDisk.Name
		dvm.StorageOSDisk.Caching = vm.StorageOSDisk.Caching
		dvm.StorageOSDisk.OSType = vm.StorageOSDisk.OSType
		dvm.StorageOSDisk.CreateOption = vm.StorageOSDisk.CreateOption
		dvm.StorageOSDisk.ImageURI = vm.StorageOSDisk.ImageURI
		dvm.StorageOSDisk.StorageAccount = osaccount
		dvm.StorageOSDisk.StorageContainer = oscontainer

		dvm.StorageDataDisk.Name = vm.StorageDataDisk.Name
		dvm.StorageDataDisk.DiskSizeGB = vm.StorageDataDisk.Size
		dvm.StorageDataDisk.CreateOption = vm.StorageDataDisk.CreateOption
		dvm.StorageDataDisk.StorageAccount = dataaccount
		dvm.StorageDataDisk.StorageContainer = datacontainer

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

		for _, cn := range g.GetComponents().ByType("network_interface") {
			ni := cn.(*components.NetworkInterface)

			if ni.VirtualMachineID != vm.ID {
				continue
			}

			nNi := definition.NetworkInterface{
				ID:                   ni.GetProviderID(),
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
				}
				if ip.PublicIPAddressID != "" {
					cpip := g.GetComponents().ByProviderID(ip.PublicIPAddress)
					if cpip != nil {
						pip := cpip.(*components.PublicIP)
						nIP.PublicIPAddressAllocation = pip.PublicIPAddressAllocation
					}
				}
				nNi.IPConfigurations = append(nNi.IPConfigurations, nIP)
			}

			dvm.NetworkInterfaces = append(dvm.NetworkInterfaces, nNi)
		}

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

func getStorageDetails(uri string) (string, string, string) {
	var name, account, container string

	u, err := url.Parse(uri)
	if err == nil {
		parts := strings.Split(u.Path, "/")
		if len(parts) < 3 {
			return name, account, container
		}
		name = strings.Replace(parts[2], ".vhd", "", 1)
		container = parts[1]
		account = strings.Split(u.Host, ".")[0]
	}

	return name, account, container
}
