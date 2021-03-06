/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/validator"
	"github.com/ernestio/ernestprovider/types/azure/manageddisk"
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// ManagedDisk : ..
type ManagedDisk struct {
	ID string `json:"id"`
	manageddisk.Event
	Base
}

// GetID : returns the component's ID
func (i *ManagedDisk) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *ManagedDisk) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *ManagedDisk) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *ManagedDisk) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *ManagedDisk) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *ManagedDisk) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *ManagedDisk) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *ManagedDisk) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *ManagedDisk) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *ManagedDisk) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *ManagedDisk) GetTags() map[string]string {
	return map[string]string{}
}

// GetTag returns a components tag
func (i *ManagedDisk) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *ManagedDisk) Diff(c graph.Component) (diff.Changelog, error) {
	cs, ok := c.(*ManagedDisk)
	if ok {
		return diff.Diff(cs, i)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (i *ManagedDisk) Update(c graph.Component) {
	cs, ok := c.(*ManagedDisk)
	if ok {
		i.ID = cs.ID
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *ManagedDisk) Rebuild(g *graph.Graph) {
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *ManagedDisk) Dependencies() (deps []string) {
	return []string{TYPERESOURCEGROUP + TYPEDELIMITER + i.ResourceGroupName}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *ManagedDisk) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *ManagedDisk) Validate() error {
	log.Println("Validating Managed disk")
	val := validator.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *ManagedDisk) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *ManagedDisk) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPEMANAGEDDISK
	i.ComponentID = TYPEMANAGEDDISK + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
