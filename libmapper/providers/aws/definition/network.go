/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// Network ...
type Network struct {
	Name             string `json:"name" yaml:"name"`
	Subnet           string `json:"subnet" yaml:"subnet"`
	Public           bool   `json:"public" yaml:"public"`
	NatGateway       string `json:"nat_gateway" yaml:"nat_gateway"`
	AvailabilityZone string `json:"availability_zone" yaml:"availability_zone"`
	VPC              string `json:"vpc" yaml:"vpc"`
}
