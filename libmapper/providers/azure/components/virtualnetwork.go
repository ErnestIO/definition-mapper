/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"
	"reflect"

	"gopkg.in/r3labs/graph.v2"
)

// VirtualNetwork : A resource group a container that holds
// related resources for an Azure solution.
type VirtualNetwork struct {
	ProviderType      string            `json:"_provider"`
	ComponentID       string            `json:"_component_id"`
	ComponentType     string            `json:"_component"`
	State             string            `json:"_state"`
	Action            string            `json:"_action"`
	DatacenterName    string            `json:"datacenter_name"`
	DatacenterType    string            `json:"datacenter_type"`
	DatacenterRegion  string            `json:"datacenter_region"`
	ID                string            `json:"id"`
	Name              string            `json:"name" validate:"required"`
	AddressSpace      []string          `json:"address_space" validate:"min=1"`
	DNSServerNames    []string          `json:"dns_server_names" validate:"dive,ip"`
	Subnets           []VNSubnet        `json:"subnets" validate:"min=1"`
	Location          string            `json:"location"`
	ResourceGroupName string            `json:"resource_group_name"`
	Tags              map[string]string `json:"tags"`
	ClientID          string            `json:"azure_client_id"`
	ClientSecret      string            `json:"azure_client_secret"`
	TenantID          string            `json:"azure_tenant_id"`
	SubscriptionID    string            `json:"azure_subscription_id"`
	Environment       string            `json:"environment"`
}

// VNSubnet ..
type VNSubnet struct {
	Name          string `json:"name"`
	AddressPrefix string `json:"address_prefix"`
	SecurityGroup string `json:"security_group"`
}

// GetID : returns the component's ID
func (vn *VirtualNetwork) GetID() string {
	return vn.ComponentID
}

// GetName returns a components name
func (vn *VirtualNetwork) GetName() string {
	return vn.Name
}

// GetProvider : returns the provider type
func (vn *VirtualNetwork) GetProvider() string {
	return vn.ProviderType
}

// GetProviderID returns a components provider id
func (vn *VirtualNetwork) GetProviderID() string {
	return vn.ID
}

// GetType : returns the type of the component
func (vn *VirtualNetwork) GetType() string {
	return vn.ComponentType
}

// GetState : returns the state of the component
func (vn *VirtualNetwork) GetState() string {
	return vn.State
}

// SetState : sets the state of the component
func (vn *VirtualNetwork) SetState(s string) {
	vn.State = s
}

// GetAction : returns the action of the component
func (vn *VirtualNetwork) GetAction() string {
	return vn.Action
}

// SetAction : Sets the action of the component
func (vn *VirtualNetwork) SetAction(s string) {
	vn.Action = s
}

// GetGroup : returns the components group
func (vn *VirtualNetwork) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (vn *VirtualNetwork) GetTags() map[string]string {
	return vn.Tags
}

// GetTag returns a components tag
func (vn *VirtualNetwork) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (vn *VirtualNetwork) Diff(c graph.Component) bool {
	cvn, ok := c.(*VirtualNetwork)
	if ok {
		return !reflect.DeepEqual(vn.DNSServerNames, cvn.DNSServerNames)
	}

	return false
}

// Update : updates the provider returned values of a component
func (vn *VirtualNetwork) Update(c graph.Component) {
	cvn, ok := c.(*VirtualNetwork)
	if ok {
		vn.ID = cvn.ID
	}

	vn.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (vn *VirtualNetwork) Rebuild(g *graph.Graph) {
	vn.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (vn *VirtualNetwork) Dependencies() (deps []string) {
	return []string{TYPERESOURCEGROUP + TYPEDELIMITER + vn.ResourceGroupName}
}

// Validate : validates the components values
func (vn *VirtualNetwork) Validate() error {
	log.Println("Validating Virtual network")
	val := NewValidator()
	return val.Validate(vn)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (vn *VirtualNetwork) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (vn *VirtualNetwork) SetDefaultVariables() {
	vn.ProviderType = PROVIDERTYPE
	vn.ComponentType = TYPEVIRTUALNETWORK
	vn.ComponentID = TYPEVIRTUALNETWORK + TYPEDELIMITER + vn.Name
	vn.DatacenterName = DATACENTERNAME
	vn.DatacenterType = DATACENTERTYPE
	vn.DatacenterRegion = DATACENTERREGION
	vn.ClientID = CLIENTID
	vn.ClientSecret = CLIENTSECRET
	vn.TenantID = TENANTID
	vn.SubscriptionID = SUBSCRIPTIONID
	vn.Environment = ENVIRONMENT
}
