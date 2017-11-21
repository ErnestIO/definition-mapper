/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// Gateway stores all information about the router and its componenets
type Gateway struct {
	Name          string         `json:"name" yaml:"name"`
	Networks      []Network      `json:"networks,omitempty" yaml:"networks,omitempty"`
	FirewallRules []FirewallRule `json:"firewall_rules,omitempty" yaml:"firewall_rules,omitempty"`
	NatRules      []NatRule      `json:"nat_rules,omitempty" yaml:"nat_rules,omitempty"`
}

// Network ...
type Network struct {
	Name   string   `json:"name" yaml:"name"`
	Subnet string   `json:"subnet" yaml:"subnet"`
	DNS    []string `json:"dns,omitempty" yaml:"dns,omitempty"`
}

// NatRule holds port forwarding information
type NatRule struct {
	Source      string `json:"source" yaml:"source"`
	Destination string `json:"destination" yaml:"destination"`
	FromPort    string `json:"from_port" yaml:"from_port"`
	ToPort      string `json:"to_port" yaml:"to_port"`
	Protocol    string `json:"protocol" yaml:"protocol"`
	Type        string `json:"type" yaml:"type"`
}

// FirewallRule ...
type FirewallRule struct {
	Name        string `json:"name" yaml:"name"`
	Source      string `json:"source" yaml:"source"`
	Destination string `json:"destination" yaml:"destination"`
	Protocol    string `json:"protocol" yaml:"protocol"`
	FromPort    string `json:"from_port" yaml:"from_port"`
	ToPort      string `json:"to_port" yaml:"to_port"`
	Action      string `json:"action,omitempty" yaml:"action.omitempty"`
}
