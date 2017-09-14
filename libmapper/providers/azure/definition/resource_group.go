/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// ResourceGroup ...
type ResourceGroup struct {
	ID               string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name             string            `json:"name,omitempty" yaml:"name,omitempty"`
	Location         string            `json:"location,omitempty" yaml:"location,omitempty"`
	Tags             map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
	VirtualNetworks  []VirtualNetwork  `json:"virtual_networks,omitempty" yaml:"virtual_networks,omitempty"`
	SecurityGroups   []SecurityGroup   `json:"security_groups,omitempty" yaml:"security_groups,omitempty"`
	LBs              []LB              `json:"loadbalancers,omitempty" yaml:"loadbalancers,omitempty"`
	VirtualMachines  []VirtualMachine  `json:"virtual_machines,omitempty" yaml:"virtual_machines,omitempty"`
	AvailabilitySets []AvailabilitySet `json:"availability_sets,omitempty" yaml:"availability_sets,omitempty"`
	StorageAccounts  []StorageAccount  `json:"storage_accounts,omitempty" yaml:"storage_accounts,omitempty"`
	SQLServers       []SQLServer       `json:"sql_servers,omitempty" yaml:"sql_servers,omitempty"`
}
