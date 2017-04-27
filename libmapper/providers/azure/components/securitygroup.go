/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/securitygroup"
	graph "gopkg.in/r3labs/graph.v2"
)

// SecurityGroup : A resource group a container that holds
// related resources for an Azure solution.
type SecurityGroup struct {
	securitygroup.Event
	Base
}

// GetID : returns the component's ID
func (i *SecurityGroup) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *SecurityGroup) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *SecurityGroup) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *SecurityGroup) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *SecurityGroup) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *SecurityGroup) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *SecurityGroup) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *SecurityGroup) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *SecurityGroup) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *SecurityGroup) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *SecurityGroup) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *SecurityGroup) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *SecurityGroup) Diff(c graph.Component) bool {
	cs, ok := c.(*SecurityGroup)
	if ok {
		if len(i.SecurityRules) != len(cs.SecurityRules) {
			return true
		}
		for j := range cs.SecurityRules {
			if i.SecurityRules[j].Name != cs.SecurityRules[j].Name {
				return true
			}
			if i.SecurityRules[j].Description != cs.SecurityRules[j].Description {
				return true
			}
			if i.SecurityRules[j].Protocol != cs.SecurityRules[j].Protocol {
				return true
			}
			if i.SecurityRules[j].SourcePort != cs.SecurityRules[j].SourcePort {
				return true
			}
			if i.SecurityRules[j].DestinationPortRange != cs.SecurityRules[j].DestinationPortRange {
				return true
			}
			if i.SecurityRules[j].SourceAddressPrefix != cs.SecurityRules[j].SourceAddressPrefix {
				return true
			}
			if i.SecurityRules[j].DestinationAddressPrefix != cs.SecurityRules[j].DestinationAddressPrefix {
				return true
			}
		}
	}

	return false
}

// Update : updates the provider returned values of a component
func (i *SecurityGroup) Update(c graph.Component) {
	cs, ok := c.(*SecurityGroup)
	if ok {
		i.ID = cs.ID
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *SecurityGroup) Rebuild(g *graph.Graph) {
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *SecurityGroup) Dependencies() (deps []string) {
	return []string{TYPERESOURCEGROUP + TYPEDELIMITER + i.ResourceGroupName}
}

// Validate : validates the components values
func (i *SecurityGroup) Validate() error {
	log.Println("Validating security groups")
	val := event.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *SecurityGroup) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *SecurityGroup) SetDefaultVariables() {
	i.ComponentType = TYPESECURITYGROUP
	i.ComponentID = TYPESECURITYGROUP + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
