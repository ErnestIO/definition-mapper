/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

import (
	"errors"
	"unicode/utf8"
)

// Datacenter ...
type Datacenter struct {
	Name            string `json:"datacenter_name"`
	Username        string `json:"datacenter_username"`
	Password        string `json:"datacenter_password"`
	Type            string `json:"datacenter_type"`
	Region          string `json:"datacenter_region"`
	ExternalNetwork string `json:"external_network"`
	VCloudURL       string `json:"vcloud_url"`
	VseURL          string `json:"vse_url"`
}

// IsValid checks if a datacenter is valid
func (d *Datacenter) IsValid() (bool, error) {
	// Check if datacenter name is null
	if d.Name == "" {
		return false, errors.New("Datacenter name should not be null")
	}
	// Check if datacenter name is > 50 characters
	if utf8.RuneCountInString(d.Name) > 50 {
		return false, errors.New("Datacenter name can't be greater than 50 characters")
	}
	return true, nil
}
