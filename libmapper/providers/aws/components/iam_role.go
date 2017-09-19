/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"strings"

	"github.com/r3labs/graph"
)

// IamRole : mapping of an iam role component
type IamRole struct {
	ProviderType         string   `json:"_provider"`
	ComponentType        string   `json:"_component"`
	ComponentID          string   `json:"_component_id"`
	State                string   `json:"_state"`
	Action               string   `json:"_action"`
	IAMRoleAWSID         string   `json:"iam_role_aws_id"`
	IAMRoleARN           string   `json:"iam_role_arn"`
	Name                 string   `json:"name"`
	AssumePolicyDocument string   `json:"assume_policy_document"`
	Policies             []string `json:"policies"`
	PolicyARNs           []string `json:"policy_arns"`
	Description          string   `json:"description"`
	Path                 string   `json:"path"`
	DatacenterType       string   `json:"datacenter_type,omitempty"`
	DatacenterName       string   `json:"datacenter_name,omitempty"`
	DatacenterRegion     string   `json:"datacenter_region"`
	AccessKeyID          string   `json:"aws_access_key_id"`
	SecretAccessKey      string   `json:"aws_secret_access_key"`
	Remove               bool     `json:"-"`
	Service              string   `json:"service"`
}

// GetID : returns the component's ID
func (i *IamRole) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *IamRole) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *IamRole) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *IamRole) GetProviderID() string {
	return i.IAMRoleARN
}

// GetType : returns the type of the component
func (i *IamRole) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *IamRole) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *IamRole) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *IamRole) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *IamRole) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *IamRole) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *IamRole) GetTags() map[string]string {
	return map[string]string{}
}

// GetTag returns a components tag
func (i *IamRole) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *IamRole) Diff(c graph.Component) bool {
	return false
}

// Update : updates the provider returned values of a component
func (i *IamRole) Update(c graph.Component) {
	ci, ok := c.(*IamRole)
	if ok {
		i.IAMRoleAWSID = ci.IAMRoleAWSID
		i.IAMRoleARN = ci.IAMRoleARN
	}

	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *IamRole) Rebuild(g *graph.Graph) {
	if i.IsReferenced(g) != true && strings.Contains(g.Action, "import") {
		i.Remove = true
	}

	if len(i.Policies) > len(i.PolicyARNs) {
		for _, policy := range i.Policies {
			i.PolicyARNs = append(i.PolicyARNs, templIAMPolicyARN(policy))
		}
	}

	if len(i.PolicyARNs) > len(i.Policies) {
		for _, arn := range i.PolicyARNs {
			cp := g.GetComponents().ByProviderID(arn)
			if cp != nil {
				i.Policies = append(i.Policies, cp.GetName())
			}
		}
	}

	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *IamRole) Dependencies() []string {
	var deps []string
	for _, policy := range i.Policies {
		deps = append(deps, TYPEIAMPOLICY+TYPEDELIMITER+policy)
	}
	return deps
}

// Validate : validates the components values
func (i *IamRole) Validate() error {
	if i.Name == "" {
		return errors.New("iam role name should not be null")
	}

	if i.AssumePolicyDocument == "" {
		return errors.New("iam role must specify a policy document")
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *IamRole) IsStateful() bool {
	return !i.Remove
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *IamRole) SetDefaultVariables() {
	i.ComponentType = TYPEIAMROLE
	i.ComponentID = TYPEIAMROLE + TYPEDELIMITER + i.Name
	i.ProviderType = PROVIDERTYPE
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.AccessKeyID = ACCESSKEYID
	i.SecretAccessKey = SECRETACCESSKEY
}

// IsReferenced : returns true if another component specifies this component directly
func (i *IamRole) IsReferenced(g *graph.Graph) bool {
	var referenced []string

	for _, c := range g.GetComponents().ByType("iam_instance_profile") {
		profile := c.(*IamInstanceProfile)
		if profile.IsReferenced(g) {
			referenced = append(referenced, profile.Roles...)
		}
	}

	if isOneOf(referenced, i.Name) {
		return true
	}

	return false
}
