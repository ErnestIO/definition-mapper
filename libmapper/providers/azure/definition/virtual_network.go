/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// VirtualNetwork ...
type VirtualNetwork struct {
	Name          string   `json:"name" yaml:"name"`
	AddressSpaces []string `json:"address_spaces" yaml:"address_spaces"`
	DNSServers    []string `json:"dns_servers" yaml:"dns_servers"`
	Subnets       []Subnet `json:"subnets" yaml:"subnets"`
}
