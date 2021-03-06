/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/validator"
	"github.com/ernestio/ernestprovider/types/azure/virtualnetwork"
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// VirtualNetwork : A resource group a container that holds
// related resources for an Azure solution.
type VirtualNetwork struct {
	virtualnetwork.Event
	Base
}

// GetID : returns the component's ID
func (vn *VirtualNetwork) GetID() string {
	return vn.ComponentID
}

// GetName returns a components name
func (vn *VirtualNetwork) GetName() string {
	return vn.Name
}

// GetProvider : returns the provider type
func (vn *VirtualNetwork) GetProvider() string {
	return vn.ProviderType
}

// GetProviderID returns a components provider id
func (vn *VirtualNetwork) GetProviderID() string {
	return vn.ID
}

// GetType : returns the type of the component
func (vn *VirtualNetwork) GetType() string {
	return vn.ComponentType
}

// GetState : returns the state of the component
func (vn *VirtualNetwork) GetState() string {
	return vn.State
}

// SetState : sets the state of the component
func (vn *VirtualNetwork) SetState(s string) {
	vn.State = s
}

// GetAction : returns the action of the component
func (vn *VirtualNetwork) GetAction() string {
	return vn.Action
}

// SetAction : Sets the action of the component
func (vn *VirtualNetwork) SetAction(s string) {
	vn.Action = s
}

// GetGroup : returns the components group
func (vn *VirtualNetwork) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (vn *VirtualNetwork) GetTags() map[string]string {
	return vn.Tags
}

// GetTag returns a components tag
func (vn *VirtualNetwork) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (vn *VirtualNetwork) Diff(c graph.Component) (diff.Changelog, error) {
	cvn, ok := c.(*VirtualNetwork)
	if ok {
		return diff.Diff(cvn, vn)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (vn *VirtualNetwork) Update(c graph.Component) {
	cvn, ok := c.(*VirtualNetwork)
	if ok {
		vn.ID = cvn.ID
	}

	vn.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (vn *VirtualNetwork) Rebuild(g *graph.Graph) {
	for x := 0; x < len(vn.Subnets); x++ {
		if vn.Subnets[x].SecurityGroupName != "" {
			vn.Subnets[x].SecurityGroup = `$(components.#[_component_id="` + TYPESECURITYGROUP + TYPEDELIMITER + vn.Subnets[x].SecurityGroupName + `"].id)`
		}
	}

	if len(vn.Subnets) < 1 {
		for _, c := range g.GetComponents().ByType("subnet") {
			s := c.(*Subnet)

			if s.ResourceGroupName == vn.ResourceGroupName && s.VirtualNetworkName == vn.Name {
				vn.Subnets = append(vn.Subnets, virtualnetwork.Subnet{
					Name:              s.Name,
					AddressPrefix:     s.AddressPrefix,
					SecurityGroup:     s.NetworkSecurityGroupID,
					SecurityGroupName: s.NetworkSecurityGroup,
				})
			}
		}
	}

	vn.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (vn *VirtualNetwork) Dependencies() (deps []string) {
	for x := 0; x < len(vn.Subnets); x++ {
		if vn.Subnets[x].SecurityGroupName != "" {
			deps = append(deps, TYPESECURITYGROUP+TYPEDELIMITER+vn.Subnets[x].SecurityGroupName)
		}
	}
	if len(deps) < 1 {
		deps = append(deps, TYPERESOURCEGROUP+TYPEDELIMITER+vn.ResourceGroupName)
	}
	return
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (vn *VirtualNetwork) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (vn *VirtualNetwork) Validate() error {
	log.Println("Validating Virtual network")
	val := validator.NewValidator()
	return val.Validate(vn)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (vn *VirtualNetwork) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (vn *VirtualNetwork) SetDefaultVariables() {
	vn.ProviderType = PROVIDERTYPE
	vn.ComponentType = TYPEVIRTUALNETWORK
	vn.ComponentID = TYPEVIRTUALNETWORK + TYPEDELIMITER + vn.Name
	vn.DatacenterName = DATACENTERNAME
	vn.DatacenterType = DATACENTERTYPE
	vn.DatacenterRegion = DATACENTERREGION
	vn.ClientID = CLIENTID
	vn.ClientSecret = CLIENTSECRET
	vn.TenantID = TENANTID
	vn.SubscriptionID = SUBSCRIPTIONID
	vn.Environment = ENVIRONMENT
}
