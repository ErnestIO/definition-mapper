/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

import (
	"errors"
	"net"
	"strconv"
)

// PortForwarding ...
type PortForwarding struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	FromPort    string `json:"from_port"`
	ToPort      string `json:"to_port"`
}

// IsValidPort checks an string to be a valid TCP port
func (rule *PortForwarding) IsValidPort(p string) (bool, bool) {
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

// IsValid checks if PortForwarding is valid
func (rule *PortForwarding) IsValid(nt Nat) (bool, error) {
	// Check if Destination is a valid IP
	ip := net.ParseIP(rule.Destination)
	if ip == nil {
		return false, errors.New("Port Forwarding must be a valid IP")
	}

	if rule.Source != "" {
		source := net.ParseIP(rule.Source)
		if source == nil {
			return false, errors.New("Port Forwarding source must be a valid IP")
		}
	}

	valid, inRange := rule.IsValidPort(rule.FromPort)
	if valid == false {
		return false, errors.New("Port Forwarding From Port is not valid")
	}
	if inRange == false {
		return false, errors.New("Port Forwarding From Port is out of range [1 - 65535]")
	}

	valid, inRange = rule.IsValidPort(rule.ToPort)
	if valid == false {
		return false, errors.New("Port Forwarding To Port is not valid")
	}
	if inRange == false {
		return false, errors.New("Port Forwarding To Port is out of range [1 - 65535]")
	}

	// Check External port is not already in use
	for _, nat := range nt.Rules {
		if rule.FromPort == nat.OriginPort {
			return false, errors.New("Port Forwarding From Port is already in use")
		}
		if rule.ToPort == nat.TranslationPort && rule.Destination == nat.TranslationIP {
			return false, errors.New("Port Forwarding To Port/Destination is already in use")
		}
	}

	return true, nil
}
