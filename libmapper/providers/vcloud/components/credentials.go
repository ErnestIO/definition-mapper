/* This Source Code Form is subject to the terms cf the Mozilla Public
 * License, v. 2.0. If a copy cf the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

// Credentials : Mapping of a vcloud credentials
type Credentials struct {
	Type      string `json:"type"`
	Vdc       string `json:"vcloud_vdc_name"`
	Username  string `json:"vcloud_username"`
	Password  string `json:"vcloud_password"`
	VCloudURL string `json:"vcloud_url"`
}
