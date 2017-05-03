/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	graph "gopkg.in/r3labs/graph.v2"
)

// StorageAccount : A resource group a container that holds
// related resources for an Azure solution.
type StorageAccount struct {
	ProviderType           string            `json:"_provider"`
	ComponentID            string            `json:"_component_id"`
	ComponentType          string            `json:"_component"`
	State                  string            `json:"_state"`
	Action                 string            `json:"_action"`
	DatacenterName         string            `json:"datacenter_name"`
	DatacenterType         string            `json:"datacenter_type"`
	DatacenterRegion       string            `json:"datacenter_region"`
	ID                     string            `json:"id"`
	Name                   string            `json:"name" validate:"required"`
	ResourceGroupName      string            `json:"resource_group_name" validate:"required"`
	Location               string            `json:"location" validate:"required"`
	AccountKind            string            `json:"account_kind"`
	AccountType            string            `json:"account_type" validate:"required"`
	PrimaryLocation        string            `json:"primary_location"`
	SecondaryLocation      string            `json:"secondary_location"`
	PrimaryBlobEndpoint    string            `json:"primary_blob_endpoint"`
	SecondaryBlobEndpoint  string            `json:"secondary_blob_endpoint"`
	PrimaryQueueEndpoint   string            `json:"primary_queue_endpoint"`
	SecondaryQueueEndpoint string            `json:"secondary_queue_endpoint"`
	PrimaryTableEndpoint   string            `json:"primary_table_endpoint"`
	SecondaryTableEndpoint string            `json:"secondary_table_endpoint"`
	PrimaryFileEndpoint    string            `json:"primary_file_endpoint"`
	PrimaryAccessKey       string            `json:"primary_access_key"`
	SecondaryAccessKey     string            `json:"secondary_access_key"`
	EnableBlobEncryption   bool              `json:"enable_blob_encryption"`
	Tags                   map[string]string `json:"tags"`
	ClientID               string            `json:"azure_client_id"`
	ClientSecret           string            `json:"azure_client_secret"`
	TenantID               string            `json:"azure_tenant_id"`
	SubscriptionID         string            `json:"azure_subscription_id"`
	Environment            string            `json:"environment"`
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
func (i *StorageAccount) Diff(c graph.Component) bool {
	cs, ok := c.(*StorageAccount)
	if ok {
		if i.AccountKind != cs.AccountKind {
			return true
		}
		if i.AccountType != cs.AccountType {
			return true
		}
		if i.EnableBlobEncryption != cs.EnableBlobEncryption {
			return true
		}
	}
	return false
}

// Update : updates the provider returned values of a component
func (i *StorageAccount) Update(c graph.Component) {
	cs, ok := c.(*StorageAccount)
	if ok {
		i.ID = cs.ID
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

// Validate : validates the components values
func (i *StorageAccount) Validate() error {
	log.Println("Validating storage accounts")
	val := NewValidator()
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
