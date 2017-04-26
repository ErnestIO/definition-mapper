/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/resourcegroup"
	graph "gopkg.in/r3labs/graph.v2"
)

// ResourceGroup : A resource group a container that holds
// related resources for an Azure solution.
type ResourceGroup struct {
	ID string `json:"id"`
	resourcegroup.Event
	Base
}

// GetID : returns the component's ID
func (i *ResourceGroup) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *ResourceGroup) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *ResourceGroup) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *ResourceGroup) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *ResourceGroup) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *ResourceGroup) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *ResourceGroup) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *ResourceGroup) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *ResourceGroup) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *ResourceGroup) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *ResourceGroup) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *ResourceGroup) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *ResourceGroup) Diff(c graph.Component) bool {
	cs, ok := c.(*ResourceGroup)
	if ok {
		if i.Location != cs.Location {
			return true
		}
	}
	return false
}

// Update : updates the provider returned values of a component
func (i *ResourceGroup) Update(c graph.Component) {
	cs, ok := c.(*SecurityGroup)
	if ok {
		i.ID = cs.ID
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *ResourceGroup) Rebuild(g *graph.Graph) {
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *ResourceGroup) Dependencies() (deps []string) {
	return
}

// Validate : validates the components values
func (i *ResourceGroup) Validate() error {
	val := event.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *ResourceGroup) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *ResourceGroup) SetDefaultVariables() {
	i.ComponentType = TYPERESOURCEGROUP
	i.ComponentID = TYPERESOURCEGROUP + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
