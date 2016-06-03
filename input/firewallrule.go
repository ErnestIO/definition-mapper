/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

import (
	"errors"
	"net"
	"strconv"
	"unicode/utf8"
)

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

// IsValidDestination checks if an string is a valid destionation
func (rule *FirewallRule) IsValidDestination(d string, networks []Network, loadbalancers []Loadbalancer) (bool, string) {
	valid, s := rule.IsValidSource(d, networks)

	// Check if it refers to a Loadbalancer Group
	isLoadbalancer := false
	for _, loadbalancer := range loadbalancers {
		if loadbalancer.Name == d {
			isLoadbalancer = true
			// Replace Destination with network Range
			s = loadbalancer.VIP
			break
		}
	}

	if valid == false && isLoadbalancer == false {
		return false, s
	}

	return true, s

}

// IsValidSource checks if an string is a valid source
func (rule *FirewallRule) IsValidSource(s string, networks []Network) (bool, string) {
	// Check if Source is a valid value or a valid IP/CIDR
	// One of: external | internal | any | named networks | CIDR
	if s != "external" && s != "internal" && s != "any" {

		// Check if it refers to an internal Network
		isNetwork := false
		for _, network := range networks {
			if network.Name == s {
				isNetwork = true
				// Replace Source with network Range
				s = network.Subnet
				break
			}
		}

		// Check if Source is a valid CIDR
		isValidCIDR := false
		_, _, cidrErr := net.ParseCIDR(s)
		if cidrErr == nil {
			isValidCIDR = true
		}

		// Check if Source is a valid IP
		isValidIP := false
		ip := net.ParseIP(s)
		if ip != nil {
			isValidIP = true
		}

		if !isValidIP && !isValidCIDR && !isNetwork {
			return false, ""
		}
	}

	return true, s
}

// IsValidPort checks an string to be a valid TCP port
func (rule *FirewallRule) IsValidPort(p string) (bool, bool) {
	valid := true
	inRange := true
	if p != "any" {
		port, err := strconv.Atoi(p)
		if err != nil {
			valid = false
		}
		if port < 1 || port > 65535 {
			inRange = false
		}
	}
	return valid, inRange
}

// IsValid if FirewallRule is rule
func (rule *FirewallRule) IsValid(networks []Network, loadbalancers []Loadbalancer) (bool, error) {
	// Check if firewall rule name is null
	if rule.Name == "" {
		err := errors.New("Firewall Rule name should not be null")
		return false, err
	}

	// Check if firewall rule name is > 50 characters
	if utf8.RuneCountInString(rule.Name) > 50 {
		err := errors.New("Firewall Rule name can't be greater than 50 characters")
		return false, err
	}

	valid, s := rule.IsValidSource(rule.Source, networks)
	if valid == false {
		return false, errors.New("Firewall Source is not valid")
	}
	rule.Source = s

	valid, s = rule.IsValidDestination(rule.Destination, networks, loadbalancers)
	if valid == false {
		return false, errors.New("Firewall Destination is not valid")
	}
	rule.Destination = s

	// Validate FromPort Port
	// Must be: [any | 1 - 65535]
	valid, inRange := rule.IsValidPort(rule.FromPort)
	if valid == false {
		return false, errors.New("Firewall From Port is not valid")
	}
	if inRange == false {
		return false, errors.New("Firewall From Port is out of range [1 - 65535]")
	}

	// Validate ToPort Port
	// Must be: [any | 1 - 65535]
	valid, inRange = rule.IsValidPort(rule.ToPort)
	if valid == false {
		return false, errors.New("Firewall To Port is not valid")
	}
	if inRange == false {
		return false, errors.New("Firewall To Port is out of range [1 - 65535]")
	}

	// Validate Protocol
	// Must be one of: tcp | udp | icmp | any | tcp & udp
	if rule.Protocol != "tcp" && rule.Protocol != "udp" && rule.Protocol != "icmp" && rule.Protocol != "any" && rule.Protocol != "tcp & udp" {
		return false, errors.New("Firewall Protocol is invalid")
	}

	return true, nil
}
