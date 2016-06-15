/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

// Client ...
type Client struct {
	ID   string `json:"client_id"`
	Name string `json:"client_name"`
}