/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// Subnet ...
type Subnet struct {
	Name          string `json:"name" yaml:"name"`
	AddressPrefix string `json:"address_prefix" yaml:"address_prefix"`
	SecurityGroup string `json:"security_group" yaml:"security_group"`
}
