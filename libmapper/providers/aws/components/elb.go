/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"fmt"
	"sort"

	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// ELBListener ...
type ELBListener struct {
	FromPort int    `json:"from_port" diff:"from_port"`
	ToPort   int    `json:"to_port" diff:"to_port"`
	Protocol string `json:"protocol" diff:"protcol"`
	SSLCert  string `json:"ssl_cert" diff:"ssl_cert"`
}

// ELB : Mapping for a elb component
type ELB struct {
	ProviderType        string            `json:"_provider" diff:"-"`
	ComponentType       string            `json:"_component" diff:"-"`
	ComponentID         string            `json:"_component_id" diff:"_component_id,immutable"`
	State               string            `json:"_state" diff:"-"`
	Action              string            `json:"_action" diff:"-"`
	Name                string            `json:"name" diff:"-"`
	IsPrivate           bool              `json:"is_private" diff:"is_private,immutable"`
	DNSName             string            `json:"dns_name" diff:"dns_name,immutable"`
	Listeners           []ELBListener     `json:"listeners" diff:"listeners"`
	Networks            []string          `json:"networks" diff:"-"`
	NetworkAWSIDs       []string          `json:"network_aws_ids" diff:"-"`
	Instances           []string          `json:"instances" diff:"instances"`
	InstanceNames       sort.StringSlice  `json:"instance_names" diff:"instance_names"`
	InstanceAWSIDs      []string          `json:"instance_aws_ids" diff:"-"`
	SecurityGroups      sort.StringSlice  `json:"security_groups" diff:"security_groups"`
	SecurityGroupAWSIDs []string          `json:"security_group_aws_ids" diff:"-"`
	Tags                map[string]string `json:"tags" diff:"tags"`
	DatacenterType      string            `json:"datacenter_type,omitempty" diff:"-"`
	DatacenterName      string            `json:"datacenter_name,omitempty" diff:"-"`
	DatacenterRegion    string            `json:"datacenter_region" diff:"-"`
	AccessKeyID         string            `json:"aws_access_key_id" diff:"-"`
	SecretAccessKey     string            `json:"aws_secret_access_key" diff:"-"`
	Service             string            `json:"service" diff:"-"`
}

// GetID : returns the component's ID
func (e *ELB) GetID() string {
	return e.ComponentID
}

// GetName returns a components name
func (e *ELB) GetName() string {
	return e.Name
}

// GetProvider : returns the provider type
func (e *ELB) GetProvider() string {
	return e.ProviderType
}

// GetProviderID returns a components provider id
func (e *ELB) GetProviderID() string {
	return e.Name
}

// GetType : returns the type of the component
func (e *ELB) GetType() string {
	return e.ComponentType
}

// GetState : returns the state of the component
func (e *ELB) GetState() string {
	return e.State
}

// SetState : sets the state of the component
func (e *ELB) SetState(s string) {
	e.State = s
}

// GetAction : returns the action of the component
func (e *ELB) GetAction() string {
	return e.Action
}

// SetAction : Sets the action of the component
func (e *ELB) SetAction(s string) {
	e.Action = s
}

// GetGroup : returns the components group
func (e *ELB) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (e *ELB) GetTags() map[string]string {
	return e.Tags
}

// GetTag returns a components tag
func (e *ELB) GetTag(tag string) string {
	return e.Tags[tag]
}

// Diff : diff's the component against another component of the same type
func (e *ELB) Diff(c graph.Component) (diff.Changelog, error) {
	ce, ok := c.(*ELB)
	if ok {
		return diff.Diff(ce, e)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (e *ELB) Update(c graph.Component) {
	ce, ok := c.(*ELB)
	if ok {
		e.DNSName = ce.DNSName
	}

	e.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (e *ELB) Rebuild(g *graph.Graph) {
	if len(e.Networks) > len(e.NetworkAWSIDs) {
		for _, nw := range e.Networks {
			e.NetworkAWSIDs = append(e.NetworkAWSIDs, templSubnetID(nw))
		}
	}

	if len(e.NetworkAWSIDs) > len(e.Networks) {
		for _, nwid := range e.NetworkAWSIDs {
			nw := g.GetComponents().ByProviderID(nwid)
			if nw != nil {
				e.Networks = append(e.Networks, nw.GetName())
			}
		}
	}

	if len(e.SecurityGroups) > len(e.SecurityGroupAWSIDs) {
		for _, sg := range e.SecurityGroups {
			e.SecurityGroupAWSIDs = append(e.SecurityGroupAWSIDs, templSecurityGroupID(sg))
		}
	}

	if len(e.SecurityGroupAWSIDs) > len(e.SecurityGroups) {
		for _, sgid := range e.SecurityGroupAWSIDs {
			sg := g.GetComponents().ByProviderID(sgid)
			if sg != nil {
				e.SecurityGroups = append(e.SecurityGroups, sg.GetName())
			}
		}
	}

	if len(e.Instances) > len(e.InstanceAWSIDs) {
		for _, ig := range e.Instances {
			for _, i := range g.GetComponents().ByGroup(GROUPINSTANCE, ig) {
				e.InstanceAWSIDs = append(e.InstanceAWSIDs, templInstanceID(i.GetName()))
			}
		}
	}

	if len(e.InstanceAWSIDs) > len(e.Instances) {
		for _, iid := range e.InstanceAWSIDs {
			i := g.GetComponents().ByProviderID(iid)
			if i != nil {
				e.Instances = appendUnique(e.Instances, i.GetTag(GROUPINSTANCE))
			}
		}
	}

	for _, ig := range e.Instances {
		for _, i := range g.GetComponents().ByGroup(GROUPINSTANCE, ig) {
			e.InstanceNames = appendUnique(e.InstanceNames, i.GetName())
		}
	}

	e.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (e *ELB) Dependencies() []string {
	var deps []string

	for _, sg := range e.SecurityGroups {
		deps = append(deps, TYPESECURITYGROUP+TYPEDELIMITER+sg)
	}

	for _, nw := range e.Networks {
		deps = append(deps, TYPENETWORK+TYPEDELIMITER+nw)
	}

	for _, in := range e.InstanceNames {
		deps = append(deps, TYPEINSTANCE+TYPEDELIMITER+in)
	}

	return deps
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (e *ELB) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (e *ELB) Validate() error {
	if e.Name == "" {
		return errors.New("ELB name should not be null")
	}

	if len(e.Listeners) < 1 {
		return errors.New("ELB must contain more than one listeners")
	}

	if e.IsPrivate != true && len(e.Networks) < 1 {
		return errors.New("ELB must specify at least one subnet if public")
	}

	/*
		if nw == n.Name && n.Public != true && e.Private != true {
			return fmt.Errorf("ELB subnet (%s) is not a public subnet", nw)
		}
	*/

	for _, listener := range e.Listeners {
		if listener.FromPort < 1 || listener.FromPort > 65535 {
			return fmt.Errorf("From Port (%d) is out of range [1 - 65535]", listener.FromPort)
		}

		if listener.ToPort < 1 || listener.ToPort > 65535 {
			return fmt.Errorf("From Port (%d) is out of range [1 - 65535]", listener.ToPort)
		}

		if listener.Protocol != "HTTP" &&
			listener.Protocol != "HTTPS" &&
			listener.Protocol != "TCP" &&
			listener.Protocol != "SSL" {
			return errors.New("ELB Protocol must be one of http, https, tcp or ssl")
		}

		if listener.Protocol == "https" && listener.SSLCert == "" || listener.Protocol == "ssl" && listener.SSLCert == "" {
			return errors.New("ELB listener must specify an ssl cert when protocol is https/ssl")
		}

	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (e *ELB) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (e *ELB) SetDefaultVariables() {
	e.ComponentType = TYPEELB
	e.ComponentID = TYPEELB + TYPEDELIMITER + e.Name
	e.ProviderType = PROVIDERTYPE
	e.DatacenterName = DATACENTERNAME
	e.DatacenterType = DATACENTERTYPE
	e.DatacenterRegion = DATACENTERREGION
	e.AccessKeyID = ACCESSKEYID
	e.SecretAccessKey = SECRETACCESSKEY
}
