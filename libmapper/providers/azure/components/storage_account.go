/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/validator"
	"github.com/ernestio/ernestprovider/types/azure/storageaccount"
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// StorageAccount : A resource group a container that holds
// related resources for an Azure solution.
type StorageAccount struct {
	storageaccount.Event
	Base
}

// GetID : returns the component's ID
func (i *StorageAccount) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *StorageAccount) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *StorageAccount) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *StorageAccount) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *StorageAccount) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *StorageAccount) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *StorageAccount) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *StorageAccount) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *StorageAccount) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *StorageAccount) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *StorageAccount) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *StorageAccount) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *StorageAccount) Diff(c graph.Component) (diff.Changelog, error) {
	cs, ok := c.(*StorageAccount)
	if ok {
		return diff.Diff(cs, i)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (i *StorageAccount) Update(c graph.Component) {
	cs, ok := c.(*StorageAccount)
	if ok {
		i.ID = cs.ID
		i.AccountKind = cs.AccountKind
		i.AccountType = cs.AccountType
		i.EnableBlobEncryption = cs.EnableBlobEncryption
		i.PrimaryLocation = cs.PrimaryLocation
		i.SecondaryLocation = cs.SecondaryLocation
		i.PrimaryBlobEndpoint = cs.PrimaryBlobEndpoint
		i.SecondaryBlobEndpoint = cs.SecondaryBlobEndpoint
		i.PrimaryQueueEndpoint = cs.PrimaryQueueEndpoint
		i.SecondaryQueueEndpoint = cs.SecondaryQueueEndpoint
		i.PrimaryTableEndpoint = cs.PrimaryTableEndpoint
		i.SecondaryTableEndpoint = cs.SecondaryTableEndpoint
		i.PrimaryFileEndpoint = cs.PrimaryFileEndpoint
		i.PrimaryAccessKey = cs.PrimaryAccessKey
		i.SecondaryAccessKey = cs.SecondaryAccessKey
		i.EnableBlobEncryption = cs.EnableBlobEncryption
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *StorageAccount) Rebuild(g *graph.Graph) {
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *StorageAccount) Dependencies() (deps []string) {
	return []string{TYPERESOURCEGROUP + TYPEDELIMITER + i.ResourceGroupName}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *StorageAccount) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *StorageAccount) Validate() error {
	log.Println("Validating storage accounts")
	val := validator.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *StorageAccount) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *StorageAccount) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPESTORAGEACCOUNT
	i.ComponentID = TYPESTORAGEACCOUNT + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
