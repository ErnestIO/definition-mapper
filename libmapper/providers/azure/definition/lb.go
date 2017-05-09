/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// LB ...
type LB struct {
	ID                       string                    `json:"id" yaml:"id"`
	Name                     string                    `json:"name" yaml:"name"`
	Location                 string                    `json:"location" yaml:"location"`
	FrontendIPConfigurations []FrontendIPConfiguration `json:"frontend_ip_configurations" validate:"required"`
	Tags                     []map[string]string       `json:"tags" yaml:"tags"`
}

// FrontendIPConfiguration : ..
type FrontendIPConfiguration struct {
	Name                       string `json:"name" validate:"required" yaml:"name"`
	Subnet                     string `json:"subnet" yaml:"subnet"`
	PublicIPAddressAllocation  string `json:"public_ip_address_allocation" yaml:"public_ip_address_allocation"`
	PrivateIPAddress           string `json:"private_ip_address" yaml:"private_ip_address"`
	PrivateIPAddressAllocation string `json:"private_ip_address_allocation" yaml:"private_ip_address_allocation"`
}
