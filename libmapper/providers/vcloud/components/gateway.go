/* This Source Code Form is subject to the terms cf the Mozilla Public
 * License, v. 2.0. If a copy cf the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"

	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// NatRule ...
type NatRule struct {
	Type            string `json:"type" diff:"type"`
	OriginIP        string `json:"origin_ip" diff:"origin_ip"`
	OriginPort      string `json:"origin_port" diff:"origin_port"`
	TranslationIP   string `json:"translation_ip" diff:"translation_ip"`
	TranslationPort string `json:"translation_port" diff:"translation_port"`
	Protocol        string `json:"protocol" diff:"protocol"`
}

// FirewallRule ...
type FirewallRule struct {
	Name            string `json:"name" diff:"name,identifier"`
	SourceIP        string `json:"source_ip" diff:"source_ip"`
	SourcePort      string `json:"source_port" diff:"source_port"`
	DestinationIP   string `json:"destination_ip" diff:"destination_ip"`
	DestinationPort string `json:"destination_port" diff:"destination_port"`
	Protocol        string `json:"protocol" diff:"protocol"`
}

// Gateway : mapping of a edge gateway component
type Gateway struct {
	Base
	ID            string         `json:"id" diff:"-"`
	Name          string         `json:"name" diff:"-"`
	NatRules      []NatRule      `json:"nat_rules" diff:"nat_rules"`
	FirewallRules []FirewallRule `json:"firewall_rules" diff:"firewall_rules"`
}

// GetID : returns the component's ID
func (gw *Gateway) GetID() string {
	return gw.ComponentID
}

// GetName returns a components name
func (gw *Gateway) GetName() string {
	return gw.Name
}

// GetProvider : returns the provider type
func (gw *Gateway) GetProvider() string {
	return gw.ProviderType
}

// GetProviderID returns a components provider id
func (gw *Gateway) GetProviderID() string {
	return gw.Name
}

// GetType : returns the type cf the component
func (gw *Gateway) GetType() string {
	return gw.ComponentType
}

// GetState : returns the state cf the component
func (gw *Gateway) GetState() string {
	return gw.State
}

// SetState : sets the state cf the component
func (gw *Gateway) SetState(s string) {
	gw.State = s
}

// GetAction : returns the action cf the component
func (gw *Gateway) GetAction() string {
	return gw.Action
}

// SetAction : Sets the action cf the component
func (gw *Gateway) SetAction(s string) {
	gw.Action = s
}

// GetGroup : returns the components group
func (gw *Gateway) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (gw *Gateway) GetTags() map[string]string {
	return map[string]string{}
}

// GetTag returns a components tag
func (gw *Gateway) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component cf the same type
func (gw *Gateway) Diff(c graph.Component) (diff.Changelog, error) {
	cgw, ok := c.(*Gateway)
	if ok {
		return diff.Diff(cgw, gw)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values cf a component
func (gw *Gateway) Update(c graph.Component) {
	cgw := c.(*Gateway)

	gw.ID = cgw.ID

	gw.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (gw *Gateway) Rebuild(g *graph.Graph) {
	gw.SetDefaultVariables()
}

// Dependencies : returns a list cf component id's upon which the component depends
func (gw *Gateway) Dependencies() []string {
	return []string{}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (gw *Gateway) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (gw *Gateway) Validate() error {
	for _, rule := range gw.FirewallRules {
		// Check if firewall rule name is null
		if rule.Name == "" {
			return errors.New("Firewall Rule name should not be null")
		}

		err := validateIP(rule.SourceIP, "Firewall Rule Source")
		if err != nil {
			return err
		}

		err = validateIP(rule.DestinationIP, "Firewall Rule Destination")
		if err != nil {
			return err
		}

		// Validate FromPort Port
		// Must be: [any | 1 - 65535]
		err = validatePort(rule.SourcePort, "Firewall Rule From")
		if err != nil {
			return err
		}

		// Validate ToPort Port
		// Must be: [any | 1 - 65535]
		err = validatePort(rule.DestinationPort, "Firewall Rule To")
		if err != nil {
			return err
		}

		// Validate Protocol
		// Must be one of: tcp | udp | icmp | any | tcp & udp
		err = validateProtocol(rule.Protocol)
		if err != nil {
			return err
		}
	}

	for _, rule := range gw.NatRules {
		// Check if Destination is a valid IP
		err := validateIP(rule.OriginIP, "Nat Rule Source")
		if err != nil {
			return err
		}

		err = validateIP(rule.TranslationIP, "Nat Rule Destination")
		if err != nil {
			return err
		}

		err = validatePort(rule.OriginPort, "Port Forwarding From")
		if err != nil {
			return err
		}

		err = validatePort(rule.TranslationPort, "Port Forwarding To")
		if err != nil {
			return err
		}
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (gw *Gateway) IsStateful() bool {
	// set to false because we can't delete a router
	return false
}

// SetDefaultVariables : sets up the default template variables for a component
func (gw *Gateway) SetDefaultVariables() {
	gw.ComponentType = TYPEROUTER
	gw.ComponentID = TYPEROUTER + TYPEDELIMITER + gw.Name
	gw.ProviderType = PROVIDERTYPE
	gw.Credentials = &Credentials{
		Type:      DATACENTERTYPE,
		Vdc:       DATACENTERNAME,
		Username:  DATACENTERUSERNAME,
		Password:  DATACENTERPASSWORD,
		VCloudURL: VCLOUDURL,
	}
}
