/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/validator"
	"github.com/ernestio/ernestprovider/types/azure/networkinterface"
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// NetworkInterface : A resource group a container that holds
// related resources for an Azure solution.
type NetworkInterface struct {
	networkinterface.Event
	Base
}

// GetID : returns the component's ID
func (i *NetworkInterface) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *NetworkInterface) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *NetworkInterface) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *NetworkInterface) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *NetworkInterface) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *NetworkInterface) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *NetworkInterface) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *NetworkInterface) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *NetworkInterface) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *NetworkInterface) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *NetworkInterface) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *NetworkInterface) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *NetworkInterface) Diff(c graph.Component) (diff.Changelog, error) {
	cs, ok := c.(*NetworkInterface)
	if ok {
		return diff.Diff(cs, i)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (i *NetworkInterface) Update(c graph.Component) {
	cs, ok := c.(*NetworkInterface)
	if ok {
		i.ID = cs.ID
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *NetworkInterface) Rebuild(g *graph.Graph) {
	if i.NetworkSecurityGroup == "" && i.NetworkSecurityGroupID != "" {
		sg := g.GetComponents().ByProviderID(i.NetworkSecurityGroupID)
		if sg != nil {
			i.NetworkSecurityGroup = sg.GetName()
		}
	}

	if i.NetworkSecurityGroupID == "" && i.NetworkSecurityGroup != "" {
		i.NetworkSecurityGroupID = templSecurityGroupID(i.NetworkSecurityGroup)
	}

	for x := 0; x < len(i.IPConfigurations); x++ {
		if i.IPConfigurations[x].SubnetID == "" && i.IPConfigurations[x].Subnet != "" {
			i.IPConfigurations[x].SubnetID = templSubnetID(i.IPConfigurations[x].Subnet)
		}

		if i.IPConfigurations[x].Subnet == "" && i.IPConfigurations[x].SubnetID != "" {
			s := g.GetComponents().ByProviderID(i.IPConfigurations[x].SubnetID)
			if s != nil {
				i.IPConfigurations[x].Subnet = s.GetName()
			}
		}

		if i.IPConfigurations[x].PublicIPAddress == "" && i.IPConfigurations[x].PublicIPAddressID != "" {
			ip := g.GetComponents().ByProviderID(i.IPConfigurations[x].PublicIPAddressID)
			if ip != nil {
				i.IPConfigurations[x].PublicIPAddress = ip.GetName()
			}
		}

		if i.IPConfigurations[x].PublicIPAddressID == "" && i.IPConfigurations[x].PublicIPAddress != "" {
			i.IPConfigurations[x].PublicIPAddressID = templPublicIPAddressID(i.IPConfigurations[x].PublicIPAddress)
		}

		if len(i.IPConfigurations[x].LoadBalancerBackendAddressPoolIDs) > len(i.IPConfigurations[x].LoadBalancerBackendAddressPools) {
			i.IPConfigurations[x].LoadBalancerBackendAddressPools = []string{}
			for _, ap := range i.IPConfigurations[x].LoadBalancerBackendAddressPoolIDs {
				ap := g.GetComponents().ByProviderID(ap)
				if ap != nil {
					i.IPConfigurations[x].LoadBalancerBackendAddressPools = append(i.IPConfigurations[x].LoadBalancerBackendAddressPools, ap.GetName())
				}
			}
		}

		if len(i.IPConfigurations[x].LoadBalancerBackendAddressPools) > len(i.IPConfigurations[x].LoadBalancerBackendAddressPoolIDs) {
			i.IPConfigurations[x].LoadBalancerBackendAddressPoolIDs = []string{}
			for _, ap := range i.IPConfigurations[x].LoadBalancerBackendAddressPools {
				i.IPConfigurations[x].LoadBalancerBackendAddressPoolIDs = append(i.IPConfigurations[x].LoadBalancerBackendAddressPoolIDs, templLoadbalancerBackendAddressPoolID(ap))
			}
		}
	}

	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *NetworkInterface) Dependencies() (deps []string) {
	if i.NetworkSecurityGroup != "" {
		deps = append(deps, TYPESECURITYGROUP+TYPEDELIMITER+i.NetworkSecurityGroup)
	}

	for _, config := range i.IPConfigurations {
		if config.Subnet != "" {
			deps = append(deps, TYPESUBNET+TYPEDELIMITER+config.Subnet)
		}
		if config.PublicIPAddress != "" {
			deps = append(deps, TYPEPUBLICIP+TYPEDELIMITER+config.PublicIPAddress)
		}
		for _, ap := range config.LoadBalancerBackendAddressPools {
			deps = append(deps, TYPELBBACKENDADDRESSPOOL+TYPEDELIMITER+ap)
		}
	}

	if len(deps) < 1 {
		return []string{TYPERESOURCEGROUP + TYPEDELIMITER + i.ResourceGroupName}
	}

	return
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *NetworkInterface) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *NetworkInterface) Validate() error {
	log.Println("Validating azure network interfaces")
	val := validator.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *NetworkInterface) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *NetworkInterface) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPENETWORKINTERFACE
	i.ComponentID = TYPENETWORKINTERFACE + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
