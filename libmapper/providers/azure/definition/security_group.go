/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// SecurityGroup ...
type SecurityGroup struct {
	ID    string              `json:"id,omitempty" yaml:"id,omitempty"`
	Name  string              `json:"name,omitempty" yaml:"name,omitempty"`
	Rules []SecurityGroupRule `json:"rules,omitempty" yaml:"rules,omitempty"`
	Tags  map[string]string   `json:"tags,omitempty" yaml:"tags,omitempty"`
}
