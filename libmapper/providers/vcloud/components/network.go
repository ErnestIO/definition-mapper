/* This Source Code Form is subject to the terms cf the Mozilla Public
 * License, v. 2.0. If a copy cf the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"net"

	"github.com/r3labs/graph"
)

// Network : Mapping of a network component
type Network struct {
	ProviderType       string   `json:"_provider"`
	ComponentType      string   `json:"_component"`
	ComponentID        string   `json:"_component_id"`
	State              string   `json:"_state"`
	Action             string   `json:"_action"`
	Name               string   `json:"name"`
	Subnet             string   `json:"range"`
	Netmask            string   `json:"netmask"`
	StartAddress       string   `json:"start_address"`
	EndAddress         string   `json:"end_address"`
	Gateway            string   `json:"gateway"`
	DNS                []string `json:"dns"`
	Router             string   `json:"router_name"`
	DatacenterType     string   `json:"datacenter_type"`
	DatacenterName     string   `json:"datacenter_name"`
	DatacenterUsername string   `json:"datacenter_username"`
	DatacenterPassword string   `json:"datacenter_password"`
	DatacenterRegion   string   `json:"datacenter_region"`
	VCloudURL          string   `json:"vcloud_url"`
	Service            string   `json:"service"`
}

// GetID : returns the component's ID
func (n *Network) GetID() string {
	return n.ComponentID
}

// GetName returns a components name
func (n *Network) GetName() string {
	return n.Name
}

// GetProvider : returns the provider type
func (n *Network) GetProvider() string {
	return n.ProviderType
}

// GetProviderID returns a components provider id
func (n *Network) GetProviderID() string {
	return n.Name
}

// GetType : returns the type cf the component
func (n *Network) GetType() string {
	return n.ComponentType
}

// GetState : returns the state cf the component
func (n *Network) GetState() string {
	return n.State
}

// SetState : sets the state cf the component
func (n *Network) SetState(s string) {
	n.State = s
}

// GetAction : returns the action cf the component
func (n *Network) GetAction() string {
	return n.Action
}

// SetAction : Sets the action cf the component
func (n *Network) SetAction(s string) {
	n.Action = s
}

// GetGroup : returns the components group
func (n *Network) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (n *Network) GetTags() map[string]string {
	return map[string]string{}
}

// GetTag returns a components tag
func (n *Network) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component cf the same type
func (n *Network) Diff(c graph.Component) bool {
	return false
}

// Update : updates the provider returned values cf a component
func (n *Network) Update(c graph.Component) {
	n.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (n *Network) Rebuild(g *graph.Graph) {
	n.SetDefaultVariables()
}

// Dependencies : returns a list cf component id's upon which the component depends
func (n *Network) Dependencies() []string {
	return []string{TYPEROUTER + TYPEDELIMITER + n.Router}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (n *Network) SequentialDependencies() []string {
	return []string{TYPEROUTER + TYPEDELIMITER + n.Router}
}

// Validate : validates the components values
func (n *Network) Validate() error {
	_, _, err := net.ParseCIDR(n.Subnet)
	if err != nil {
		return errors.New("Network CIDR is not valid")
	}

	if n.Name == "" {
		return errors.New("Network name should not be null")
	}

	for _, val := range n.DNS {
		if ok := net.ParseIP(val); ok == nil {
			return errors.New("DNS " + val + " is not a valid CIDR")
		}
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (n *Network) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (n *Network) SetDefaultVariables() {
	n.ComponentType = TYPENETWORK
	n.ComponentID = TYPENETWORK + TYPEDELIMITER + n.Name
	n.ProviderType = PROVIDERTYPE
	n.DatacenterName = DATACENTERNAME
	n.DatacenterType = DATACENTERTYPE
	n.DatacenterRegion = DATACENTERREGION
	n.DatacenterUsername = DATACENTERUSERNAME
	n.DatacenterPassword = DATACENTERPASSWORD
	n.VCloudURL = VCLOUDURL
}
