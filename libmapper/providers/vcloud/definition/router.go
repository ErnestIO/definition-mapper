/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// Router stores all information about the router and its componenets
type Router struct {
	Name          string         `json:"name"`
	Networks      []Network      `json:"networks"`
	FirewallRules []FirewallRule `json:"firewall_rules"`
	NatRules      []NatRule      `json:"nat_rules"`
}

// Network ...
type Network struct {
	Name   string   `json:"name"`
	Router string   `json:"router"`
	Subnet string   `json:"subnet"`
	DNS    []string `json:"dns"`
}

// NatRule holds port forwarding information
type NatRule struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	FromPort    string `json:"from_port"`
	ToPort      string `json:"to_port"`
	Protocol    string `json:"protocol"`
	Type        string `json:"type"`
}

// FirewallRule ...
type FirewallRule struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Protocol    string `json:"protocol"`
	FromPort    string `json:"from_port"`
	ToPort      string `json:"to_port"`
	Action      string `json:"action"`
}
