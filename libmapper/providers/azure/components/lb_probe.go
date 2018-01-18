/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/lbprobe"
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// LBProbe : ..
type LBProbe struct {
	ID string `json:"id"`
	lbprobe.Event
	Base
}

// GetID : returns the component's ID
func (i *LBProbe) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *LBProbe) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *LBProbe) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *LBProbe) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *LBProbe) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *LBProbe) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *LBProbe) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *LBProbe) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *LBProbe) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *LBProbe) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *LBProbe) GetTags() map[string]string {
	return map[string]string{}
}

// GetTag returns a components tag
func (i *LBProbe) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *LBProbe) Diff(c graph.Component) (diff.Changelog, error) {
	cs, ok := c.(*LBProbe)
	if ok {
		return diff.Diff(cs, i)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (i *LBProbe) Update(c graph.Component) {
	cs, ok := c.(*LBProbe)
	if ok {
		i.ID = cs.ID
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *LBProbe) Rebuild(g *graph.Graph) {
	if i.Loadbalancer == "" && i.LoadbalancerID != "" {
		lb := g.GetComponents().ByProviderID(i.LoadbalancerID)
		if lb != nil {
			i.Loadbalancer = lb.GetName()
		}
	}

	if i.LoadbalancerID == "" && i.Loadbalancer != "" {
		i.LoadbalancerID = templLoadbalancerID(i.Loadbalancer)
	}

	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *LBProbe) Dependencies() (deps []string) {
	return append(deps, TYPELB+TYPEDELIMITER+i.Loadbalancer)
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *LBProbe) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *LBProbe) Validate() error {
	log.Println("Validating LB")
	val := event.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *LBProbe) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *LBProbe) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPELBPROBE
	i.ComponentID = TYPELBPROBE + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
