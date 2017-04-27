/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// PublicIP ...
type PublicIP struct {
	ID                        string              `json:"id" yaml:"id"`
	Name                      string              `json:"name" yaml:"name"`
	Location                  string              `json:"location" yaml:"location"`
	PublicIPAddressAllocation string              `json:"public_ip_address_allocation" yaml:"public_ip_address_allocation"`
	Tags                      []map[string]string `json:"tags" yaml:"tags"`
}
