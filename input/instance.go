/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// Instance ...
type Instance struct {
	Count       int              `json:"count"`
	Cpus        int              `json:"cpus"`
	Image       string           `json:"image"`
	Memory      string           `json:"memory"`
	Disks       []string         `json:"disks"`
	Name        string           `json:"name"`
	Networks    InstanceNetworks `json:"networks"`
	Provisioner []Exec           `json:"provisioner"`
}

// IsValid : Validates the instance returning true or false if is valid or not
func (i *Instance) IsValid() (bool, error) {
	if i.Name == "" {
		return false, errors.New("Instance name should not be null")
	}

	if utf8.RuneCountInString(i.Name) > 50 {
		return false, errors.New("Instance name can't be greater than 50 characters")
	}

	if i.Image == "" {
		return false, errors.New("Instance image should not be null")
	}

	imageParts := strings.Split(i.Image, "/")
	if len(imageParts) < 2 {
		return false, errors.New("Instance image invalid, use format <catalog>/<image>")
	}

	if imageParts[0] == "" {
		return false, errors.New("Instance image catalog should not be null, use format <catalog>/<image>")
	}

	if imageParts[1] == "" {
		return false, errors.New("Instance image image should not be null, use format <catalog>/<image>")
	}

	if i.Cpus < 1 {
		return false, errors.New("Instance cpus should not be < 1")
	}

	if i.Memory == "" {
		return false, errors.New("Instance memory should not be null")
	}

	if i.Count < 1 {
		return false, errors.New("Instance count should not be < 1")
	}

	if i.Networks.Name == "" {
		return false, errors.New("Instance network name should not be null")
	}

	if i.Networks.StartIP == nil {
		return false, errors.New("Instance network start_ip should not be null")
	}

	return true, nil
}
