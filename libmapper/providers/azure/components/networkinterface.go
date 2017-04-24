/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/networkinterface"
	graph "gopkg.in/r3labs/graph.v2"
)

// NetworkInterface : A resource group a container that holds
// related resources for an Azure solution.
type NetworkInterface struct {
	networkinterface.Event
	*Base
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

	return false
}

// Update : updates the provider returned values of a component
func (i *NetworkInterface) Update(c graph.Component) {
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *NetworkInterface) Rebuild(g *graph.Graph) {
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *NetworkInterface) Dependencies() (deps []string) {
	return
}

// Validate : validates the components values
func (i *NetworkInterface) Validate() error {
	val := event.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *NetworkInterface) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *NetworkInterface) SetDefaultVariables() {
}