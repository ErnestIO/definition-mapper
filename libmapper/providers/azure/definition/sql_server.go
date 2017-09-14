/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// SQLServer ...
type SQLServer struct {
	ID                         string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name                       string            `json:"name,omitempty" yaml:"name,omitempty"`
	Version                    string            `json:"version,omitempty" yaml:"version,omitempty"`
	AdministratorLogin         string            `json:"administrator_login,omitempty" yaml:"administrator_login,omitempty"`
	AdministratorLoginPassword string            `json:"administrator_login_password,omitempty" yaml:"administrator_login_password,omitempty"`
	Tags                       map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Databases                  []SQLDatabase     `json:"databases,omitempty" yaml:"databases,omitempty"`
	FirewallRules              []SQLFirewallRule `json:"rules,omitempty" yaml:"rules,omitempty"`
}
