/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/storagecontainer"
	"github.com/r3labs/graph"
)

// StorageContainer : A resource group a container that holds
// related resources for an Azure solution.
type StorageContainer struct {
	storagecontainer.Event
	Base
}

// GetID : returns the component's ID
func (i *StorageContainer) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *StorageContainer) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *StorageContainer) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *StorageContainer) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *StorageContainer) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *StorageContainer) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *StorageContainer) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *StorageContainer) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *StorageContainer) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *StorageContainer) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *StorageContainer) GetTags() (tags map[string]string) {
	return
}

// GetTag returns a components tag
func (i *StorageContainer) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *StorageContainer) Diff(c graph.Component) bool {
	return false
}

// Update : updates the provider returned values of a component
func (i *StorageContainer) Update(c graph.Component) {
	cs, ok := c.(*StorageContainer)
	if ok {
		i.ID = cs.ID
		i.ContainerAccessType = cs.ContainerAccessType
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *StorageContainer) Rebuild(g *graph.Graph) {
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *StorageContainer) Dependencies() (deps []string) {
	return []string{TYPESTORAGEACCOUNT + TYPEDELIMITER + i.StorageAccountName}
}

// Validate : validates the components values
func (i *StorageContainer) Validate() error {
	log.Println("Validating storage containers")
	val := event.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *StorageContainer) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *StorageContainer) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPESTORAGECONTAINER
	i.ComponentID = TYPESTORAGECONTAINER + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
