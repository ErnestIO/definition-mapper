/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// SecurityGroupRule ...
type SecurityGroupRule struct {
	Name                     string `json:"name" yaml:"name"`
	Description              string `json:"description" yaml:"description"`
	Priority                 int    `json:"priority" yaml:"priority"`
	Direction                string `json:"direction" yaml:"direction"`
	Access                   string `json:"access" yaml:"access"`
	Protocol                 string `json:"protocol" yaml:"protocol"`
	SourcePortRange          string `json:"source_port_range" yaml:"source_port_range"`
	DestinationPortRange     string `json:"destination_port_range" yaml:"destination_port_range"`
	SourceAddressPrefix      string `json:"source_address_prefix" yaml:"source_address_prefix"`
	DestinationAddressPrefix string `json:"destination_address_prefix" yaml:"destination_address_prefix"`
}
