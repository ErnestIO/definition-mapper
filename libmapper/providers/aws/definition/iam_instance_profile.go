/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// IamInstanceProfile ...
type IamInstanceProfile struct {
	Name  string   `json:"name" yaml:"name"`
	Path  string   `json:"path" yaml:"path"`
	Roles []string `json:"roles" yaml:"roles"`
}
