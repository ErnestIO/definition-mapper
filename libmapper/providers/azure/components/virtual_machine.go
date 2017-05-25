/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"reflect"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/virtualmachine"
	graph "gopkg.in/r3labs/graph.v2"
)

// VirtualMachine : A resource group a container that holds
// related resources for an Azure solution.
type VirtualMachine struct {
	virtualmachine.Event
	Base
}

// GetID : returns the component's ID
func (i *VirtualMachine) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *VirtualMachine) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *VirtualMachine) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *VirtualMachine) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *VirtualMachine) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *VirtualMachine) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *VirtualMachine) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *VirtualMachine) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *VirtualMachine) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *VirtualMachine) GetGroup() string {
	return i.Tags["ernest.instance_group"]
}

// GetTags returns a components tags
func (i *VirtualMachine) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *VirtualMachine) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *VirtualMachine) Diff(c graph.Component) bool {
	cvm, ok := c.(*VirtualMachine)
	if ok {
		if i.VMSize != cvm.VMSize {
			return true
		}

		if i.StorageDataDisk.Size != cvm.StorageDataDisk.Size {
			return true
		}

		if reflect.DeepEqual(i.NetworkInterfaces, cvm.NetworkInterfaces) != true {
			return true
		}

		if reflect.DeepEqual(i.Tags, cvm.Tags) != true {
			return true
		}
	}

	return false
}

// Update : updates the provider returned values of a component
func (i *VirtualMachine) Update(c graph.Component) {
	cvm, ok := c.(*VirtualMachine)
	if ok {
		i.ID = cvm.ID
		// ???
		i.StorageDataDisk.Lun = cvm.StorageDataDisk.Lun
	}

	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *VirtualMachine) Rebuild(g *graph.Graph) {
	if len(i.NetworkInterfaces) > len(i.NetworkInterfaceIDs) {
		for _, iface := range i.NetworkInterfaces {
			i.NetworkInterfaceIDs = append(i.NetworkInterfaceIDs, templNetworkInterfaceID(iface))
		}
	}

	if len(i.NetworkInterfaceIDs) > len(i.NetworkInterfaces) {
		for _, id := range i.NetworkInterfaceIDs {
			iface := g.GetComponents().ByProviderID(id)
			if iface != nil {
				i.NetworkInterfaces = append(i.NetworkInterfaces, iface.GetName())
			}
		}
	}

	if i.AvailabilitySet == "" && i.AvailabilitySetID != "" {
		as := g.GetComponents().ByProviderID(i.AvailabilitySetID)
		if as != nil {
			i.AvailabilitySet = as.GetName()
		}
	}

	if i.AvailabilitySetID == "" && i.AvailabilitySet != "" {
		i.AvailabilitySetID = templAvailabilitySetID(i.AvailabilitySet)
	}

	if i.StorageOSDisk.ManagedDisk == "" && i.StorageOSDisk.ManagedDiskID != "" {
		md := g.GetComponents().ByProviderID(i.StorageOSDisk.ManagedDiskID)
		if md != nil {
			i.StorageOSDisk.ManagedDisk = md.GetName()
		}
	}

	if i.StorageOSDisk.ManagedDiskID == "" && i.StorageOSDisk.ManagedDisk != "" {
		i.StorageOSDisk.ManagedDiskID = templManagedDiskID(i.StorageOSDisk.ManagedDisk)
	}

	if i.StorageDataDisk.ManagedDisk == "" && i.StorageDataDisk.ManagedDiskID != "" {
		md := g.GetComponents().ByProviderID(i.StorageDataDisk.ManagedDiskID)
		if md != nil {
			i.StorageDataDisk.ManagedDisk = md.GetName()
		}
	}

	if i.StorageDataDisk.ManagedDiskID == "" && i.StorageDataDisk.ManagedDisk != "" {
		i.StorageDataDisk.ManagedDiskID = templManagedDiskID(i.StorageDataDisk.ManagedDisk)
	}

	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *VirtualMachine) Dependencies() (deps []string) {
	for _, iface := range i.NetworkInterfaces {
		deps = append(deps, TYPENETWORKINTERFACE+TYPEDELIMITER+iface)
	}

	if i.StorageOSDisk.StorageContainer != "" {
		deps = append(deps, TYPESTORAGECONTAINER+TYPEDELIMITER+i.StorageOSDisk.StorageContainer)
	}

	if i.StorageOSDisk.ManagedDisk != "" {
		deps = append(deps, TYPEMANAGEDDISK+TYPEDELIMITER+i.StorageOSDisk.ManagedDisk)
	}

	if i.StorageDataDisk.StorageContainer != "" && i.StorageDataDisk.StorageContainer != i.StorageOSDisk.StorageContainer {
		deps = append(deps, TYPESTORAGECONTAINER+TYPEDELIMITER+i.StorageDataDisk.StorageContainer)
	}

	if i.StorageDataDisk.ManagedDisk != "" {
		deps = append(deps, TYPEMANAGEDDISK+TYPEDELIMITER+i.StorageDataDisk.ManagedDisk)
	}

	if i.AvailabilitySet != "" {
		deps = append(deps, TYPEAVAILABILITYSET+TYPEDELIMITER+i.AvailabilitySet)
	}

	if len(deps) < 1 {
		return []string{TYPERESOURCEGROUP + TYPEDELIMITER + i.ResourceGroupName}
	}

	return
}

// Validate : validates the components values
func (i *VirtualMachine) Validate() error {
	val := event.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *VirtualMachine) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *VirtualMachine) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPEVIRTUALMACHINE
	i.ComponentID = TYPEVIRTUALMACHINE + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
