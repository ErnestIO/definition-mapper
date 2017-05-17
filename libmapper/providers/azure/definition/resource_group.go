/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// ResourceGroup ...
type ResourceGroup struct {
	ID              string            `json:"id" yaml:"id"`
	Name            string            `json:"name" yaml:"name"`
	Location        string            `json:"location" yaml:"location"`
	Tags            map[string]string `json:"tags" yaml:"tags"`
	VirtualNetworks []VirtualNetwork  `json:"virtual_networks" yaml:"virtual_networks"`
	SecurityGroups  []SecurityGroup   `json:"security_groups" yaml:"security_groups"`
	LBs             []LB              `json:"loadbalancers" yaml:"loadbalancers"`
	VirtualMachines []VirtualMachine  `json:"virtual_machines" yaml:"virtual_machines"`
	StorageAccounts []StorageAccount  `json:"storage_accounts" yaml:"storage_accounts"`
	SQLServers      []SQLServer       `json:"sql_servers" yaml:"sql_servers"`
}
