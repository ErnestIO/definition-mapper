/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// VirtualMachine ...
type VirtualMachine struct {
	Name                   string                  `json:"name,omitempty" yaml:"name,omitempty"`
	Count                  int                     `json:"count,omitempty" yaml:"count,omitempty"`
	Size                   string                  `json:"size,omitempty" yaml:"size,omitempty"`
	Image                  string                  `json:"image,omitempty" yaml:"image,omitempty"`
	AvailabilitySet        string                  `json:"availability_set,omitempty" yaml:"availability_set,omitempty"`
	Authentication         Authentication          `json:"authentication,omitempty" yaml:"authentication,omitempty"`
	StorageOSDisk          StorageOSDisk           `json:"storage_os_disk,omitempty" yaml:"storage_os_disk,omitempty"`
	OSProfile              OSProfile               `json:"os_profile,omitempty" yaml:"os_profile,omitempty"`
	OSProfileWindowsConfig *OSProfileWindowsConfig `json:"os_profile_windows_config,omitempty" yaml:"os_profile_windows_config,omitempty"`
	NetworkInterfaces      []NetworkInterface      `json:"network_interfaces,omitempty" yaml:"network_interfaces,omitempty"`
	Plan                   struct {
		Name      string `json:"name" yaml:"name,omitempty"`
		Publisher string `json:"publisher,omitempty" yaml:"publisher,omitempty"`
		Product   string `json:"product,omitempty" yaml:"product,omitempty"`
	} `json:"plan,omitempty" yaml:"plan,omitempty"`
	BootDiagnostics struct {
		Enabled    bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
		StorageURI string `json:"storage_uri,omitempty" yaml:"storage_uri,omitempty"`
	} `json:"boot_diagnostics,omitempty" yaml:"boot_diagnostics,omitempty"`
	StorageDataDisk struct {
		Name                    string `json:"name,omitempty" yaml:"name,omitempty"`
		StorageAccount          string `json:"storage_account,omitempty" yaml:"storage_account,omitempty"`
		StorageContainer        string `json:"storage_container,omitempty" yaml:"storage_container,omitempty"`
		ManagedDiskType         string `json:"managed_disk_type,omitempty" yaml:"managed_disk_type,omitempty"`
		CreateOption            string `json:"create_option,omitempty" yaml:"create_option,omitempty"`
		Caching                 string `json:"caching,omitempty" yaml:"caching,omitempty"`
		ImageURI                string `json:"image_uri,omitempty" yaml:"image_uri,omitempty"`
		StorageSourceResourceID string `json:"storage_source_resource_id,omitempty" yaml:"storage_source_resource_id,omitempty"`
		OSType                  string `json:"os_type,omitempty" yaml:"os_type,omitempty"`
		DiskSizeGB              *int32 `json:"disk_size_gb,omitempty" yaml:"disk_size_gb,omitempty"`
	} `json:"storage_data_disk,omitempty" yaml:"storage_data_disk,omitempty"`
	DeleteOSDiskOnTermination    bool              `json:"delete_os_disk_on_termination,omitempty" yaml:"delete_os_disk_on_termination,omitempty"`
	DeleteDataDisksOnTermination bool              `json:"delete_data_disks_on_termination,omitempty" yaml:"delete_data_disks_on_termination,omitempty"`
	LicenseType                  string            `json:"license_type,omitempty" yaml:"license_type,omitempty"`
	Tags                         map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Authentication ...
type Authentication struct {
	AdminUsername                 string            `json:"admin_username,omitempty" yaml:"admin_username,omitempty"`
	AdminPassword                 string            `json:"admin_password,omitempty" yaml:"admin_password,omitempty"`
	SSHKeys                       map[string]string `json:"ssh_keys,omitempty" yaml:"ssh_keys,omitempty"`
	DisablePasswordAuthentication *bool             `json:"disable_password_authentication,omitempty" yaml:"disable_password_authentication,omitempty"`
}

// StorageOSDisk ...
type StorageOSDisk struct {
	Name             string `json:"name,omitempty" yaml:"name,omitempty"`
	StorageAccount   string `json:"storage_account,omitempty" yaml:"storage_account,omitempty"`
	StorageContainer string `json:"storage_container,omitempty" yaml:"storage_container,omitempty"`
	CreateOption     string `json:"create_option,omitempty" yaml:"create_option,omitempty"`
	Caching          string `json:"caching,omitempty" yaml:"caching,omitempty"`
	ImageURI         string `json:"image_uri,omitempty" yaml:"image_uri,omitempty"`
	OSType           string `json:"os_type,omitempty" yaml:"os_type,omitempty"`
	DiskSizeGB       int32  `json:"disk_size_gb,omitempty" yaml:"disk_size_gb,omitempty"`
	ManagedDiskType  string `json:"managed_disk_type,omitempty" yaml:"managed_disk_type,omitempty"`
}

// OSProfile ...
type OSProfile struct {
	ComputerName string `json:"computer_name,omitempty" yaml:"computer_name,omitempty"`
	CustomData   string `json:"custom_data,omitempty" yaml:"custom_data,omitempty"`
}

// OSProfileWindowsConfig ...
type OSProfileWindowsConfig struct {
	ProvisionVMAgent         bool    `json:"provision_vm_agent,omitempty" yaml:"provision_vm_agent,omitempty"`
	EnableAutomaticUpgrades  bool    `json:"enable_automatic_upgrades,omitempty" yaml:"enable_automatic_upgrades,omitempty"`
	WinRM                    []WinRM `json:"winrm,omitempty" yaml:"winrm,omitempty"`
	AdditionalUnattendConfig struct {
		Pass        string `json:"pass,omitempty" yaml:"pass,omitempty"`
		Component   string `json:"component,omitempty" yaml:"component,omitempty"`
		SettingName string `json:"setting_name,omitempty" yaml:"setting_name,omitempty"`
		Content     string `json:"content,omitempty" yaml:"content,omitempty"`
	} `json:"additional_unattend_config,omitempty" yaml:"additional_unattend_config,omitempty"`
}

// WinRM ...
type WinRM struct {
	Protocol       string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	CertificateURL string `json:"certificate_url,omitempty" yaml:"certificate_url,omitempty"`
}
