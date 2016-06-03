/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

// Loadbalancer ...
type Loadbalancer struct {
	Name     string `json:"name"`
	Pool     string `json:"pool"`
	VIP      string `json:"vip"`
	FromPort string `json:"from_port"`
	ToPort   string `json:"to_port"`
}
