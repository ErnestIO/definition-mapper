/* This Source Code Form is subject to the terms cf the Mozilla Public
 * License, v. 2.0. If a copy cf the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"net"

	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// Network : Mapping of a network component
type Network struct {
	Base
	ID            string   `json:"id" diff:"-"`
	Name          string   `json:"name" diff:"-"`
	Subnet        string   `json:"range" diff:"-"`
	Netmask       string   `json:"netmask" diff:"-"`
	StartAddress  string   `json:"start_address" diff:"-"`
	EndAddress    string   `json:"end_address" diff:"-"`
	Gateway       string   `json:"gateway" diff:"-"`
	DNS           []string `json:"dns" diff:"dns"`
	EdgeGateway   string   `json:"edge_gateway" diff:"-"`
	EdgeGatewayID string   `json:"edge_gateway_id" diff:"-"`
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
func (n *Network) Diff(c graph.Component) (diff.Changelog, error) {
	cn, ok := c.(*Network)
	if ok {
		diff.Diff(cn, n)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values cf a component
func (n *Network) Update(c graph.Component) {
	cn := c.(*Network)

	n.ID = cn.ID
	n.EdgeGatewayID = cn.EdgeGatewayID

	n.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (n *Network) Rebuild(g *graph.Graph) {
	n.SetDefaultVariables()
}

// Dependencies : returns a list cf component id's upon which the component depends
func (n *Network) Dependencies() []string {
	return []string{TYPEROUTER + TYPEDELIMITER + n.EdgeGateway}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (n *Network) SequentialDependencies() []string {
	return []string{TYPEROUTER + TYPEDELIMITER + n.EdgeGateway}
}

// Validate : validates the components values
func (n *Network) Validate() error {
	if n.Name == "" {
		return errors.New("Network name should not be null")
	}

	for _, val := range n.DNS {
		if val == "" {
			continue
		}
		if ok := net.ParseIP(val); ok == nil {
			return errors.New("DNS " + val + " is not a valid CIDR")
		}
	}

	if n.Subnet == "" && n.Gateway != "" && n.Netmask != "" {
		return nil
	}

	_, _, err := net.ParseCIDR(n.Subnet)
	if err != nil {
		return errors.New("Network CIDR is not valid")
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
	n.Credentials = &Credentials{
		Type:      DATACENTERTYPE,
		Vdc:       DATACENTERNAME,
		Username:  DATACENTERUSERNAME,
		Password:  DATACENTERPASSWORD,
		VCloudURL: VCLOUDURL,
	}
}
