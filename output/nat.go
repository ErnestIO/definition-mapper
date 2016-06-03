/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

// Nat ...
type Nat struct {
	Name       string    `json:"name"`
	Service    string    `json:"service"`
	RouterName string    `json:"router_name"`
	Rules      []NatRule `json:"rules"`
}
