/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// AvailabilitySet ...
type AvailabilitySet struct {
	Name              string `json:"name,omitempty" yaml:"name,omitempty"`
	FaultDomainCount  int    `json:"fault_domain_count,omitempty" yaml:"fault_domain_count,omitempty"`
	UpdateDomainCount int    `json:"update_domain_count,omitempty" yaml:"update_domain_count,omitempty"`
	Managed           bool   `json:"managed,omitempty" yaml:"managed,omitempty"`
}
