/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// NetworkInterface ...
type NetworkInterface struct {
	Name                 string              `json:"name" yaml:"name"`
	SecurityGroup        string              `json:"security_group" yaml:"security_group"`
	InternalDNSNameLabel string              `json:"internal_dns_name_label" yaml:"internal_dns_name_label"`
	EnableIPForwarding   string              `json:"enable_ip_forwarding" yaml:"enable_ip_forwarding"`
	DNSServers           []string            `json:"dns_servers" yaml:"dns_servers"`
	IPConfigurations     []IPConfiguration   `json:"ip_configurations" yaml:"ip_configurations"`
	Tags                 []map[string]string `json:"tags" yaml:"tags"`
}

// IPConfiguration ...
type IPConfiguration struct {
	Name                       string `json:"name" yaml:"name"`
	Subnet                     string `json:"subnet" yaml:"subnet"`
	PrivateIPAddressAllocation string `json:"private_ip_address_allocation" yaml:"private_ip_address_allocation"`
	PrivateIPAddress           string `json:"private_ip_address" yaml:"private_ip_address"`
	PublicIPAddressID          string `json:"public_ip_address_id" yaml:"public_ip_address_id"`
}
