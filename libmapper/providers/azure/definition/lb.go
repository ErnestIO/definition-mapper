/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// LB ...
type LB struct {
	ID                       string                    `json:"id,omitempty" yaml:"id,omitempty"`
	Name                     string                    `json:"name,omitempty" yaml:"name,omitempty"`
	Location                 string                    `json:"location,omitempty" yaml:"location,omitempty"`
	FrontendIPConfigurations []FrontendIPConfiguration `json:"frontend_ip_configurations,omitempty" yaml:"frontend_ip_configurations,omitempty" validate:"required"`
	Probes                   []LoadbalancerProbe       `json:"probes,omitempty" yaml:"probes,omitempty"`
	BackendAddressPools      []string                  `json:"backend_address_pools,omitempty" yaml:"backend_address_pools,omitempty"`
	Tags                     map[string]string         `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FrontendIPConfiguration : ..
type FrontendIPConfiguration struct {
	Name                       string             `json:"name,omitempty" validate:"required" yaml:"name,omitempty"`
	Subnet                     string             `json:"subnet,omitempty" yaml:"subnet,omitempty"`
	PublicIPAddressAllocation  string             `json:"public_ip_address_allocation,omitempty" yaml:"public_ip_address_allocation,omitempty"`
	PrivateIPAddress           string             `json:"private_ip_address,omitempty" yaml:"private_ip_address,omitempty"`
	PrivateIPAddressAllocation string             `json:"private_ip_address_allocation,omitempty" yaml:"private_ip_address_allocation,omitempty"`
	Rules                      []LoadbalancerRule `json:"rules,omitempty" yaml:"rules,omitempty"`
}

// LoadbalancerRule ...
type LoadbalancerRule struct {
	Name               string `json:"name,omitempty" yaml:"name,omitempty"`
	Protocol           string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	FrontendPort       int    `json:"frontend_port,omitempty" yaml:"frontend_port,omitempty"`
	BackendPort        int    `json:"backend_port,omitempty" yaml:"backend_port,omitempty"`
	BackendAddressPool string `json:"backend_address_pool,omitempty" yaml:"backend_address_pool,omitempty"`
	Probe              string `json:"probe,omitempty" yaml:"probe,omitempty"`
	FloatingIP         bool   `json:"floating_ip,omitempty" yaml:"floating_ip,omitempty"`
	IdleTimeout        int    `json:"idle_timeout,omitempty" yaml:"idle_timeout,omitempty"`
	LoadDistribution   string `json:"load_distribution,omitempty" yaml:"load_distribution,omitempty"`
}

// LoadbalancerProbe ...
type LoadbalancerProbe struct {
	Name            string `json:"name,omitempty" yaml:"name,omitempty"`
	Port            int    `json:"port,omitempty" yaml:"port,omitempty"`
	Protocol        string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	RequestPath     string `json:"request_path,omitempty" yaml:"request_path,omitempty"`
	Interval        int    `json:"interval,omitempty" yaml:"interval,omitempty"`
	MaximumFailures int    `json:"max_failures,omitempty" yaml:"max_failures,omitempty"`
}
