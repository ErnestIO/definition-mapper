/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

import (
	"errors"
	"net"
	"unicode/utf8"
)

// Network ...
type Network struct {
	Name   string   `json:"name"`
	Router string   `json:"router"`
	Subnet string   `json:"subnet"`
	DNS    []string `json:"dns"`
}

// IsValid checks if a Network is valid
func (n *Network) IsValid() (bool, error) {
	_, _, err := net.ParseCIDR(n.Subnet)
	if err != nil {
		return false, errors.New("Network CIDR is not valid")
	}

	if n.Name == "" {
		return false, errors.New("Network name should not be null")
	}

	for _, val := range n.DNS {
		if ok := net.ParseIP(val); ok == nil {
			return false, errors.New("DNS " + val + " is not a valid CIDR")
		}
	}

	// Check if network name is > 50 characters
	if utf8.RuneCountInString(n.Name) > 50 {
		return false, errors.New("Network name can't be greater than 50 characters")
	}

	return true, nil
}
