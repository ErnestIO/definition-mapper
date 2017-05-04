/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	graph "gopkg.in/r3labs/graph.v2"
)

// PublicIP : A resource group a container that holds
// related resources for an Azure solution.
type PublicIP struct {
	ProviderType              string            `json:"_provider"`
	ComponentID               string            `json:"_component_id"`
	ComponentType             string            `json:"_component"`
	State                     string            `json:"_state"`
	Action                    string            `json:"_action"`
	DatacenterName            string            `json:"datacenter_name"`
	DatacenterType            string            `json:"datacenter_type"`
	DatacenterRegion          string            `json:"datacenter_region"`
	ID                        string            `json:"id"`
	Name                      string            `json:"name" validate:"required"`
	Location                  string            `json:"location" validate:"required"`
	ResourceGroupName         string            `json:"resource_group_name" validate:"required"`
	PublicIPAddressAllocation string            `json:"public_ip_address_allocation" validate:"required"`
	IdleTimeoutInMinutes      int               `json:"idle_timeout_in_minutes"`
	DomainNameLabel           string            `json:"domain_name_label"`
	ReverseFQDN               string            `json:"reverse_fqdn"`
	FQDN                      string            `json:"fqdn"`
	IP                        string            `json:"ip_address"`
	Tags                      map[string]string `json:"tags"`
	ClientID                  string            `json:"azure_client_id"`
	ClientSecret              string            `json:"azure_client_secret"`
	TenantID                  string            `json:"azure_tenant_id"`
	SubscriptionID            string            `json:"azure_subscription_id"`
	Environment               string            `json:"environment"`
}

// GetID : returns the component's ID
func (i *PublicIP) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *PublicIP) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *PublicIP) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *PublicIP) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *PublicIP) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *PublicIP) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *PublicIP) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *PublicIP) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *PublicIP) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *PublicIP) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *PublicIP) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *PublicIP) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *PublicIP) Diff(c graph.Component) bool {
	cs, ok := c.(*PublicIP)
	if ok {
		if i.Location != cs.Location {
			return true
		}
	}

	return false
}

// Update : updates the provider returned values of a component
func (i *PublicIP) Update(c graph.Component) {
	cs, ok := c.(*PublicIP)
	if ok {
		i.ID = cs.ID
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *PublicIP) Rebuild(g *graph.Graph) {
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *PublicIP) Dependencies() (deps []string) {
	return []string{TYPERESOURCEGROUP + TYPEDELIMITER + i.ResourceGroupName}
}

// Validate : validates the components values
func (i *PublicIP) Validate() error {
	log.Println("Validating public IP")
	val := NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *PublicIP) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *PublicIP) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPEPUBLICIP
	i.ComponentID = TYPEPUBLICIP + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
