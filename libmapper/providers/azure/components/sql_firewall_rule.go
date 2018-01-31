/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/validator"
	"github.com/ernestio/ernestprovider/types/azure/sqlfirewallrule"
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// SQLFirewallRule : A resource group a container that holds
// related resources for an Azure solution.
type SQLFirewallRule struct {
	sqlfirewallrule.Event
	Base
}

// GetID : returns the component's ID
func (i *SQLFirewallRule) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *SQLFirewallRule) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *SQLFirewallRule) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *SQLFirewallRule) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *SQLFirewallRule) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *SQLFirewallRule) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *SQLFirewallRule) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *SQLFirewallRule) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *SQLFirewallRule) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *SQLFirewallRule) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *SQLFirewallRule) GetTags() (t map[string]string) {
	return
}

// GetTag returns a components tag
func (i *SQLFirewallRule) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *SQLFirewallRule) Diff(c graph.Component) (diff.Changelog, error) {
	cv, ok := c.(*SQLFirewallRule)
	if ok {
		return diff.Diff(cv, i)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (i *SQLFirewallRule) Update(c graph.Component) {
	cs, ok := c.(*SQLFirewallRule)
	if ok {
		i.ID = cs.ID
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *SQLFirewallRule) Rebuild(g *graph.Graph) {
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *SQLFirewallRule) Dependencies() (deps []string) {
	return []string{TYPESQLSERVER + TYPEDELIMITER + i.ServerName}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *SQLFirewallRule) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *SQLFirewallRule) Validate() error {
	log.Println("Validating SQL firewall rules")
	val := validator.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *SQLFirewallRule) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *SQLFirewallRule) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPESQLFIREWALLRULE
	i.ComponentID = TYPESQLFIREWALLRULE + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
