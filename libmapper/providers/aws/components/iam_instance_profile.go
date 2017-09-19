/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"strings"

	"github.com/r3labs/graph"
)

// IamInstanceProfile : mapping of an iam instance profile component
type IamInstanceProfile struct {
	ProviderType            string   `json:"_provider"`
	ComponentType           string   `json:"_component"`
	ComponentID             string   `json:"_component_id"`
	State                   string   `json:"_state"`
	Action                  string   `json:"_action"`
	IAMInstanceProfileAWSID string   `json:"iam_instance_profile_aws_id"`
	IAMInstanceProfileARN   string   `json:"iam_instance_profile_arn"`
	Name                    string   `json:"name"`
	Roles                   []string `json:"roles"`
	Path                    string   `json:"path"`
	DatacenterType          string   `json:"datacenter_type,omitempty"`
	DatacenterName          string   `json:"datacenter_name,omitempty"`
	DatacenterRegion        string   `json:"datacenter_region"`
	AccessKeyID             string   `json:"aws_access_key_id"`
	SecretAccessKey         string   `json:"aws_secret_access_key"`
	Remove                  bool     `json:"-"`
	Service                 string   `json:"service"`
}

// GetID : returns the component's ID
func (i *IamInstanceProfile) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *IamInstanceProfile) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *IamInstanceProfile) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *IamInstanceProfile) GetProviderID() string {
	return i.IAMInstanceProfileARN
}

// GetType : returns the type of the component
func (i *IamInstanceProfile) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *IamInstanceProfile) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *IamInstanceProfile) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *IamInstanceProfile) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *IamInstanceProfile) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *IamInstanceProfile) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *IamInstanceProfile) GetTags() map[string]string {
	return map[string]string{}
}

// GetTag returns a components tag
func (i *IamInstanceProfile) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *IamInstanceProfile) Diff(c graph.Component) bool {
	return false
}

// Update : updates the provider returned values of a component
func (i *IamInstanceProfile) Update(c graph.Component) {
	ci, ok := c.(*IamInstanceProfile)
	if ok {
		i.IAMInstanceProfileAWSID = ci.IAMInstanceProfileAWSID
		i.IAMInstanceProfileARN = ci.IAMInstanceProfileARN
	}

	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *IamInstanceProfile) Rebuild(g *graph.Graph) {
	if i.IsReferenced(g) != true && strings.Contains(g.Action, "import") {
		i.Remove = true
	}

	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *IamInstanceProfile) Dependencies() []string {
	var deps []string

	for _, role := range i.Roles {
		deps = append(deps, TYPEIAMROLE+TYPEDELIMITER+role)
	}

	return deps
}

// Validate : validates the components values
func (i *IamInstanceProfile) Validate() error {
	if i.Name == "" {
		return errors.New("iam instance profile name should not be null")
	}

	if len(i.Roles) < 1 {
		return errors.New("iam instance profile should specify at least one role")
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *IamInstanceProfile) IsStateful() bool {
	return !i.Remove
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *IamInstanceProfile) SetDefaultVariables() {
	i.ComponentType = TYPEIAMINSTANCEPROFILE
	i.ComponentID = TYPEIAMINSTANCEPROFILE + TYPEDELIMITER + i.Name
	i.ProviderType = PROVIDERTYPE
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.AccessKeyID = ACCESSKEYID
	i.SecretAccessKey = SECRETACCESSKEY
}

// IsReferenced : returns true if another component specifies this component directly
func (i *IamInstanceProfile) IsReferenced(g *graph.Graph) bool {
	var referenced []string

	for _, c := range g.GetComponents().ByType("instance") {
		instance := c.(*Instance)
		if instance.IAMInstanceProfile != nil {
			referenced = append(referenced, *instance.IAMInstanceProfile)
		}
	}

	if isOneOf(referenced, i.Name) {
		return true
	}

	return false
}
