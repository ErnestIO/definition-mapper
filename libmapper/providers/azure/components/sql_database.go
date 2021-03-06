/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/validator"
	"github.com/ernestio/ernestprovider/types/azure/sqldatabase"
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// SQLDatabase : A resource group a container that holds
// related resources for an Azure solution.
type SQLDatabase struct {
	sqldatabase.Event
	Base
}

// GetID : returns the component's ID
func (i *SQLDatabase) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *SQLDatabase) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *SQLDatabase) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *SQLDatabase) GetProviderID() string {
	return i.ID
}

// GetType : returns the type of the component
func (i *SQLDatabase) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *SQLDatabase) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *SQLDatabase) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *SQLDatabase) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *SQLDatabase) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *SQLDatabase) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (i *SQLDatabase) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *SQLDatabase) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (i *SQLDatabase) Diff(c graph.Component) (diff.Changelog, error) {
	cv, ok := c.(*SQLDatabase)
	if ok {
		return diff.Diff(cv, i)
	}

	return diff.Changelog{}, nil
}

func diffTags(t1, t2 map[string]string) bool {
	blackList := []string{"ernest_firewall_rules"}
	tags1 := map[string]string{}
	tags2 := map[string]string{}
	for k, t := range t1 {
		for _, bt := range blackList {
			if k != bt {
				tags1[k] = t
			}
		}
	}

	for k, t := range t2 {
		for _, bt := range blackList {
			if k != bt {
				tags2[k] = t
			}
		}
	}

	if len(tags1) != len(tags2) {
		return true
	}

	for kk1, tt1 := range tags1 {
		sw := false
		if val, ok := tags2[kk1]; ok {
			if val == tt1 {
				sw = true
			}
		}
		if sw == false {
			return true
		}
	}

	return false
}

// Update : updates the provider returned values of a component
func (i *SQLDatabase) Update(c graph.Component) {
	cs, ok := c.(*SQLDatabase)
	if ok {
		i.ID = cs.ID
		i.CreateMode = cs.CreateMode
		i.SourceDatabaseID = cs.SourceDatabaseID
		i.RestorePointInTime = cs.RestorePointInTime
		i.Edition = cs.Edition
		i.Collation = cs.Collation
		i.MaxSizeBytes = cs.MaxSizeBytes
		i.RequestedServiceObjectiveID = cs.RequestedServiceObjectiveID
		i.RequestedServiceObjectiveName = cs.RequestedServiceObjectiveName
		i.SourceDatabaseDeletionData = cs.SourceDatabaseDeletionData

		i.CreationDate = cs.CreationDate
		i.DefaultSecondaryLocation = cs.DefaultSecondaryLocation
	}
	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *SQLDatabase) Rebuild(g *graph.Graph) {
	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *SQLDatabase) Dependencies() (deps []string) {
	return []string{TYPESQLSERVER + TYPEDELIMITER + i.ServerName}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *SQLDatabase) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *SQLDatabase) Validate() error {
	log.Println("Validating SQL databases")
	val := validator.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *SQLDatabase) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *SQLDatabase) SetDefaultVariables() {
	i.ProviderType = PROVIDERTYPE
	i.ComponentType = TYPESQLDATABASE
	i.ComponentID = TYPESQLDATABASE + TYPEDELIMITER + i.Name
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.ClientID = CLIENTID
	i.ClientSecret = CLIENTSECRET
	i.TenantID = TENANTID
	i.SubscriptionID = SUBSCRIPTIONID
	i.Environment = ENVIRONMENT
}
