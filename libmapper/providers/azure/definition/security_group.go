/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// SecurityGroup ...
type SecurityGroup struct {
	ID    string              `json:"id" yaml:"id"`
	Name  string              `json:"name" yaml:"name"`
	Rules []SecurityGroupRule `json:"rules" yaml:"rules"`
	Tags  map[string]string   `json:"tags" yaml:"tags"`
}
