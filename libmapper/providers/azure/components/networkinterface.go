/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"
	"strings"

	graph "gopkg.in/r3labs/graph.v2"
)

// IPConfiguration : ...
type IPConfiguration struct {
	Name                            string   `json:"name" validate:"required"`
	Subnet                          string   `json:"subnet_id" validate:"required"`
	PrivateIPAddress                string   `json:"private_ip_address"`
	PrivateIPAddressAllocation      string   `json:"private_ip_address_allocation" validate:"required"`
	PublicIPAddress                 string   `json:"public_ip_address_id"`
	LoadBalancerBackendAddressPools []string `json:"load_balancer_backend_address_pools_ids"`
	LoadBalancerInboundNatRules     []string `json:"load_balancer_inbound_nat_rules_ids"`
}

// NetworkInterface : A resource group a container that holds
// related resources for an Azure solution.
type NetworkInterface struct {
	ProviderType         string            `json:"_provider"`
	ComponentID          string            `json:"_component_id"`
	ComponentType        string            `json:"_component"`
	State                string            `json:"_state"`
	Action               string            `json:"_action"`
	DatacenterName       string            `json:"datacenter_name"`
	DatacenterType       string            `json:"datacenter_type"`
	DatacenterRegion     string            `json:"datacenter_region"`
	ID                   string            `json:"id"`
	Name                 string            `json:"name" validate:"required"`
	ResourceGroupName    string            `json:"resource_group_name" validate:"required"`
	Location             string            `json:"location" validate:"required"`
	NetworkSecurityGroup string            `json:"network_security_group_id"`
	MacAddress           string            `json:"mac_address"`
	PrivateIPAddress     string            `json:"private_ip_address"`
	VirtualMachineID     string            `json:"virtual_machine_id"`
	IPConfigurations     []IPConfiguration `json:"ip_configuration" validate:"min=1,dive"`
	DNSServers           []string          `json:"dns_servers" validate:"dive,ip"`
	InternalDNSNameLabel string            `json:"internal_dns_name_label"`
	AppliedDNSServers    []string          `json:"applied_dns_servers"`
	InternalFQDN         string            `json:"internal_fqdn"`
	EnableIPForwarding   bool              `json:"enable_ip_forwarding"`
	Tags                 map[string]string `json:"tags"`
	ClientID             string            `json:"azure_client_id"`
	ClientSecret         string            `json:"azure_client_secret"`
	TenantID             string            `json:"azure_tenant_id"`
	SubscriptionID       string            `json:"azure_subscription_id"`
	Environment          string            `json:"environment"`
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
func (i *NetworkInterface) Diff(c graph.Component) bool {
	cs, ok := c.(*NetworkInterface)
	if ok {
		if i.Name != cs.Name {
			return true
		}
		if i.NetworkSecurityGroup != cs.NetworkSecurityGroup {
			return true
		}
		if i.InternalDNSNameLabel != cs.InternalDNSNameLabel {
			return true
		}
		if i.ResourceGroupName != cs.ResourceGroupName {
			return true
		}
		if len(i.IPConfigurations) != len(cs.IPConfigurations) {
			return true
		}
		if len(i.DNSServers) != len(cs.DNSServers) {
			return true
		}
		for j := range i.DNSServers {
			if i.DNSServers[j] != cs.DNSServers[j] {
				return true
			}
		}
		for j := range i.IPConfigurations {
			if i.IPConfigurations[j].Name != cs.IPConfigurations[j].Name {
				return true
			}
			if i.IPConfigurations[j].PrivateIPAddress != cs.IPConfigurations[j].PrivateIPAddress {
				return true
			}
			if i.IPConfigurations[j].PrivateIPAddressAllocation != cs.IPConfigurations[j].PrivateIPAddressAllocation {
				return true
			}
			if i.IPConfigurations[j].PublicIPAddress != cs.IPConfigurations[j].PublicIPAddress {
				return true
			}
		}
	}
	return false
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
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *NetworkInterface) Dependencies() (deps []string) {
	if i.NetworkSecurityGroup != "" {
		deps = append(deps, TYPESECURITYGROUP+TYPEDELIMITER+i.NetworkSecurityGroup)
	}

	for _, config := range i.IPConfigurations {
		subnet := strings.Split(config.Subnet, "::")[1]
		subnet = strings.Split(subnet, `"]`)[0]
		deps = append(deps, TYPESUBNET+TYPEDELIMITER+subnet)
	}

	if len(deps) < 1 {
		return []string{TYPERESOURCEGROUP + TYPEDELIMITER + i.ResourceGroupName}
	}

	return
}

// Validate : validates the components values
func (i *NetworkInterface) Validate() error {
	log.Println("Validating azure network interfaces")
	val := NewValidator()
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
