/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

import (
	"errors"
	"unicode/utf8"
)

// ServiceDefinition ...
type ServiceDefinition struct {
	Name          string         `json:"name"`
	Datacenter    string         `json:"datacenter"`
	Bootstrapping string         `json:"bootstrapping"`
	ErnestIP      []string       `json:"ernest_ip"`
	ServiceIP     string         `json:"service_ip"`
	Routers       []Router       `json:"routers"`
	Instances     []Instance     `json:"instances"`
	Loadbalancers []Loadbalancer `json:"loadbalancers"`
}

// IsNameValid checks if service is valid
func (service *ServiceDefinition) IsNameValid() (bool, error) {
	// Check if service name is null
	if service.Name == "" {
		return false, errors.New("Service name should not be null")
	}
	// Check if service name is > 50 characters
	if utf8.RuneCountInString(service.Name) > 50 {
		return false, errors.New("Service name can't be greater than 50 characters")
	}
	return true, nil
}

// IsSaltBootstrapped : Return a boolean describing if bootstrapping option is salt
func (service *ServiceDefinition) IsSaltBootstrapped() bool {
	if service.Bootstrapping == "salt" {
		return true
	}
	return false
}
