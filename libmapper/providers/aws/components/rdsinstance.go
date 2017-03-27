/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	graph "gopkg.in/r3labs/graph.v2"
)

// Licenses stores all valid license types for rds
var Licenses = []string{"license-included", "bring-your-own-license", "general-public-license"}

// StorageTypes stores all of the valid types of storage that can be allocated to a RDS Instance
var StorageTypes = []string{"standard", "gp2", "io1"}

// EngineTypeAurora ...
var EngineTypeAurora = "aurora"

// RDSInstance ...
type RDSInstance struct {
	ProviderType        string            `json:"_provider"`
	ComponentType       string            `json:"_component"`
	ComponentID         string            `json:"_component_id"`
	State               string            `json:"_state"`
	Action              string            `json:"_action"`
	ARN                 string            `json:"arn"`
	Name                string            `json:"name"`
	Size                string            `json:"size"`
	Engine              string            `json:"engine"`
	EngineVersion       string            `json:"engine_version,omitempty"`
	Port                *int64            `json:"port,omitempty"`
	Cluster             string            `json:"cluster,omitempty"`
	Public              bool              `json:"public"`
	Endpoint            string            `json:"endpoint,omitempty"`
	MultiAZ             bool              `json:"multi_az"`
	PromotionTier       *int64            `json:"promotion_tier,omitempty"`
	StorageType         string            `json:"storage_type,omitempty"`
	StorageSize         *int64            `json:"storage_size,omitempty"`
	StorageIops         *int64            `json:"storage_iops,omitempty"`
	AvailabilityZone    string            `json:"availability_zone,omitempty"`
	SecurityGroups      []string          `json:"security_groups"`
	SecurityGroupAWSIDs []string          `json:"security_group_aws_ids"`
	Networks            []string          `json:"networks"`
	NetworkAWSIDs       []string          `json:"network_aws_ids"`
	DatabaseName        string            `json:"database_name,omitempty"`
	DatabaseUsername    string            `json:"database_username,omitempty"`
	DatabasePassword    string            `json:"database_password,omitempty"`
	AutoUpgrade         bool              `json:"auto_upgrade"`
	BackupRetention     *int64            `json:"backup_retention,omitempty"`
	BackupWindow        string            `json:"backup_window,omitempty"`
	MaintenanceWindow   string            `json:"maintenance_window,omitempty"`
	FinalSnapshot       bool              `json:"final_snapshot"`
	ReplicationSource   string            `json:"replication_source,omitempty"`
	License             string            `json:"license,omitempty"`
	Timezone            string            `json:"timezone,omitempty"`
	Tags                map[string]string `json:"tags"`
	DatacenterType      string            `json:"datacenter_type"`
	DatacenterName      string            `json:"datacenter_name"`
	DatacenterRegion    string            `json:"datacenter_region"`
	AccessKeyID         string            `json:"aws_access_key_id"`
	SecretAccessKey     string            `json:"aws_secret_access_key"`
	Service             string            `json:"service"`
}

// GetID : returns the component's ID
func (r *RDSInstance) GetID() string {
	return r.ComponentID
}

// GetName returns a components name
func (r *RDSInstance) GetName() string {
	return r.Name
}

// GetProvider : returns the provider type
func (r *RDSInstance) GetProvider() string {
	return r.ProviderType
}

// GetProviderID returns a components provider id
func (r *RDSInstance) GetProviderID() string {
	return r.ARN
}

// GetType : returns the type of the component
func (r *RDSInstance) GetType() string {
	return r.ComponentType
}

// GetState : returns the state of the component
func (r *RDSInstance) GetState() string {
	return r.State
}

// SetState : sets the state of the component
func (r *RDSInstance) SetState(s string) {
	r.State = s
}

// GetAction : returns the action of the component
func (r *RDSInstance) GetAction() string {
	return r.Action
}

// SetAction : Sets the action of the component
func (r *RDSInstance) SetAction(s string) {
	r.Action = s
}

// GetGroup : returns the components group
func (r *RDSInstance) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (r *RDSInstance) GetTags() map[string]string {
	return r.Tags
}

// GetTag returns a components tag
func (r *RDSInstance) GetTag(tag string) string {
	return r.Tags[tag]
}

// Diff : diff's the component against another component of the same type
func (r *RDSInstance) Diff(c graph.Component) bool {
	cr, ok := c.(*RDSInstance)
	if ok {
		if r.Size != cr.Size {
			return true
		}

		if r.EngineVersion != cr.EngineVersion {
			return true
		}

		if r.Port != nil && cr.Port != nil {
			if *r.Port != *cr.Port {
				return true
			}
		}

		if r.StorageSize != nil && cr.StorageSize != nil {
			if *r.StorageSize != *cr.StorageSize {
				return true
			}
		}

		if r.StorageIops != nil && cr.StorageIops != nil {
			if *r.StorageIops != *cr.StorageIops {
				return true
			}
		}

		if r.StorageType != cr.StorageType {
			return true
		}

		if r.MultiAZ != cr.MultiAZ {
			return true
		}

		if r.PromotionTier != nil && cr.PromotionTier != nil {
			if *r.PromotionTier != *cr.PromotionTier {
				return true
			}
		}

		if r.AutoUpgrade != cr.AutoUpgrade {
			return true
		}

		if r.BackupRetention != nil && cr.BackupRetention != nil {
			if *r.BackupRetention != *cr.BackupRetention {
				return true
			}
		}

		if r.BackupWindow != cr.BackupWindow {
			return true
		}

		if r.DatabasePassword != cr.DatabasePassword {
			return true
		}

		if r.Public != cr.Public {
			return true
		}

		if reflect.DeepEqual(r.SecurityGroups, cr.SecurityGroups) != true {
			return true
		}

		return !reflect.DeepEqual(r.Networks, cr.Networks)
	}

	return false
}

// Update : updates the provider returned values of a component
func (r *RDSInstance) Update(c graph.Component) {
	cr, ok := c.(*RDSInstance)
	if ok {
		r.ARN = cr.ARN
		r.Endpoint = cr.Endpoint
	}

	r.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (r *RDSInstance) Rebuild(g *graph.Graph) {
	if r.Cluster != "" {
		r.DatabaseName = ""
		r.DatabaseUsername = ""
		r.EngineVersion = ""
		r.StorageType = ""
		r.StorageIops = nil
		r.StorageSize = nil
		r.Port = nil
	}

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
func (r *RDSInstance) Dependencies() []string {
	var deps []string

	for _, sg := range r.SecurityGroups {
		deps = append(deps, TYPESECURITYGROUP+TYPEDELIMITER+sg)
	}

	for _, nw := range r.Networks {
		deps = append(deps, TYPENETWORK+TYPEDELIMITER+nw)
	}

	if r.Cluster != "" {
		deps = append(deps, TYPERDSCLUSTER+TYPEDELIMITER+r.Cluster)
	}

	return deps
}

// Validate : validates the components values
func (r *RDSInstance) Validate() error {
	if r.Name == "" {
		return errors.New("RDS Instance name should not be null")
	}

	if len(r.Name) > 255 {
		return errors.New("RDS Instance name should not exceed 255 characters")
	}

	if r.Size == "" {
		return errors.New("RDS Instance size should not be null")
	}

	if r.Size[:3] != "db." {
		return errors.New("RDS Instance size should be a valid resource size. i.e. 'db.r3.large'")
	}

	err := r.validateReplication()
	if err != nil {
		return err
	}

	err = r.validateDatabase()
	if err != nil {
		return err
	}

	err = r.validateEngine()
	if err != nil {
		return err
	}

	err = r.validatePort()
	if err != nil {
		return err
	}

	err = r.validateStorage()
	if err != nil {
		return err
	}

	err = r.validateBackups()
	if err != nil {
		return err
	}

	err = r.validateOther()
	if err != nil {
		return err
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (r *RDSInstance) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (r *RDSInstance) SetDefaultVariables() {
	r.ComponentType = TYPERDSINSTANCE
	r.ComponentID = TYPERDSINSTANCE + TYPEDELIMITER + r.Name
	r.ProviderType = PROVIDERTYPE
	r.DatacenterName = DATACENTERNAME
	r.DatacenterType = DATACENTERTYPE
	r.DatacenterRegion = DATACENTERREGION
	r.AccessKeyID = ACCESSKEYID
	r.SecretAccessKey = SECRETACCESSKEY
}

func (r *RDSInstance) validateReplication() error {
	if r.ReplicationSource != "" {
		if r.Engine != "" {
			return errors.New("RDS Instance must not specify an engine if a replication source is set")
		}

		if r.EngineVersion != "" {
			return errors.New("RDS Instance must not specify an engine version if a replication source is set")
		}

		if r.StorageSize != nil {
			return errors.New("RDS Instance must not specify storage size if a replication source is set")
		}

		if r.Cluster != "" {
			return errors.New("RDS Instance must not specify a cluster if a replication source is set")
		}

		if r.MultiAZ == true {
			return errors.New("RDS Instance must not specify multi az standby instance if a replication source is set")
		}

		if r.PromotionTier != nil {
			return errors.New("RDS Instance must not specify promotion tier if a replication source is set")
		}

		if r.DatabaseName != "" {
			return errors.New("RDS Instance must not specify database name if a replication source is set")
		}

		if r.DatabaseUsername != "" {
			return errors.New("RDS Instance must not specify database username if a replication source is set")
		}

		if r.DatabasePassword != "" {
			return errors.New("RDS Instance must not specify database password if a replication source is set")
		}

		if r.License != "" {
			return errors.New("RDS Instance must not specify a license type if a replication source is set")
		}

		if r.Timezone != "" {
			return errors.New("RDS Instance must not specify a timezone if a replication source is set")
		}
	}

	return nil
}

func (r *RDSInstance) validateBackups() error {
	if r.BackupRetention != nil {
		if *r.BackupRetention < 1 || *r.BackupRetention > 35 {
			return errors.New("RDS Instance backup retention should be between 1 and 35 days")
		}
	}

	if r.BackupWindow != "" {
		parts := strings.Split(r.BackupWindow, "-")

		err := validateTimeFormat(parts[0])
		if err != nil {
			return errors.New("RDS Instance backup window: " + err.Error())
		}

		err = validateTimeFormat(parts[1])
		if err != nil {
			return errors.New("RDS Instance backup window: " + err.Error())
		}
	}

	return nil
}

func (r *RDSInstance) validatePort() error {
	if r.Cluster != "" && r.Port != nil {
		return fmt.Errorf("RDS Instance port should be set on cluster")
	}

	if r.Port != nil {
		if *r.Port < 1150 || *r.Port > 65535 {
			return errors.New("RDS Instance port number should be between 1150 and 65535")
		}
	}

	return nil
}

func (r *RDSInstance) validateDatabase() error {
	if r.Cluster != "" {
		if r.DatabaseName != "" {
			return errors.New("RDS Instance database name should be set on cluster")
		}

		if r.DatabaseUsername != "" {
			return errors.New("RDS Instance database username should be set on cluster")
		}

		if r.DatabasePassword != "" {
			return errors.New("RDS Instance database password should be set on cluster")
		}
	} else {
		if r.DatabaseName == "" {
			return errors.New("RDS Instance database name should not be null")
		}

		if len(r.DatabaseName) > 64 {
			return errors.New("RDS Instance database name should not exceed 64 characters")
		}

		for _, c := range r.DatabaseName {
			if unicode.IsLetter(c) != true && unicode.IsNumber(c) != true {
				return errors.New("RDS Instance database name can only contain alphanumeric characters")
			}
		}

		if r.DatabaseUsername == "" {
			return errors.New("RDS Instance database username should not be null")
		}

		if len(r.DatabaseUsername) > 16 {
			return errors.New("RDS Instance database username should not exceed 16 characters")
		}

		for _, c := range r.DatabasePassword {
			if unicode.IsSymbol(c) || unicode.IsMark(c) {
				return fmt.Errorf("RDS Instance database password contains an offending character: '%c'", c)
			}
		}
	}

	return nil
}

func (r *RDSInstance) validateEngine() error {
	if r.Cluster != "" {
		if r.EngineVersion != "" {
			return fmt.Errorf("RDS Instance engine version should be set on cluster")
		}
	} else {
		if r.Engine == "" {
			return errors.New("RDS Instance engine type should not be null")
		}
	}

	return nil
}

func (r *RDSInstance) validateStorage() error {
	if r.Engine != EngineTypeAurora {
		if r.StorageType != "" && isOneOf(StorageTypes, r.StorageType) != true {
			return errors.New("RDS Instance storage type must be either 'standard', 'gp2' or 'io1'")
		}
		if r.StorageSize != nil {
			if *r.StorageSize < 5 || *r.StorageSize > 6144 {
				return errors.New("RDS Instance storage size must be between 5 - 6144 GB")
			}
		}
		if r.StorageIops != nil {
			if (*r.StorageIops % 1000) != 0 {
				return errors.New("RDS Instance storage iops must be a multiple of 1000")
			}
		}
	} else {
		if r.StorageType != "" || r.StorageSize != nil || r.StorageIops != nil {
			return errors.New("RDS Instance storage options cannot be set if the engine type is 'aurora'")
		}
	}

	return nil
}

func (r *RDSInstance) validateOther() error {
	if r.PromotionTier != nil {
		if r.Engine != EngineTypeAurora {
			return errors.New("RDS Instance promotion tier should only be specified when using the aurora engine")
		}
		if *r.PromotionTier < 0 || *r.PromotionTier > 15 {
			return errors.New("RDS Instance promotion tier should be between 0 - 15")
		}
	}

	if r.AvailabilityZone != "" && r.MultiAZ {
		return errors.New("RDS Instance cannot specify both an availability zone and a multi az standby instance")
	}

	if mwerr := validateTimeWindow(r.MaintenanceWindow); r.MaintenanceWindow != "" && mwerr != nil {
		return fmt.Errorf("RDS Instance maintenance window: %s", mwerr.Error())
	}

	if r.Public == false && len(r.Networks) < 1 && r.Cluster == "" {
		return errors.New("RDS Instance should specify at least one network if not set to public")
	}

	if r.Engine != EngineTypeAurora && r.Engine != "" && isOneOf(Licenses, r.License) != true {
		return errors.New("RDS Instance license must be one of 'license-included', 'bring-your-own-license', 'general-public-license'")
	}

	return nil
}
