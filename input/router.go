/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

type Router struct {
	Name           string           `json:"name"`
	Networks       []Network        `json:"networks"`
	Rules          []FirewallRule   `json:"rules"`
	PortForwarding []PortForwarding `json:"port_forwarding"`
}
