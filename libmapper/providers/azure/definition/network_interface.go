/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// NetworkInterface ...
type NetworkInterface struct {
	ID                   string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name                 string            `json:"name,omitempty" yaml:"name,omitempty"`
	SecurityGroup        string            `json:"security_group,omitempty" yaml:"security_group,omitempty"`
	InternalDNSNameLabel string            `json:"internal_dns_name_label,omitempty" yaml:"internal_dns_name_label,omitempty"`
	EnableIPForwarding   bool              `json:"enable_ip_forwarding,omitempty" yaml:"enable_ip_forwarding,omitempty"`
	DNSServers           []string          `json:"dns_servers,omitempty" yaml:"dns_servers,omitempty"`
	IPConfigurations     []IPConfiguration `json:"ip_configurations,omitempty" yaml:"ip_configurations,omitempty"`
	Tags                 map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// IPConfiguration ...
type IPConfiguration struct {
	Name                            string   `json:"name,omitempty" yaml:"name,omitempty"`
	Subnet                          string   `json:"subnet,omitempty" yaml:"subnet,omitempty"`
	PublicIPAddressAllocation       string   `json:"public_ip_address_allocation,omitempty" yaml:"public_ip_address_allocation,omitempty"`
	PrivateIPAddressAllocation      string   `json:"private_ip_address_allocation,omitempty" yaml:"private_ip_address_allocation,omitempty"`
	PrivateIPAddress                string   `json:"private_ip_address,omitempty" yaml:"private_ip_address,omitempty"`
	LoadBalancerBackendAddressPools []string `json:"load_balancer_backend_address_pools,omitempty" yaml:"load_balancer_backend_address_pools,omitempty"`
}
