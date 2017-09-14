/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// VirtualNetwork ...
type VirtualNetwork struct {
	Name          string   `json:"name,omitempty" yaml:"name,omitempty"`
	AddressSpaces []string `json:"address_spaces,omitempty" yaml:"address_spaces,omitempty"`
	DNSServers    []string `json:"dns_servers,omitempty" yaml:"dns_servers,omitempty"`
	Subnets       []Subnet `json:"subnets,omitempty" yaml:"subnets,omitempty"`
}
