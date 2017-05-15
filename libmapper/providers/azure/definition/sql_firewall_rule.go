/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// SQLFirewallRule ..
type SQLFirewallRule struct {
	ID             string `json:"id" yaml:"id"`
	Name           string `json:"name" yaml:"name"`
	StartIPAddress string `json:"start_ip_address" yaml:"start_ip_address"`
	EndIPAddress   string `json:"end_ip_address" yaml:"end_ip_address"`
}
