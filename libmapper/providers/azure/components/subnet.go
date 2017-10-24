/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/subnet"
	"github.com/r3labs/graph"
)

// Subnet : A resource group a container that holds
// related resources for an Azure solution.
type Subnet struct {
	subnet.Event
	Base
}

// GetID : returns the component's ID
func (s *Subnet) GetID() string {
	return s.ComponentID
}

// GetName returns a components name
func (s *Subnet) GetName() string {
	return s.Name
}

// GetProvider : returns the provider type
func (s *Subnet) GetProvider() string {
	return s.ProviderType
}

// GetProviderID returns a components provider id
func (s *Subnet) GetProviderID() string {
	return s.ID
}

// GetType : returns the type of the component
func (s *Subnet) GetType() string {
	return s.ComponentType
}

// GetState : returns the state of the component
func (s *Subnet) GetState() string {
	return s.State
}

// SetState : sets the state of the component
func (s *Subnet) SetState(state string) {
	s.State = state
}

// GetAction : returns the action of the component
func (s *Subnet) GetAction() string {
	return s.Action
}

// SetAction : Sets the action of the component
func (s *Subnet) SetAction(action string) {
	s.Action = action
}

// GetGroup : returns the components group
func (s *Subnet) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (s *Subnet) GetTags() (tags map[string]string) {
	return
}

// GetTag returns a components tag
func (s *Subnet) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (s *Subnet) Diff(c graph.Component) bool {
	cs, ok := c.(*Subnet)
	if ok {
		if s.NetworkSecurityGroup != cs.NetworkSecurityGroup {
			return true
		}
	}

	return false
}

// Update : updates the provider returned values of a component
func (s *Subnet) Update(c graph.Component) {
	cs, ok := c.(*Subnet)
	if ok {
		s.ID = cs.ID
		s.IPConfigurations = cs.IPConfigurations
		s.RouteTable = cs.RouteTable
	}
	s.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (s *Subnet) Rebuild(g *graph.Graph) {
	if s.NetworkSecurityGroup == "" && s.NetworkSecurityGroupID != "" {
		sg := g.GetComponents().ByProviderID(s.NetworkSecurityGroupID)
		if sg != nil {
			s.NetworkSecurityGroup = sg.GetName()
		}
	}

	if s.NetworkSecurityGroupID == "" && s.NetworkSecurityGroup != "" {
		s.NetworkSecurityGroupID = templSecurityGroupID(s.NetworkSecurityGroup)
	}

	s.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (s *Subnet) Dependencies() (deps []string) {
	if s.NetworkSecurityGroup != "" {
		deps = append(deps, TYPESECURITYGROUP+TYPEDELIMITER+s.NetworkSecurityGroup)
	}

	return append(deps, TYPEVIRTUALNETWORK+TYPEDELIMITER+s.VirtualNetworkName)
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (s *Subnet) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (s *Subnet) Validate() error {
	log.Println("Validating subnets")
	val := event.NewValidator()
	return val.Validate(s)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (s *Subnet) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (s *Subnet) SetDefaultVariables() {
	s.ProviderType = PROVIDERTYPE
	s.ComponentType = TYPESUBNET
	s.ComponentID = TYPESUBNET + TYPEDELIMITER + s.Name
	s.DatacenterName = DATACENTERNAME
	s.DatacenterType = DATACENTERTYPE
	s.DatacenterRegion = DATACENTERREGION
	s.ClientID = CLIENTID
	s.ClientSecret = CLIENTSECRET
	s.TenantID = TENANTID
	s.SubscriptionID = SUBSCRIPTIONID
	s.Environment = ENVIRONMENT
}
