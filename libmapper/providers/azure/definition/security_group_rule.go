/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// SecurityGroupRule ...
type SecurityGroupRule struct {
	Name                     string `json:"name,omitempty" yaml:"name,omitempty"`
	Description              string `json:"description,omitempty" yaml:"description,omitempty"`
	Priority                 int    `json:"priority,omitempty" yaml:"priority,omitempty"`
	Direction                string `json:"direction,omitempty" yaml:"direction,omitempty"`
	Access                   string `json:"access,omitempty" yaml:"access,omitempty"`
	Protocol                 string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	SourcePortRange          string `json:"source_port_range,omitempty" yaml:"source_port_range,omitempty"`
	DestinationPortRange     string `json:"destination_port_range,omitempty" yaml:"destination_port_range,omitempty"`
	SourceAddressPrefix      string `json:"source_address_prefix,omitempty" yaml:"source_address_prefix,omitempty"`
	DestinationAddressPrefix string `json:"destination_address_prefix,omitempty" yaml:"destination_address_prefix,omitempty"`
}
