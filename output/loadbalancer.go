/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import "net"

// Loadbalancer ...
type Loadbalancer struct {
	Name     string   `json:"name"`
	Router   string   `json:"router"`
	VIP      net.IP   `json:"vip"`
	Network  string   `json:"network"`
	Protocol string   `json:"protocol"`
	Port     string   `json:"port"`
	Servers  []Server `json:"servers"`
}
