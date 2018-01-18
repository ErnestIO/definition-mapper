/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, q. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// Query : mapping of an query component
type Query struct {
	ID                string            `json:"id" diff:"-"`
	ProviderType      string            `json:"_provider" diff:"-"`
	ComponentType     string            `json:"_component" diff:"-"`
	ComponentID       string            `json:"_component_id" diff:"-"`
	State             string            `json:"_state" diff:"-"`
	Action            string            `json:"_action" diff:"-"`
	ResourceGroupName string            `json:"resource_group_name" diff:"-"`
	Tags              map[string]string `json:"tags" diff:"-"`
	DatacenterType    string            `json:"datacenter_type,omitempty" diff:"-"`
	DatacenterName    string            `json:"datacenter_name,omitempty" diff:"-"`
	DatacenterRegion  string            `json:"datacenter_region" diff:"-"`
	ClientID          string            `json:"azure_client_id" diff:"-"`
	ClientSecret      string            `json:"azure_client_secret" diff:"-"`
	TenantID          string            `json:"azure_tenant_id" diff:"-"`
	SubscriptionID    string            `json:"azure_subscription_id" diff:"-"`
	Environment       string            `json:"environment" diff:"-"`
}

// GetID : returns the component's ID
func (q *Query) GetID() string {
	return q.ComponentID
}

// GetName returns a components name
func (q *Query) GetName() string {
	return "query"
}

// GetProvider : returns the provider type
func (q *Query) GetProvider() string {
	return q.ProviderType
}

// GetProviderID returns a components provider id
func (q *Query) GetProviderID() string {
	return ""
}

// GetType : returns the type of the component
func (q *Query) GetType() string {
	return q.ComponentType
}

// GetState : returns the state of the component
func (q *Query) GetState() string {
	return q.State
}

// SetState : sets the state of the component
func (q *Query) SetState(s string) {
	q.State = s
}

// GetAction : returns the action of the component
func (q *Query) GetAction() string {
	return q.Action
}

// SetAction : Sets the action of the component
func (q *Query) SetAction(s string) {
	q.Action = s
}

// GetGroup : returns the components group
func (q *Query) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (q *Query) GetTags() map[string]string {
	return q.Tags
}

// GetTag returns a components tag
func (q *Query) GetTag(tag string) string {
	return ""
}

// Diff : diff's the component against another component of the same type
func (q *Query) Diff(c graph.Component) (diff.Changelog, error) {
	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (q *Query) Update(c graph.Component) {
	q.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (q *Query) Rebuild(g *graph.Graph) {
	q.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (q *Query) Dependencies() []string {
	return []string{}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (q *Query) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (q *Query) Validate() error {
	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (q *Query) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (q *Query) SetDefaultVariables() {
	q.ComponentID = q.ComponentType[0:len(q.ComponentType)-1] + TYPEDELIMITER + "query" + q.ResourceGroupName
	q.ProviderType = PROVIDERTYPE
	q.DatacenterType = DATACENTERTYPE
	q.DatacenterRegion = DATACENTERREGION
	q.ClientID = CLIENTID
	q.ClientSecret = CLIENTSECRET
	q.TenantID = TENANTID
	q.SubscriptionID = SUBSCRIPTIONID
	q.Environment = ENVIRONMENT
}
