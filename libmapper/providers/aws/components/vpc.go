/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"

	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// Vpc : mapping of an instance component
type Vpc struct {
	ProviderType     string            `json:"_provider" diff:"-"`
	ComponentType    string            `json:"_component" diff:"-"`
	ComponentID      string            `json:"_component_id" diff:"-"`
	State            string            `json:"_state" diff:"-"`
	Action           string            `json:"_action" diff:"-"`
	VpcAWSID         string            `json:"vpc_aws_id" diff:"-"`
	Subnet           string            `json:"subnet" diff:"-"`
	Name             string            `json:"name" diff:"name,immutable"`
	AutoRemove       bool              `json:"auto_remove" diff:"-"`
	Tags             map[string]string `json:"tags" diff:"-"`
	DatacenterType   string            `json:"datacenter_type,omitempty" diff:"-"`
	DatacenterName   string            `json:"datacenter_name,omitempty" diff:"-"`
	DatacenterRegion string            `json:"datacenter_region" diff:"-"`
	AccessKeyID      string            `json:"aws_access_key_id" diff:"-"`
	SecretAccessKey  string            `json:"aws_secret_access_key" diff:"-"`
	Service          string            `json:"service" diff:"-"`
}

// GetID : returns the component's ID
func (v *Vpc) GetID() string {
	return v.ComponentID
}

// GetName returns a components name
func (v *Vpc) GetName() string {
	return v.Name
}

// GetProvider : returns the provider type
func (v *Vpc) GetProvider() string {
	return v.ProviderType
}

// GetProviderID returns a components provider id
func (v *Vpc) GetProviderID() string {
	return v.VpcAWSID
}

// GetType : returns the type of the component
func (v *Vpc) GetType() string {
	return v.ComponentType
}

// GetState : returns the state of the component
func (v *Vpc) GetState() string {
	return v.State
}

// SetState : sets the state of the component
func (v *Vpc) SetState(s string) {
	v.State = s
}

// GetAction : returns the action of the component
func (v *Vpc) GetAction() string {
	return v.Action
}

// SetAction : Sets the action of the component
func (v *Vpc) SetAction(s string) {
	v.Action = s
}

// GetGroup : returns the components group
func (v *Vpc) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (v *Vpc) GetTags() map[string]string {
	return v.Tags
}

// GetTag returns a components tag
func (v *Vpc) GetTag(tag string) string {
	return v.Tags[tag]
}

// Diff : diff's the component against another component of the same type
func (v *Vpc) Diff(c graph.Component) (diff.Changelog, error) {
	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (v *Vpc) Update(c graph.Component) {
	v.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (v *Vpc) Rebuild(g *graph.Graph) {
	v.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (v *Vpc) Dependencies() []string {
	return []string{}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (v *Vpc) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (v *Vpc) Validate() error {
	if v.Name == "" {
		return errors.New("vpc must specify a name")
	}
	if v.Subnet == "" && v.VpcAWSID == "" {
		return errors.New("vpc must specify either subnet or an existing vpc id")
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (v *Vpc) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (v *Vpc) SetDefaultVariables() {
	v.ComponentType = TYPEVPC
	v.ComponentID = TYPEVPC + TYPEDELIMITER + v.Name
	v.ProviderType = PROVIDERTYPE
	v.DatacenterType = DATACENTERTYPE
	v.DatacenterRegion = DATACENTERREGION
	v.AccessKeyID = ACCESSKEYID
	v.SecretAccessKey = SECRETACCESSKEY
}
