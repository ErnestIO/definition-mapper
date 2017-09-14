/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// SQLFirewallRule ..
type SQLFirewallRule struct {
	ID             string `json:"id,omitempty" yaml:"id,omitempty"`
	Name           string `json:"name,omitempty" yaml:"name,omitempty"`
	StartIPAddress string `json:"start_ip_address,omitempty" yaml:"start_ip_address,omitempty"`
	EndIPAddress   string `json:"end_ip_address,omitempty" yaml:"end_ip_address,omitempty"`
}
