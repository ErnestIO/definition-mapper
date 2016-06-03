/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

// Datacenter ...
type Datacenter struct {
	Name            string `json:"name"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Region          string `json:"region"`
	Type            string `json:"type"`
	ExternalNetwork string `json:"external_network"`
	VCloudURL       string `json:"vcloud_url"`
	VseURL          string `json:"vse_url"`
}
