/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"
	"strings"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/lbrule"
	"github.com/r3labs/graph"
)

// LBRule : ..
type LBRule struct {
	ID string `json:"id"`
	lbrule.Event
	Base
}

// GetID : returns the component's ID
func (i *LBRule) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *LBRule) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *LBRule) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *LBRule) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *LBRule) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *LBRule) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *LBRule) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *LBRule) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *LBRule) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *LBRule) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *LBRule) GetTags() map[string]string {
	return map[string]string{}
}

// GetTag returns a components tag
func (i *LBRule) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *LBRule) Diff(c graph.Component) bool {
	cs, ok := c.(*LBRule)
	if ok {
		if i.FrontendIPConfigurationName != cs.FrontendIPConfigurationName {
			return true
		}
		if strings.EqualFold(i.Protocol, cs.Protocol) == false {
			return true
		}
		if i.FrontendPort != cs.FrontendPort {
			return true
		}
		if i.BackendPort != cs.BackendPort {
			return true
		}
		if i.BackendAddressPool != cs.BackendAddressPool {
			return true
		}
		if i.Probe != cs.Probe {
			return true
		}
		if i.EnableFloatingIP != cs.EnableFloatingIP {
			return true
		}
		if i.IdleTimeoutInMinutes != cs.IdleTimeoutInMinutes {
			return true
		}
		if i.LoadDistribution != cs.LoadDistribution {
			return true
		}
	}

	return false
}

// Update : updates the provider returned values of a component
func (i *LBRule) Update(c graph.Component) {
	cs, ok := c.(*LBRule)
	if ok {
		i.ID = cs.ID
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *LBRule) Rebuild(g *graph.Graph) {
	if i.Loadbalancer == "" && i.LoadbalancerID != "" {
		lb := g.GetComponents().ByProviderID(i.LoadbalancerID)
		if lb != nil {
			i.Loadbalancer = lb.GetName()
		}
	}

	if i.LoadbalancerID == "" && i.Loadbalancer != "" {
		i.LoadbalancerID = templLoadbalancerID(i.Loadbalancer)
	}

	if i.Probe == "" && i.ProbeID != "" {
		probe := g.GetComponents().ByProviderID(i.ProbeID)
		if probe != nil {
			i.Probe = probe.GetName()
		}
	}

	if i.ProbeID == "" && i.Probe != "" {
		i.ProbeID = templLoadbalancerProbeID(i.Probe)
	}

	if i.BackendAddressPool == "" && i.BackendAddressPoolID != "" {
		ap := g.GetComponents().ByProviderID(i.BackendAddressPoolID)
		if ap != nil {
			i.BackendAddressPool = ap.GetName()
		}
	}

	if i.BackendAddressPoolID == "" && i.BackendAddressPool != "" {
		i.BackendAddressPoolID = templLoadbalancerBackendAddressPoolID(i.BackendAddressPool)
	}

	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *LBRule) Dependencies() (deps []string) {
	if i.Probe != "" {
		deps = append(deps, TYPELBPROBE+TYPEDELIMITER+i.Probe)
	}

	if i.BackendAddressPool != "" {
		deps = append(deps, TYPELBBACKENDADDRESSPOOL+TYPEDELIMITER+i.BackendAddressPool)
	}

	if len(deps) < 1 {
		deps = append(deps, TYPELB+TYPEDELIMITER+i.Loadbalancer)
	}

	return
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *LBRule) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *LBRule) Validate() error {
	log.Println("Validating LB")
	val := event.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *LBRule) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *LBRule) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPELBRULE
	i.ComponentID = TYPELBRULE + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
