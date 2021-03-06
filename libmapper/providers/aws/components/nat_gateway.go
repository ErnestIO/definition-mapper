/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"

	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// NatGateway : mapping of a nat component
type NatGateway struct {
	ProviderType           string            `json:"_provider" diff:"-"`
	ComponentType          string            `json:"_component" diff:"-"`
	ComponentID            string            `json:"_component_id" diff:"_component_id,immutable"`
	State                  string            `json:"_state" diff:"-"`
	Action                 string            `json:"_action" diff:"-"`
	NatGatewayAWSID        string            `json:"nat_gateway_aws_id" diff:"-"`
	Name                   string            `json:"name" diff:"-"`
	PublicNetwork          string            `json:"public_network" diff:"-"`
	RoutedNetworks         []string          `json:"routed_networks" diff:"routed_networks"`
	RoutedNetworkAWSIDs    []string          `json:"routed_networks_aws_ids" diff:"-"`
	PublicNetworkAWSID     string            `json:"public_network_aws_id" diff:"-"`
	NatGatewayAllocationID string            `json:"nat_gateway_allocation_id" diff:"-"`
	NatGatewayAllocationIP string            `json:"nat_gateway_allocation_ip" diff:"-"`
	InternetGatewayID      string            `json:"internet_gateway_id" diff:"-"`
	DatacenterType         string            `json:"datacenter_type" diff:"-"`
	DatacenterName         string            `json:"datacenter_name" diff:"-"`
	DatacenterRegion       string            `json:"datacenter_region" diff:"-"`
	AccessKeyID            string            `json:"aws_access_key_id" diff:"-"`
	SecretAccessKey        string            `json:"aws_secret_access_key" diff:"-"`
	VpcID                  string            `json:"vpc_id" diff:"-"`
	Remove                 bool              `json:"-" diff:"-"`
	Tags                   map[string]string `json:"tags" diff:"-"`
	Service                string            `json:"service" diff:"-"`
}

// GetID : returns the component's ID
func (n *NatGateway) GetID() string {
	return n.ComponentID
}

// GetName returns a components name
func (n *NatGateway) GetName() string {
	return n.Name
}

// GetProvider : returns the provider type
func (n *NatGateway) GetProvider() string {
	return n.ProviderType
}

// GetProviderID returns a components provider id
func (n *NatGateway) GetProviderID() string {
	return n.NatGatewayAWSID
}

// GetType : returns the type of the component
func (n *NatGateway) GetType() string {
	return n.ComponentType
}

// GetState : returns the state of the component
func (n *NatGateway) GetState() string {
	return n.State
}

// SetState : sets the state of the component
func (n *NatGateway) SetState(s string) {
	n.State = s
}

// GetAction : returns the action of the component
func (n *NatGateway) GetAction() string {
	return n.Action
}

// SetAction : Sets the action of the component
func (n *NatGateway) SetAction(s string) {
	n.Action = s
}

// GetGroup : returns the components group
func (n *NatGateway) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (n *NatGateway) GetTags() map[string]string {
	return n.Tags
}

// GetTag returns a components tag
func (n *NatGateway) GetTag(tag string) string {
	return n.Tags[tag]
}

// Diff : diff's the component against another component of the same type
func (n *NatGateway) Diff(c graph.Component) (diff.Changelog, error) {
	cn, ok := c.(*NatGateway)
	if ok {
		return diff.Diff(cn, n)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (n *NatGateway) Update(c graph.Component) {
	cn, ok := c.(*NatGateway)
	if ok {
		n.NatGatewayAWSID = cn.NatGatewayAWSID
		n.NatGatewayAllocationID = cn.NatGatewayAllocationID
		n.NatGatewayAllocationIP = cn.NatGatewayAllocationIP
	}

	n.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (n *NatGateway) Rebuild(g *graph.Graph) {
	for _, nw := range n.RoutedNetworkAWSIDs {
		fn := g.GetComponents().ByProviderID(nw)
		if fn != nil {
			n.Name = fn.GetTag("ernest.nat_gateway")
		}
	}

	if n.PublicNetwork == "" && n.PublicNetworkAWSID != "" {
		pn := g.GetComponents().ByProviderID(n.PublicNetworkAWSID)
		if pn != nil {
			n.PublicNetwork = pn.GetName()
		} else {
			n.Remove = true
			return
		}
	}

	if n.PublicNetworkAWSID == "" && n.PublicNetwork != "" {
		n.PublicNetworkAWSID = templSubnetID(n.PublicNetwork)
	}

	if len(n.RoutedNetworks) > len(n.RoutedNetworkAWSIDs) {
		for _, nw := range n.RoutedNetworks {
			n.RoutedNetworkAWSIDs = append(n.RoutedNetworkAWSIDs, templSubnetID(nw))
		}
	}

	if len(n.RoutedNetworkAWSIDs) > len(n.RoutedNetworks) {
		for _, nwid := range n.RoutedNetworkAWSIDs {
			nw := g.GetComponents().ByProviderID(nwid)
			if nw != nil {
				n.RoutedNetworks = append(n.RoutedNetworks, nw.GetName())
			}
		}
	}

	n.VpcID = templSubnetVPCID(n.PublicNetwork)

	n.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (n *NatGateway) Dependencies() []string {
	var deps []string

	for _, nw := range n.RoutedNetworks {
		deps = append(deps, TYPENETWORK+TYPEDELIMITER+nw)
	}

	deps = append(deps, TYPENETWORK+TYPEDELIMITER+n.PublicNetwork)

	return deps
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (n *NatGateway) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (n *NatGateway) Validate() error {
	if n.Name == "" {
		return errors.New("Nat Gateway name should not be null")
	}

	if n.PublicNetwork == "" {
		return errors.New("Nat Gateway should specify a public network")
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (n *NatGateway) IsStateful() bool {
	return !n.Remove
}

// SetDefaultVariables : sets up the default template variables for a component
func (n *NatGateway) SetDefaultVariables() {
	n.ComponentType = TYPENATGATEWAY
	n.ComponentID = TYPENATGATEWAY + TYPEDELIMITER + n.Name
	n.ProviderType = PROVIDERTYPE
	n.DatacenterName = DATACENTERNAME
	n.DatacenterType = DATACENTERTYPE
	n.DatacenterRegion = DATACENTERREGION
	n.AccessKeyID = ACCESSKEYID
	n.SecretAccessKey = SECRETACCESSKEY
}
