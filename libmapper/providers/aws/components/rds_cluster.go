/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// RDSCluster ...
type RDSCluster struct {
	ProviderType        string            `json:"_provider" diff:"-"`
	ComponentType       string            `json:"_component" diff:"-"`
	ComponentID         string            `json:"_component_id" diff:"-"`
	State               string            `json:"_state" diff:"-"`
	Action              string            `json:"_action" diff:"-"`
	ARN                 string            `json:"arn" diff:"-"`
	Name                string            `json:"name" diff:"name,immutable"`
	Engine              string            `json:"engine" diff:"-"`
	EngineVersion       string            `json:"engine_version,omitempty" diff:"-"`
	Port                *int64            `json:"port,omitempty" diff:"port"`
	Endpoint            string            `json:"endpoint,omitempty" diff:"-"`
	AvailabilityZones   []string          `json:"availability_zones" diff:"-"`
	SecurityGroups      []string          `json:"security_groups" diff:"security_groups"`
	SecurityGroupAWSIDs []string          `json:"security_group_aws_ids" diff:"-"`
	Networks            []string          `json:"networks" diff:"networks"`
	NetworkAWSIDs       []string          `json:"network_aws_ids" diff:"-"`
	DatabaseName        string            `json:"database_name,omitempty" diff:"-"`
	DatabaseUsername    string            `json:"database_username,omitempty" diff:"-"`
	DatabasePassword    string            `json:"database_password,omitempty" diff:"database_password"`
	BackupRetention     *int64            `json:"backup_retention,omitempty" diff:"backup_retention"`
	BackupWindow        string            `json:"backup_window,omitempty" diff:"backup_window"`
	MaintenanceWindow   string            `json:"maintenance_window,omitempty" diff:"maintenance_window"`
	ReplicationSource   string            `json:"replication_source,omitempty" diff:"-"`
	FinalSnapshot       bool              `json:"final_snapshot" diff:"-"`
	Tags                map[string]string `json:"tags" diff:"-"`
	DatacenterType      string            `json:"datacenter_type" diff:"-"`
	DatacenterName      string            `json:"datacenter_name" diff:"-"`
	DatacenterRegion    string            `json:"datacenter_region" diff:"-"`
	AccessKeyID         string            `json:"aws_access_key_id" diff:"-"`
	SecretAccessKey     string            `json:"aws_secret_access_key" diff:"-"`
	Service             string            `json:"service" diff:"-"`
}

// GetID : returns the component's ID
func (r *RDSCluster) GetID() string {
	return r.ComponentID
}

// GetName returns a components name
func (r *RDSCluster) GetName() string {
	return r.Name
}

// GetProvider : returns the provider type
func (r *RDSCluster) GetProvider() string {
	return r.ProviderType
}

// GetProviderID returns a components provider id
func (r *RDSCluster) GetProviderID() string {
	return r.ARN
}

// GetType : returns the type of the component
func (r *RDSCluster) GetType() string {
	return r.ComponentType
}

// GetState : returns the state of the component
func (r *RDSCluster) GetState() string {
	return r.State
}

// SetState : sets the state of the component
func (r *RDSCluster) SetState(s string) {
	r.State = s
}

// GetAction : returns the action of the component
func (r *RDSCluster) GetAction() string {
	return r.Action
}

// SetAction : Sets the action of the component
func (r *RDSCluster) SetAction(s string) {
	r.Action = s
}

// GetGroup : returns the components group
func (r *RDSCluster) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (r *RDSCluster) GetTags() map[string]string {
	return r.Tags
}

// GetTag returns a components tag
func (r *RDSCluster) GetTag(tag string) string {
	return r.Tags[tag]
}

// Diff : diff's the component against another component of the same type
func (r *RDSCluster) Diff(c graph.Component) (diff.Changelog, error) {
	cr, ok := c.(*RDSCluster)
	if ok {
		return diff.Diff(cr, r)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (r *RDSCluster) Update(c graph.Component) {
	cr, ok := c.(*RDSCluster)
	if ok {
		r.ARN = cr.ARN
		r.Endpoint = cr.Endpoint
	}

	r.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (r *RDSCluster) Rebuild(g *graph.Graph) {
	if len(r.Networks) > len(r.NetworkAWSIDs) {
		for _, nw := range r.Networks {
			r.NetworkAWSIDs = append(r.NetworkAWSIDs, templSubnetID(nw))
		}
	}

	if len(r.NetworkAWSIDs) > len(r.Networks) {
		for _, nwid := range r.NetworkAWSIDs {
			nw := g.GetComponents().ByProviderID(nwid)
			if nw != nil {
				r.Networks = append(r.Networks, nw.GetName())
			}
		}
	}

	if len(r.SecurityGroups) > len(r.SecurityGroupAWSIDs) {
		for _, sg := range r.SecurityGroups {
			r.SecurityGroupAWSIDs = append(r.SecurityGroupAWSIDs, templSecurityGroupID(sg))
		}
	}

	if len(r.SecurityGroupAWSIDs) > len(r.SecurityGroups) {
		for _, sgid := range r.SecurityGroupAWSIDs {
			sg := g.GetComponents().ByProviderID(sgid)
			if sg != nil {
				r.SecurityGroups = append(r.SecurityGroups, sg.GetName())
			}
		}
	}

	r.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (r *RDSCluster) Dependencies() []string {
	var deps []string

	for _, sg := range r.SecurityGroups {
		deps = append(deps, TYPESECURITYGROUP+TYPEDELIMITER+sg)
	}

	for _, nw := range r.Networks {
		deps = append(deps, TYPENETWORK+TYPEDELIMITER+nw)
	}

	return deps
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (r *RDSCluster) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (r *RDSCluster) Validate() error {
	if r.Name == "" {
		return errors.New("RDS Cluster name should not be null")
	}

	if len(r.Name) > 255 {
		return errors.New("RDS Cluster name should not exceed 255 characters")
	}

	if r.Engine == "" {
		return errors.New("RDS Cluster engine type should not be null")
	}

	if r.ReplicationSource != "" {
		if len(r.ReplicationSource) < 12 || r.ReplicationSource[:12] != "arn:aws:rds:" {
			return errors.New("RDS Cluster replication source should be a valid amazon resource name (ARN), i.e. 'arn:aws:rds:us-east-1:123456789012:cluster:my-aurora-cluster'")
		}
	}

	if r.DatabaseName == "" {
		return errors.New("RDS Cluster database name should not be null")
	}

	if len(r.DatabaseName) > 64 {
		return errors.New("RDS Cluster database name should not exceed 64 characters")
	}

	for _, c := range r.DatabaseName {
		if unicode.IsLetter(c) != true && unicode.IsNumber(c) != true {
			return errors.New("RDS Cluster database name can only contain alphanumeric characters")
		}
	}

	if r.DatabaseUsername == "" {
		return errors.New("RDS Cluster database username should not be null")
	}

	if len(r.DatabaseUsername) > 16 {
		return errors.New("RDS Cluster database username should not exceed 16 characters")
	}

	if r.DatabasePassword != "" {
		if len(r.DatabasePassword) < 8 || len(r.DatabasePassword) > 41 {
			return errors.New("RDS Cluster database password should be between 8 and 41 characters")
		}

		for _, c := range r.DatabasePassword {
			if unicode.IsSymbol(c) || unicode.IsMark(c) {
				return fmt.Errorf("RDS Cluster database password contains an offending character: '%c'", c)
			}
		}
	}

	if r.Port != nil {
		if *r.Port < 1150 || *r.Port > 65535 {
			return errors.New("RDS Cluster port number should be between 1150 and 65535")
		}
	}

	if r.BackupRetention != nil {
		if *r.BackupRetention < 1 || *r.BackupRetention > 35 {
			return errors.New("RDS Cluster backup retention should be between 1 and 35 days")
		}
	}

	if r.BackupWindow != "" {
		parts := strings.Split(r.BackupWindow, "-")

		err := validateTimeFormat(parts[0])
		if err != nil {
			return errors.New("RDS Cluster backup window: " + err.Error())
		}

		err = validateTimeFormat(parts[1])
		if err != nil {
			return errors.New("RDS Cluster backup window: " + err.Error())
		}
	}

	if mwerr := validateTimeWindow(r.MaintenanceWindow); r.MaintenanceWindow != "" && mwerr != nil {
		return fmt.Errorf("RDS Cluster maintenance window: %s", mwerr.Error())
	}

	return nil

}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (r *RDSCluster) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (r *RDSCluster) SetDefaultVariables() {
	r.ComponentType = TYPERDSCLUSTER
	r.ComponentID = TYPERDSCLUSTER + TYPEDELIMITER + r.Name
	r.ProviderType = PROVIDERTYPE
	r.DatacenterName = DATACENTERNAME
	r.DatacenterType = DATACENTERTYPE
	r.DatacenterRegion = DATACENTERREGION
	r.AccessKeyID = ACCESSKEYID
	r.SecretAccessKey = SECRETACCESSKEY
}
