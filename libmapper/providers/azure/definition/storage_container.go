/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// StorageContainer ...
type StorageContainer struct {
	Name       string `json:"name" yaml:"name"`
	AccessType string `json:"access_type" yaml:"access_type"`
}
