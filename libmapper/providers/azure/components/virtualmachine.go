/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"reflect"

	graph "gopkg.in/r3labs/graph.v2"
)

// VirtualMachine : A resource group a container that holds
// related resources for an Azure solution.
type VirtualMachine struct {
	ProviderType      string `json:"_provider"`
	ComponentID       string `json:"_component_id"`
	ComponentType     string `json:"_component"`
	State             string `json:"_state"`
	Action            string `json:"_action"`
	DatacenterName    string `json:"datacenter_name"`
	DatacenterType    string `json:"datacenter_type"`
	DatacenterRegion  string `json:"datacenter_region"`
	ID                string `json:"id"`
	Name              string `json:"name" validate:"required"`
	ResourceGroupName string `json:"resource_group_name" validate:"required"`
	Location          string `json:"location" validate:"required"`
	Plan              struct {
		Name      string `json:"name"`
		Publisher string `json:"publisher"`
		Product   string `json:"product"`
	} `json:"plan"`
	AvailabilitySetID     string `json:"availability_set_id"`
	LicenseType           string `json:"license_type"`
	VMSize                string `json:"vm_size"`
	StorageImageReference struct {
		Publisher string `json:"publisher" structs:"publisher"`
		Offer     string `json:"offer" structs:"offer"`
		Sku       string `json:"sku" structs:"sku"`
		Version   string `json:"version" structs:"version"`
	} `json:"storage_image_reference" validate:"dive"`
	StorageOSDisk struct {
		Name         string `json:"name" structs:"name"`
		VhdURI       string `json:"vhd_uri" structs:"vhd_uri"`
		CreateOption string `json:"create_option" structs:"create_option"`
		OSType       string `json:"os_type" structs:"os_type"`
		ImageURI     string `json:"image_uri" structs:"image_uri"`
		Caching      string `json:"caching" structs:"caching"`
	} `json:"storage_os_disk" validate:"dive"`
	DeleteOSDiskOnTermination bool `json:"delete_os_disk_on_termination"`
	StorageDataDisk           struct {
		Name         string `json:"name" structs:"name"`
		VhdURI       string `json:"vhd_uri" structs:"vhd_uri"`
		CreateOption string `json:"create_option" structs:"create_option"`
		Size         int32  `json:"disk_size_gb" structs:"disk_size_gb"`
		Lun          int32  `json:"lun" structs:"lun"`
	} `json:"storage_data_disk"`
	DeleteDataDisksOnTermination bool             `json:"delete_data_disks_on_termination"`
	BootDiagnostics              []BootDiagnostic `json:"boot_diagnostics"`
	OSProfile                    struct {
		ComputerName  string `json:"computer_name" structs:"computer_name"`
		AdminUsername string `json:"admin_username" structs:"admin_username"`
		AdminPassword string `json:"admin_password" structs:"admin_password"`
		CustomData    string `json:"custom_data" structs:"custom_data"`
	} `json:"os_profile"`
	OSProfileWindowsConfig struct {
		ProvisionVMAgent         bool               `json:"provision_vm_agent" structs:"provision_vm_agent"`
		EnableAutomaticUpgrades  bool               `json:"enable_automatic_upgrades" structs:"enable_automatic_upgrades"`
		WinRm                    []WinRM            `json:"winrm" structs:"winrm"`
		AdditionalUnattendConfig []UnattendedConfig `json:"additional_unattend_config" structs:"additional_unattend_config"`
	} `json:"os_profile_windows_config"`
	OSProfileLinuxConfig struct {
		DisablePasswordAuthentication bool     `json:"disable_password_authentication" structs:"disable_password_authentication"`
		SSHKeys                       []SSHKey `json:"ssh_keys" structs:"ssh_keys"`
	} `json:"os_profile_linux_config"`
	NetworkInterfaces   []string          `json:"network_interfaces"`
	NetworkInterfaceIDs []string          `json:"network_interface_ids"`
	Tags                map[string]string `json:"tags"`
	ClientID            string            `json:"azure_client_id"`
	ClientSecret        string            `json:"azure_client_secret"`
	TenantID            string            `json:"azure_tenant_id"`
	SubscriptionID      string            `json:"azure_subscription_id"`
	Environment         string            `json:"environment"`
}

// WinRM ...
type WinRM struct {
	Protocol       string `json:"protocol" structs:"protocol"`
	CertificateURL string `json:"certificate_url" structs:"certification_url"`
}

// SSHKey ...
type SSHKey struct {
	Path    string `json:"path" structs:"path"`
	KeyData string `json:"key_data" structs:"key_data"`
}

// BootDiagnostic ...
type BootDiagnostic struct {
	Enabled bool   `json:"enabled"`
	URI     string `json:"storage_uri"`
}

// UnattendedConfig ...
type UnattendedConfig struct {
	Pass        string `json:"pass" structs:"pass"`
	Component   string `json:"component" structs:"component"`
	SettingName string `json:"setting_name" structs:"setting_name"`
	Content     string `json:"content" structs:"content"`
}

// GetID : returns the component's ID
func (i *VirtualMachine) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *VirtualMachine) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *VirtualMachine) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *VirtualMachine) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *VirtualMachine) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *VirtualMachine) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *VirtualMachine) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *VirtualMachine) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *VirtualMachine) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *VirtualMachine) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *VirtualMachine) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *VirtualMachine) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *VirtualMachine) Diff(c graph.Component) bool {
	cvm, ok := c.(*VirtualMachine)
	if ok {
		if i.VMSize != cvm.VMSize {
			return true
		}

		if i.StorageDataDisk.Size != cvm.StorageDataDisk.Size {
			return true
		}

		if reflect.DeepEqual(i.NetworkInterfaces, cvm.NetworkInterfaces) != true {
			return true
		}
	}

	return false
}

// Update : updates the provider returned values of a component
func (i *VirtualMachine) Update(c graph.Component) {
	cvm, ok := c.(*VirtualMachine)
	if ok {
		i.ID = cvm.ID
		// ???
		i.StorageDataDisk.Lun = cvm.StorageDataDisk.Lun
	}

	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *VirtualMachine) Rebuild(g *graph.Graph) {
	if len(i.NetworkInterfaces) > len(i.NetworkInterfaceIDs) {
		for _, iface := range i.NetworkInterfaces {
			i.NetworkInterfaceIDs = append(i.NetworkInterfaceIDs, templNetworkInterfaceID(iface))
		}
	}

	if len(i.NetworkInterfaceIDs) > len(i.NetworkInterfaces) {
		for _, id := range i.NetworkInterfaceIDs {
			iface := g.GetComponents().ByProviderID(id)
			if iface != nil {
				i.NetworkInterfaces = append(i.NetworkInterfaces, iface.GetName())
			}
		}
	}

	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *VirtualMachine) Dependencies() (deps []string) {
	for _, iface := range i.NetworkInterfaces {
		deps = append(deps, TYPENETWORKINTERFACE+TYPEDELIMITER+iface)
	}

	if len(deps) < 1 {
		return []string{TYPERESOURCEGROUP + TYPEDELIMITER + i.ResourceGroupName}
	}

	return
}

// Validate : validates the components values
func (i *VirtualMachine) Validate() error {
	val := NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *VirtualMachine) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *VirtualMachine) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPEVIRTUALMACHINE
	i.ComponentID = TYPEVIRTUALMACHINE + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
