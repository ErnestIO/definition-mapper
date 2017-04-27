/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"log"

	"github.com/ernestio/ernestprovider/event"
	"github.com/ernestio/ernestprovider/providers/azure/sqldatabase"
	graph "gopkg.in/r3labs/graph.v2"
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
func (i *SQLDatabase) Diff(c graph.Component) bool {
	cv, ok := c.(*SQLDatabase)
	if ok {
		if i.ServerName != cv.Name {
			return true
		}
		if i.CreateMode != cv.CreateMode {
			return true
		}
		if i.SourceDatabaseID != cv.SourceDatabaseID {
			return true
		}
		if i.RestorePointInTime != cv.RestorePointInTime {
			return true
		}
		if i.Edition != cv.Edition {
			return true
		}
		if i.Collation != cv.Collation {
			return true
		}
		if i.MaxSizeBytes != cv.MaxSizeBytes {
			return true
		}
		if i.RequestedServiceObjectiveID != cv.RequestedServiceObjectiveID {
			return true
		}
		if i.RequestedServiceObjectiveName != cv.RequestedServiceObjectiveName {
			return true
		}
		if i.SourceDatabaseDeletionData != cv.SourceDatabaseDeletionData {
			return true
		}
		if i.ElasticPoolName != cv.ElasticPoolName {
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

// Validate : validates the components values
func (i *SQLDatabase) Validate() error {
	log.Println("Validating SQL databases")
	val := event.NewValidator()
	return val.Validate(i)
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *SQLDatabase) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *SQLDatabase) SetDefaultVariables() {
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
