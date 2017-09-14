/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// Subnet ...
type Subnet struct {
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	AddressPrefix string `json:"address_prefix,omitempty" yaml:"address_prefix,omitempty"`
	SecurityGroup string `json:"security_group,omitempty" yaml:"security_group,omitempty"`
}
