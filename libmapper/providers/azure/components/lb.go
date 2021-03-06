/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/validator"
	"github.com/ernestio/ernestprovider/types/azure/lb"
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// LB : ..
type LB struct {
	ID string `json:"id"`
	lb.Event
	Base
}

// GetID : returns the component's ID
func (i *LB) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *LB) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *LB) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *LB) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *LB) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *LB) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *LB) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *LB) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *LB) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *LB) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *LB) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *LB) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *LB) Diff(c graph.Component) (diff.Changelog, error) {
	cs, ok := c.(*LB)
	if ok {
		return diff.Diff(cs, i)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (i *LB) Update(c graph.Component) {
	cs, ok := c.(*LB)
	if ok {
		i.ID = cs.ID
		i.FrontendIPConfigurations = cs.FrontendIPConfigurations
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *LB) Rebuild(g *graph.Graph) {
	for x := 0; x < len(i.FrontendIPConfigurations); x++ {
		if i.FrontendIPConfigurations[x].PublicIPAddress == "" && i.FrontendIPConfigurations[x].PublicIPAddressID != "" {
			ip := g.GetComponents().ByProviderID(i.FrontendIPConfigurations[x].PublicIPAddressID)
			if ip != nil {
				i.FrontendIPConfigurations[x].PublicIPAddress = ip.GetName()
			}
		}

		if i.FrontendIPConfigurations[x].PublicIPAddressID == "" && i.FrontendIPConfigurations[x].PublicIPAddress != "" {
			i.FrontendIPConfigurations[x].PublicIPAddressID = templPublicIPAddressID(i.FrontendIPConfigurations[x].PublicIPAddress)
		}

		if i.FrontendIPConfigurations[x].Subnet == "" && i.FrontendIPConfigurations[x].SubnetID != "" {
			sub := g.GetComponents().ByProviderID(i.FrontendIPConfigurations[x].SubnetID)
			if sub != nil {
				i.FrontendIPConfigurations[x].Subnet = sub.GetName()
			}
		}

		if i.FrontendIPConfigurations[x].SubnetID == "" && i.FrontendIPConfigurations[x].Subnet != "" {
			i.FrontendIPConfigurations[x].SubnetID = templSubnetID(i.FrontendIPConfigurations[x].Subnet)
		}
	}

	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *LB) Dependencies() (deps []string) {
	deps = append(deps, TYPERESOURCEGROUP+TYPEDELIMITER+i.ResourceGroupName)
	for _, ip := range i.FrontendIPConfigurations {
		if ip.PublicIPAddress != "" {
			deps = append(deps, TYPEPUBLICIP+TYPEDELIMITER+ip.PublicIPAddress)
		}
		if ip.Subnet != "" {
			deps = append(deps, TYPESUBNET+TYPEDELIMITER+ip.Subnet)
		}
	}
	return
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *LB) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *LB) Validate() error {
	log.Println("Validating LB")
	val := validator.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *LB) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *LB) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPELB
	i.ComponentID = TYPELB + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
