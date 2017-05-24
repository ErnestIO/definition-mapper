/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// AvailabilitySet ...
type AvailabilitySet struct {
	Name              string `json:"name" yaml:"name"`
	FaultDomainCount  int    `json:"fault_domain_count" yaml:"fault_domain_count"`
	UpdateDomainCount int    `json:"update_domain_count" yaml:"update_domain_count"`
}
