/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"reflect"

	"github.com/r3labs/graph"
)

// InstanceDisk an instance disk
type InstanceDisk struct {
	ID   int `json:"id"`
	Size int `json:"size"`
}

// Instance : mapping of an instance component
type Instance struct {
	ProviderType       string            `json:"_provider"`
	ComponentType      string            `json:"_component"`
	ComponentID        string            `json:"_component_id"`
	State              string            `json:"_state"`
	Action             string            `json:"_action"`
	Name               string            `json:"name"`
	Hostname           string            `json:"hostname"`
	Catalog            string            `json:"reference_catalog"`
	Image              string            `json:"reference_image"`
	Cpus               int               `json:"cpus"`
	Memory             int               `json:"ram"`
	Network            string            `json:"network"`
	IP                 string            `json:"ip"`
	Disks              []InstanceDisk    `json:"disks"`
	InstanceOnly       bool              `json:"-"`
	ShellCommands      []string          `json:"shell_commands"`
	Tags               map[string]string `json:"tags"`
	DatacenterType     string            `json:"datacenter_type"`
	DatacenterName     string            `json:"datacenter_name"`
	DatacenterUsername string            `json:"datacenter_username"`
	DatacenterPassword string            `json:"datacenter_password"`
	DatacenterRegion   string            `json:"datacenter_region"`
	VCloudURL          string            `json:"vcloud_url"`
	Service            string            `json:"service"`
}

// GetID : returns the component's ID
func (i *Instance) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *Instance) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *Instance) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *Instance) GetProviderID() string {
	return i.Name
}

// GetType : returns the type of the component
func (i *Instance) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *Instance) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *Instance) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *Instance) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *Instance) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *Instance) GetGroup() string {
	return i.Tags["ernest.instance_group"]
}

// GetTags returns a components tags
func (i *Instance) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *Instance) GetTag(tag string) string {
	return i.Tags[tag]
}

// Diff : diff's the component against another component of the same type
func (i *Instance) Diff(c graph.Component) bool {
	ci, ok := c.(*Instance)
	if ok {
		if i.Hostname != ci.Hostname {
			return true
		}

		if i.Cpus != ci.Cpus {
			return true
		}

		if i.Memory != ci.Memory {
			return true
		}

		if i.Network != ci.Network {
			return true
		}

		if i.IP != ci.IP {
			return true
		}

		if reflect.DeepEqual(i.Disks, ci.Disks) != true {
			return true
		}
	}

	return false
}

// Update : updates the provider returned values of a component
func (i *Instance) Update(c graph.Component) {
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *Instance) Rebuild(g *graph.Graph) {
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *Instance) Dependencies() []string {
	if i.InstanceOnly {
		return []string{}
	}
	return []string{TYPENETWORK + TYPEDELIMITER + i.Network}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *Instance) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *Instance) Validate() error {
	if i.Name == "" {
		return errors.New("Instance name should not be null")
	}

	if i.Image == "" {
		return errors.New("Instance image should not be null")
	}

	if i.Catalog == "" {
		return errors.New("Instance image catalog should not be null, use format <catalog>/<image>")
	}

	if i.Image == "" {
		return errors.New("Instance image image should not be null, use format <catalog>/<image>")
	}

	if i.Cpus < 1 {
		return errors.New("Instance cpus should not be < 1")
	}

	if i.Memory < 1 {
		return errors.New("Instance memory should not be null")
	}

	if i.Network == "" {
		return errors.New("Instance network name should not be null")
	}

	if i.IP == "" {
		return errors.New("Instance network start_ip should not be null")
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *Instance) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *Instance) SetDefaultVariables() {
	i.ComponentType = TYPEINSTANCE
	i.ComponentID = TYPEINSTANCE + TYPEDELIMITER + i.Name
	i.ProviderType = PROVIDERTYPE
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.DatacenterUsername = DATACENTERUSERNAME
	i.DatacenterPassword = DATACENTERPASSWORD
	i.VCloudURL = VCLOUDURL
}
