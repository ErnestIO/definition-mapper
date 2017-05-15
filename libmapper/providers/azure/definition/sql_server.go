/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// SQLServer ...
type SQLServer struct {
	ID                         string            `json:"id" yaml:"id"`
	Name                       string            `json:"name" yaml:"name"`
	Version                    string            `json:"version" yaml:"version"`
	AdministratorLogin         string            `json:"administrator_login" yaml:"administrator_login"`
	AdministratorLoginPassword string            `json:"administrator_login_password" yaml:"administrator_login_password"`
	Tags                       map[string]string `json:"tags" yaml:"tags"`
	Databases                  []SQLDatabase     `json:"databases" yaml:"databases"`
	FirewallRules              []SQLFirewallRule `json:"rules" yaml:"rules"`
}
