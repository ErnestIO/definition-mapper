/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import "strings"

// Instance ...
type Instance struct {
	Name        string   `json:"name" yaml:"name"`
	Count       int      `json:"count" yaml:"count"`
	Cpus        int      `json:"cpus" yaml:"cpus"`
	Image       string   `json:"image" yaml:"image"`
	Memory      string   `json:"memory" yaml:"memory"`
	RootDisk    string   `json:"root_disk,omitempty" yaml:"root_disk,omitempty"`
	Disks       []string `json:"disks,omitempty" yaml:"disks,omitempty"`
	Network     string   `json:"network" yaml:"network"`
	StartIP     string   `json:"start_ip" yaml:"start_ip"`
	Provisioner []*Exec  `json:"provisioner,omitempty" yaml:"provisioner,omitempty"`
}

// Exec ...
type Exec struct {
	Shell    []string `json:"shell,omitempty"  yaml:"shell,omitempty"`
	Commands []string `json:"exec,omitempty" yaml:"exec,omitempty"`
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
