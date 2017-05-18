/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// LB ...
type LB struct {
	ID                       string                    `json:"id" yaml:"id"`
	Name                     string                    `json:"name" yaml:"name"`
	Location                 string                    `json:"location" yaml:"location"`
	FrontendIPConfigurations []FrontendIPConfiguration `json:"frontend_ip_configurations" yaml:"frontend_ip_configurations" validate:"required"`
	Probes                   []LoadbalancerProbe       `json:"probes" yaml:"probes"`
	BackendAddressPools      []string                  `json:"backend_address_pools" yaml:"backend_address_pools"`
	Tags                     []map[string]string       `json:"tags" yaml:"tags"`
}

// FrontendIPConfiguration : ..
type FrontendIPConfiguration struct {
	Name                       string             `json:"name" validate:"required" yaml:"name"`
	Subnet                     string             `json:"subnet" yaml:"subnet"`
	PublicIPAddressAllocation  string             `json:"public_ip_address_allocation" yaml:"public_ip_address_allocation"`
	PrivateIPAddress           string             `json:"private_ip_address" yaml:"private_ip_address"`
	PrivateIPAddressAllocation string             `json:"private_ip_address_allocation" yaml:"private_ip_address_allocation"`
	Rules                      []LoadbalancerRule `json:"rules" yaml:"rules"`
}

// LoadbalancerRule ...
type LoadbalancerRule struct {
	Name               string `json:"name" yaml:"name"`
	Protocol           string `json:"protocol" yaml:"protocol"`
	FrontendPort       int    `json:"frontend_port" yaml:"frontend_port"`
	BackendPort        int    `json:"backend_port" yaml:"backend_port"`
	BackendAddressPool string `json:"backend_address_pool" yaml:"backend_address_pool"`
	Probe              string `json:"probe" yaml:"probe"`
	FloatingIP         bool   `json:"floating_ip" yaml:"floating_ip"`
	IdleTimeout        int    `json:"idle_timeout" yaml:"idle_timeout"`
	LoadDistribution   string `json:"load_distribution" yaml:"load_distribution"`
}

// LoadbalancerProbe ...
type LoadbalancerProbe struct {
	Name            string `json:"name" yaml:"name"`
	Port            int    `json:"port" yaml:"port"`
	Protocol        string `json:"protocol" yaml:"protocol"`
	RequestPath     string `json:"request_path" yaml:"request_path"`
	Interval        int    `json:"interval" yaml:"interval"`
	MaximumFailures int    `json:"max_failures" yaml:"max_failures"`
}
