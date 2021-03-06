/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	"github.com/ernestio/ernestprovider/types/azure/virtualmachine"
	"github.com/r3labs/graph"
)

// MapVirtualMachines ...
func MapVirtualMachines(d *definition.Definition) (vms []*components.VirtualMachine) {
	for _, rg := range d.ResourceGroups {
		for _, vm := range rg.VirtualMachines {
			image := getImageParts(vm.Image)
			if vm.Count == 0 {
				vm.Count = 1
			}

			for i := 1; i < vm.Count+1; i++ {
				cvm := &components.VirtualMachine{}
				cvm.Name = vm.Name + "-" + strconv.Itoa(i)
				cvm.VMSize = vm.Size
				cvm.AvailabilitySet = vm.AvailabilitySet

				if len(image) == 4 {
					cvm.StorageImageReference.Publisher = image[0]
					cvm.StorageImageReference.Offer = image[1]
					cvm.StorageImageReference.Sku = image[2]
					cvm.StorageImageReference.Version = image[3]
				}

				for _, ni := range vm.NetworkInterfaces {
					cvm.NetworkInterfaces = append(cvm.NetworkInterfaces, ni.Name+"-"+strconv.Itoa(i))
				}

				cvm.StorageOSDisk.Name = vm.StorageOSDisk.Name + "-" + cvm.Name
				cvm.StorageOSDisk.Caching = vm.StorageOSDisk.Caching
				cvm.StorageOSDisk.OSType = vm.StorageOSDisk.OSType
				cvm.StorageOSDisk.CreateOption = vm.StorageOSDisk.CreateOption
				cvm.StorageOSDisk.ImageURI = vm.StorageOSDisk.ImageURI
				cvm.StorageOSDisk.StorageAccountType = vm.StorageOSDisk.ManagedDiskType
				if vm.StorageOSDisk.StorageAccount != "" && vm.StorageOSDisk.StorageContainer != "" {
					cvm.StorageOSDisk.VhdURI = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s.vhd", vm.StorageOSDisk.StorageAccount, vm.StorageOSDisk.StorageContainer, vm.StorageOSDisk.Name+"-"+strconv.Itoa(i))
				}
				cvm.StorageOSDisk.StorageAccount = vm.StorageOSDisk.StorageAccount
				cvm.StorageOSDisk.StorageContainer = vm.StorageOSDisk.StorageContainer
				if vm.StorageOSDisk.ManagedDiskType != "" {
					cvm.StorageOSDisk.ManagedDisk = cvm.StorageOSDisk.Name
				}

				if vm.StorageDataDisk.Name != "" {
					cvm.StorageDataDisk.Name = vm.StorageDataDisk.Name + "-" + cvm.Name
					cvm.StorageDataDisk.Size = vm.StorageDataDisk.DiskSizeGB
					cvm.StorageDataDisk.CreateOption = vm.StorageDataDisk.CreateOption
					if vm.StorageDataDisk.StorageAccount != "" && vm.StorageDataDisk.StorageContainer != "" {
						cvm.StorageDataDisk.VhdURI = fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s.vhd", vm.StorageDataDisk.StorageAccount, vm.StorageDataDisk.StorageContainer, vm.StorageDataDisk.Name+"-"+strconv.Itoa(i))
					}
					cvm.StorageDataDisk.StorageAccount = vm.StorageDataDisk.StorageAccount
					cvm.StorageDataDisk.StorageContainer = vm.StorageDataDisk.StorageContainer
					cvm.StorageDataDisk.StorageAccountType = vm.StorageDataDisk.ManagedDiskType
					if vm.StorageDataDisk.ManagedDiskType != "" {
						cvm.StorageDataDisk.ManagedDisk = cvm.StorageDataDisk.Name
					}
				}

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

				if vm.OSProfile.ComputerName != "" {
					cvm.OSProfile.ComputerName = vm.OSProfile.ComputerName + "-" + strconv.Itoa(i)
				}
				cvm.OSProfile.CustomData = base64.StdEncoding.EncodeToString([]byte(vm.OSProfile.CustomData))

				cvm.OSProfileLinuxConfig.SSHKeys = mapSSHKeys(vm.Authentication.SSHKeys)
				cvm.OSProfileLinuxConfig.DisablePasswordAuthentication = vm.Authentication.DisablePasswordAuthentication
				cvm.OSProfile.AdminUsername = vm.Authentication.AdminUsername
				cvm.OSProfile.AdminPassword = vm.Authentication.AdminPassword

				if vm.OSProfileWindowsConfig != nil {
					config := virtualmachine.OSProfileWindowsConfig{}
					config.ProvisionVMAgent = vm.OSProfileWindowsConfig.ProvisionVMAgent
					config.EnableAutomaticUpgrades = vm.OSProfileWindowsConfig.EnableAutomaticUpgrades
					for _, winrm := range vm.OSProfileWindowsConfig.WinRM {
						config.WinRm = append(config.WinRm, virtualmachine.WinRM{
							Protocol:       winrm.Protocol,
							CertificateURL: winrm.CertificateURL,
						})
					}

					if vm.OSProfileWindowsConfig.AdditionalUnattendConfig.Pass != "" {
						config.AdditionalUnattendConfig = append(cvm.OSProfileWindowsConfig.AdditionalUnattendConfig, virtualmachine.UnattendedConfig{
							Pass:        vm.OSProfileWindowsConfig.AdditionalUnattendConfig.Pass,
							Component:   vm.OSProfileWindowsConfig.AdditionalUnattendConfig.Component,
							SettingName: vm.OSProfileWindowsConfig.AdditionalUnattendConfig.SettingName,
							Content:     vm.OSProfileWindowsConfig.AdditionalUnattendConfig.Content,
						})
					}
					cvm.OSProfileWindowsConfig = &config
				}

				tags := make(map[string]string)
				if vm.Tags != nil {
					tags = vm.Tags
				}
				cvm.Powered = true
				cvm.Tags = mapVMTags(vm.Name, d.Name, tags)
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
	ci := g.GetComponents().ByType("virtual_machine")

	for _, ig := range ci.TagValues("ernest.instance_group") {
		is := ci.ByGroup("ernest.instance_group", ig)

		if len(is) < 1 {
			continue
		}

		firstInstance := is[0].(*components.VirtualMachine)
		if firstInstance.ResourceGroupName != rg.Name {
			continue
		}

		image := firstInstance.StorageImageReference

		dvm := definition.VirtualMachine{
			Name:            ig,
			Size:            firstInstance.VMSize,
			AvailabilitySet: firstInstance.AvailabilitySet,
			Image:           strings.Join([]string{image.Publisher, image.Offer, image.Sku, image.Version}, ":"),
			Count:           len(is),
			Tags:            firstInstance.Tags,
			LicenseType:     firstInstance.LicenseType,
		}

		_, osaccount, oscontainer := getStorageDetails(firstInstance.StorageOSDisk.VhdURI)
		_, dataaccount, datacontainer := getStorageDetails(firstInstance.StorageDataDisk.VhdURI)

		nameParts := strings.Split(firstInstance.StorageOSDisk.Name, "-"+dvm.Name)
		dvm.StorageOSDisk.Name = nameParts[0]
		dvm.StorageOSDisk.Caching = firstInstance.StorageOSDisk.Caching
		dvm.StorageOSDisk.OSType = firstInstance.StorageOSDisk.OSType
		dvm.StorageOSDisk.CreateOption = firstInstance.StorageOSDisk.CreateOption
		dvm.StorageOSDisk.ImageURI = firstInstance.StorageOSDisk.ImageURI
		dvm.StorageOSDisk.StorageAccount = osaccount
		dvm.StorageOSDisk.StorageContainer = oscontainer
		dvm.StorageOSDisk.ManagedDiskType = firstInstance.StorageOSDisk.StorageAccountType

		dvm.StorageDataDisk.Name = firstInstance.StorageDataDisk.Name
		dvm.StorageDataDisk.DiskSizeGB = firstInstance.StorageDataDisk.Size
		dvm.StorageDataDisk.CreateOption = firstInstance.StorageDataDisk.CreateOption
		dvm.StorageDataDisk.StorageAccount = dataaccount
		dvm.StorageDataDisk.StorageContainer = datacontainer

		if len(firstInstance.BootDiagnostics) > 0 {
			dvm.BootDiagnostics.Enabled = firstInstance.BootDiagnostics[0].Enabled
			dvm.BootDiagnostics.StorageURI = firstInstance.BootDiagnostics[0].URI
		}

		dvm.Plan.Name = firstInstance.Plan.Name
		dvm.Plan.Product = firstInstance.Plan.Product
		dvm.Plan.Publisher = firstInstance.Plan.Publisher

		dvm.Authentication.SSHKeys = mapDefinitionSSHKeys(firstInstance.OSProfileLinuxConfig.SSHKeys)
		dvm.Authentication.DisablePasswordAuthentication = firstInstance.OSProfileLinuxConfig.DisablePasswordAuthentication
		if firstInstance.OSProfileWindowsConfig != nil {
			dvm.OSProfileWindowsConfig = &definition.OSProfileWindowsConfig{}
			dvm.OSProfileWindowsConfig.ProvisionVMAgent = firstInstance.OSProfileWindowsConfig.ProvisionVMAgent
			dvm.OSProfileWindowsConfig.EnableAutomaticUpgrades = firstInstance.OSProfileWindowsConfig.EnableAutomaticUpgrades
			for _, v := range firstInstance.OSProfileWindowsConfig.WinRm {
				winrm := definition.WinRM{}
				winrm.Protocol = v.Protocol
				winrm.CertificateURL = v.CertificateURL
				dvm.OSProfileWindowsConfig.WinRM = append(dvm.OSProfileWindowsConfig.WinRM, winrm)
			}
		}
		dvm.Authentication.AdminUsername = firstInstance.OSProfile.AdminUsername
		dvm.Authentication.AdminPassword = firstInstance.OSProfile.AdminPassword

		dvm.OSProfile = definition.OSProfile{}
		dvm.OSProfile.ComputerName = firstInstance.OSProfile.ComputerName
		dvm.OSProfile.CustomData = firstInstance.OSProfile.CustomData

		dvm.DeleteOSDiskOnTermination = firstInstance.DeleteOSDiskOnTermination
		dvm.DeleteDataDisksOnTermination = firstInstance.DeleteDataDisksOnTermination

		for _, cn := range g.GetComponents().ByType("network_interface") {
			ni := cn.(*components.NetworkInterface)

			if ni.VirtualMachineID != firstInstance.ID {
				continue
			}

			nNi := definition.NetworkInterface{
				Name:                 ig,
				EnableIPForwarding:   ni.EnableIPForwarding,
				SecurityGroup:        ni.NetworkSecurityGroup,
				DNSServers:           ni.DNSServers,
				InternalDNSNameLabel: ni.InternalDNSNameLabel,
			}

			for _, ip := range ni.IPConfigurations {
				parts := strings.Split(ip.SubnetID, "/")
				network := parts[len(parts)-3]
				nIP := definition.IPConfiguration{
					Name:                            ip.Name,
					Subnet:                          network + ":" + ip.Subnet,
					PrivateIPAddress:                ip.PrivateIPAddress,
					PrivateIPAddressAllocation:      ip.PrivateIPAddressAllocation,
					LoadBalancerBackendAddressPools: ip.LoadBalancerBackendAddressPools,
				}
				if ip.PublicIPAddressID != "" {
					cpip := g.GetComponents().ByProviderID(ip.PublicIPAddressID)
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

func mapVMTags(group, service string, tags map[string]string) map[string]string {
	if _, ok := tags["ernest.service"]; !ok {
		tags["ernest.service"] = service
	}
	tags["ernest.instance_group"] = group

	return tags
}
