/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// InternetGateway : Mapping of a network component
type InternetGateway struct {
	ProviderType         string            `json:"_provider" diff:"-"`
	ComponentType        string            `json:"_component" diff:"-"`
	ComponentID          string            `json:"_component_id" diff:"component_id,identifier"`
	State                string            `json:"_state" diff:"-"`
	Action               string            `json:"_action" diff:"-"`
	InternetGatewayAWSID string            `json:"internet_gateway_aws_id" diff:"-"`
	Name                 string            `json:"name" diff:"-"`
	Tags                 map[string]string `json:"tags" diff:"-"`
	DatacenterType       string            `json:"datacenter_type" diff:"-"`
	DatacenterName       string            `json:"datacenter_name" diff:"-"`
	DatacenterRegion     string            `json:"datacenter_region" diff:"-"`
	AccessKeyID          string            `json:"aws_access_key_id" diff:"-"`
	SecretAccessKey      string            `json:"aws_secret_access_key" diff:"-"`
	Vpc                  string            `json:"vpc" diff:"-"`
	VpcID                string            `json:"vpc_id" diff:"-"`
	Service              string            `json:"service" diff:"-"`
}

// GetID : returns the component's ID
func (i *InternetGateway) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *InternetGateway) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *InternetGateway) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *InternetGateway) GetProviderID() string {
	return i.InternetGatewayAWSID
}

// GetType : returns the type of the component
func (i *InternetGateway) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *InternetGateway) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *InternetGateway) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *InternetGateway) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *InternetGateway) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *InternetGateway) GetGroup() string {
	return ""
}

// GetTags : returns the components tags
func (i *InternetGateway) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *InternetGateway) GetTag(tag string) string {
	return i.Tags[tag]
}

// Diff : diff's the component against another component of the same type
func (i *InternetGateway) Diff(c graph.Component) (diff.Changelog, error) {
	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (i *InternetGateway) Update(c graph.Component) {
	ci, ok := c.(*InternetGateway)
	if ok {
		i.InternetGatewayAWSID = ci.InternetGatewayAWSID
	}

	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *InternetGateway) Rebuild(g *graph.Graph) {
	if i.Vpc == "" && i.VpcID != "" {
		v := g.GetComponents().ByProviderID(i.VpcID)
		if v != nil {
			i.Vpc = v.GetName()
		}
	}

	if i.Vpc != "" && i.VpcID == "" {
		i.VpcID = templVpcID(i.Vpc)
	}

	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *InternetGateway) Dependencies() []string {
	return []string{"vpc::" + i.Vpc}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *InternetGateway) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *InternetGateway) Validate() error {
	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *InternetGateway) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *InternetGateway) SetDefaultVariables() {
	i.ComponentType = TYPEINTERNETGATEWAY
	i.ComponentID = TYPEINTERNETGATEWAY + TYPEDELIMITER + i.Name
	i.ProviderType = PROVIDERTYPE
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.AccessKeyID = ACCESSKEYID
	i.SecretAccessKey = SECRETACCESSKEY
}
