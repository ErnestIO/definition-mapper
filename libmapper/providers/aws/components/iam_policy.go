/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"

	graph "gopkg.in/r3labs/graph.v2"
)

// IamPolicy : mapping of an iam policy component
type IamPolicy struct {
	ProviderType     string `json:"_provider"`
	ComponentType    string `json:"_component"`
	ComponentID      string `json:"_component_id"`
	State            string `json:"_state"`
	Action           string `json:"_action"`
	IAMPolicyAWSID   string `json:"iam_policy_aws_id"`
	IAMPolicyARN     string `json:"iam_policy_arn"`
	Name             string `json:"name"`
	PolicyDocument   string `json:"policy_document"`
	Description      string `json:"description"`
	Path             string `json:"path"`
	DatacenterType   string `json:"datacenter_type,omitempty"`
	DatacenterName   string `json:"datacenter_name,omitempty"`
	DatacenterRegion string `json:"datacenter_region"`
	AccessKeyID      string `json:"aws_access_key_id"`
	SecretAccessKey  string `json:"aws_secret_access_key"`
	Service          string `json:"service"`
}

// GetID : returns the component's ID
func (i *IamPolicy) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *IamPolicy) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *IamPolicy) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *IamPolicy) GetProviderID() string {
	return i.IAMPolicyARN
}

// GetType : returns the type of the component
func (i *IamPolicy) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *IamPolicy) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *IamPolicy) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *IamPolicy) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *IamPolicy) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *IamPolicy) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *IamPolicy) GetTags() map[string]string {
	return map[string]string{}
}

// GetTag returns a components tag
func (i *IamPolicy) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *IamPolicy) Diff(c graph.Component) bool {
	return false
}

// Update : updates the provider returned values of a component
func (i *IamPolicy) Update(c graph.Component) {
	ci, ok := c.(*IamPolicy)
	if ok {
		i.IAMPolicyAWSID = ci.IAMPolicyAWSID
		i.IAMPolicyARN = ci.IAMPolicyARN
	}

	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *IamPolicy) Rebuild(g *graph.Graph) {
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *IamPolicy) Dependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *IamPolicy) Validate() error {
	if i.Name == "" {
		return errors.New("iam policy name should not be null")
	}

	if i.PolicyDocument == "" {
		return errors.New("iam policy document should not be null")
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *IamPolicy) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *IamPolicy) SetDefaultVariables() {
	i.ComponentType = TYPEIAMPOLICY
	i.ComponentID = TYPEIAMPOLICY + TYPEDELIMITER + i.Name
	i.ProviderType = PROVIDERTYPE
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.AccessKeyID = ACCESSKEYID
	i.SecretAccessKey = SECRETACCESSKEY
}
