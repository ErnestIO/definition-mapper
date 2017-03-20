/* This Source Code Form is subject to the terms cf the Mozilla Public
 * License, v. 2.0. If a copy cf the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"net"
	"reflect"

	graph "gopkg.in/r3labs/graph.v2"
)

// NatRule ...
type NatRule struct {
	Type            string `json:"type"`
	OriginIP        string `json:"origin_ip"`
	OriginPort      string `json:"origin_port"`
	TranslationIP   string `json:"translation_ip"`
	TranslationPort string `json:"translation_port"`
	Protocol        string `json:"protocol"`
	Network         string `json:"network"`
}

// FirewallRule ...
type FirewallRule struct {
	Name            string `json:"name"`
	SourceIP        string `json:"source_ip"`
	SourcePort      string `json:"source_port"`
	DestinationIP   string `json:"destination_ip"`
	DestinationPort string `json:"destination_port"`
	Protocol        string `json:"protocol"`
}

// Router : mapping of a router component
type Router struct {
	ProviderType       string            `json:"_provider"`
	ComponentType      string            `json:"_component"`
	ComponentID        string            `json:"_component_id"`
	State              string            `json:"_state"`
	Action             string            `json:"_action"`
	Name               string            `json:"name"`
	IP                 string            `json:"ip"`
	NatRules           []NatRule         `json:"nat_rules"`
	FirewallRules      []FirewallRule    `json:"firewall_rules"`
	Tags               map[string]string `json:"tags"`
	DatacenterType     string            `json:"datacenter_type"`
	DatacenterName     string            `json:"datacenter_name"`
	DatacenterUsername string            `json:"datacenter_username"`
	DatacenterPassword string            `json:"datacenter_password"`
	DatacenterRegion   string            `json:"datacenter_region"`
	VCloudURL          string            `json:"vcloud_url"`
	Service            string            `json:"service"`
}

// GetID : returns the component's ID
func (r *Router) GetID() string {
	return r.ComponentID
}

// GetName returns a components name
func (r *Router) GetName() string {
	return r.Name
}

// GetProvider : returns the provider type
func (r *Router) GetProvider() string {
	return r.ProviderType
}

// GetProviderID returns a components provider id
func (r *Router) GetProviderID() string {
	return r.Name
}

// GetType : returns the type cf the component
func (r *Router) GetType() string {
	return r.ComponentType
}

// GetState : returns the state cf the component
func (r *Router) GetState() string {
	return r.State
}

// SetState : sets the state cf the component
func (r *Router) SetState(s string) {
	r.State = s
}

// GetAction : returns the action cf the component
func (r *Router) GetAction() string {
	return r.Action
}

// SetAction : Sets the action cf the component
func (r *Router) SetAction(s string) {
	r.Action = s
}

// GetGroup : returns the components group
func (r *Router) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (r *Router) GetTags() map[string]string {
	return r.Tags
}

// GetTag returns a components tag
func (r *Router) GetTag(tag string) string {
	return r.Tags[tag]
}

// Diff : diff's the component against another component cf the same type
func (r *Router) Diff(c graph.Component) bool {
	cr, ok := c.(*Router)
	if ok {
		if len(r.FirewallRules) != len(cr.FirewallRules) {
			return true
		}

		if reflect.DeepEqual(r.FirewallRules, cr.FirewallRules) != true {
			return true
		}

		if len(r.NatRules) != len(cr.NatRules) {
			return true
		}

		if reflect.DeepEqual(r.NatRules, cr.NatRules) != true {
			return true
		}
	}

	return false
}

// Update : updates the provider returned values cf a component
func (r *Router) Update(c graph.Component) {
	r.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (r *Router) Rebuild(g *graph.Graph) {
	for i := 0; i < len(r.NatRules); i++ {
		r.NatRules[i].Network = EXTERNALNETWORK
	}

	r.SetDefaultVariables()
}

// Dependencies : returns a list cf component id's upon which the component depends
func (r *Router) Dependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (r *Router) Validate() error {
	for _, rule := range r.FirewallRules {
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

	for _, rule := range r.NatRules {
		// Check if Destination is a valid IP
		ip := net.ParseIP(rule.TranslationIP)
		if ip == nil {
			return errors.New("Port Forwarding must be a valid IP")
		}

		if rule.OriginIP != "" {
			source := net.ParseIP(rule.OriginIP)
			if source == nil {
				return errors.New("Port Forwarding source must be a valid IP")
			}
		}

		err := validatePort(rule.OriginPort, "Port Forwarding From")
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
func (r *Router) IsStateful() bool {
	// set to false because we can't delete a router
	return false
}

// SetDefaultVariables : sets up the default template variables for a component
func (r *Router) SetDefaultVariables() {
	r.ComponentType = TYPEROUTER
	r.ComponentID = TYPEROUTER + TYPEDELIMITER + r.Name
	r.ProviderType = PROVIDERTYPE
	r.DatacenterName = DATACENTERNAME
	r.DatacenterType = DATACENTERTYPE
	r.DatacenterRegion = DATACENTERREGION
	r.DatacenterUsername = DATACENTERUSERNAME
	r.DatacenterPassword = DATACENTERPASSWORD
	r.VCloudURL = VCLOUDURL
}
