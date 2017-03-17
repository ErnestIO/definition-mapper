/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import "strings"

// Instance ...
type Instance struct {
	Count       int      `json:"count"`
	Cpus        int      `json:"cpus"`
	Image       string   `json:"image"`
	Memory      string   `json:"memory"`
	RootDisk    string   `json:"root_disk"`
	Disks       []string `json:"disks"`
	Name        string   `json:"name"`
	Network     string   `json:"network"`
	StartIP     string   `json:"start_ip"`
	Provisioner []Exec   `json:"provisioner"`
}

// Exec ...
type Exec struct {
	Shell    []string `json:"shell"`
	Commands []string `json:"exec"`
}

// Catalog ...
func (i *Instance) Catalog() string {
	parts := strings.Split(i.Image, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

// Template ...
func (i *Instance) Template() string {
	parts := strings.Split(i.Image, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
