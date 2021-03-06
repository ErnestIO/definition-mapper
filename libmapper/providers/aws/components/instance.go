/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"

	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

// InstanceVolume ...
type InstanceVolume struct {
	VolumeAWSID string `json:"volume_aws_id" diff:"-"`
	Volume      string `json:"volume" diff:"volume"`
	Device      string `json:"device" diff:"device"`
}

// Instance : mapping of an instance component
type Instance struct {
	ProviderType          string            `json:"_provider" diff:"-"`
	ComponentType         string            `json:"_component" diff:"-"`
	ComponentID           string            `json:"_component_id" diff:"_component_id,immutable"`
	State                 string            `json:"_state" diff:"-"`
	Action                string            `json:"_action" diff:"-"`
	InstanceAWSID         string            `json:"instance_aws_id" diff:"-"`
	Name                  string            `json:"name" diff:"-"`
	Type                  string            `json:"instance_type" diff:"instance_type"`
	Image                 string            `json:"image" diff:"image,immutable"`
	IP                    string            `json:"ip" diff:"ip,immutable"`
	PublicIP              string            `json:"public_ip" diff:"public_ip,immutable"`
	ElasticIP             string            `json:"elastic_ip" diff:"elastic_ip,immutable"`
	ElasticIPAWSID        *string           `json:"elastic_ip_aws_id,omitempty" diff:"-"`
	AssignElasticIP       bool              `json:"assign_elastic_ip" diff:"-"`
	KeyPair               string            `json:"key_pair" diff:"-"`
	UserData              string            `json:"user_data" diff:"-"`
	Network               string            `json:"network_name" diff:"network,immutable"`
	NetworkAWSID          string            `json:"network_aws_id" diff:"-"`
	NetworkIsPublic       bool              `json:"network_is_public" diff:"-"`
	SecurityGroups        []string          `json:"security_groups" diff:"security_groups"`
	SecurityGroupAWSIDs   []string          `json:"security_group_aws_ids" diff:"-"`
	IAMInstanceProfile    *string           `json:"iam_instance_profile" diff:"-"`
	IAMInstanceProfileARN *string           `json:"iam_instance_profile_arn" diff:"-"`
	Volumes               []InstanceVolume  `json:"volumes" diff:"volumes"`
	Tags                  map[string]string `json:"tags" diff:"tags"`
	Powered               bool              `json:"powered" diff:"powered"`
	DatacenterType        string            `json:"datacenter_type,omitempty" diff:"-"`
	DatacenterName        string            `json:"datacenter_name,omitempty" diff:"-"`
	DatacenterRegion      string            `json:"datacenter_region" diff:"-"`
	AccessKeyID           string            `json:"aws_access_key_id" diff:"-"`
	SecretAccessKey       string            `json:"aws_secret_access_key" diff:"-"`
	Service               string            `json:"service" diff:"-"`
}

// GetID : returns the component's ID
func (i *Instance) GetID() string {
	return i.ComponentID
}

// GetName returns a components name
func (i *Instance) GetName() string {
	return i.Name
}

// GetProvider : returns the provider type
func (i *Instance) GetProvider() string {
	return i.ProviderType
}

// GetProviderID returns a components provider id
func (i *Instance) GetProviderID() string {
	return i.InstanceAWSID
}

// GetType : returns the type of the component
func (i *Instance) GetType() string {
	return i.ComponentType
}

// GetState : returns the state of the component
func (i *Instance) GetState() string {
	return i.State
}

// SetState : sets the state of the component
func (i *Instance) SetState(s string) {
	i.State = s
}

// GetAction : returns the action of the component
func (i *Instance) GetAction() string {
	return i.Action
}

// SetAction : Sets the action of the component
func (i *Instance) SetAction(s string) {
	i.Action = s
}

// GetGroup : returns the components group
func (i *Instance) GetGroup() string {
	return i.Tags["ernest.instance_group"]
}

// GetTags returns a components tags
func (i *Instance) GetTags() map[string]string {
	return i.Tags
}

// GetTag returns a components tag
func (i *Instance) GetTag(tag string) string {
	return i.Tags[tag]
}

// Diff : diff's the component against another component of the same type
func (i *Instance) Diff(c graph.Component) (diff.Changelog, error) {
	ci, ok := c.(*Instance)
	if ok {
		return diff.Diff(ci, i)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (i *Instance) Update(c graph.Component) {
	ci, ok := c.(*Instance)
	if ok {
		i.InstanceAWSID = ci.InstanceAWSID
		i.PublicIP = ci.PublicIP
		i.ElasticIP = ci.ElasticIP
		i.ElasticIPAWSID = ci.ElasticIPAWSID
		i.Powered = ci.Powered
	}

	i.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (i *Instance) Rebuild(g *graph.Graph) {
	if i.Network == "" && i.NetworkAWSID != "" {
		n := g.GetComponents().ByProviderID(i.NetworkAWSID)
		if n != nil {
			i.Network = n.GetName()
		}
	}

	if i.Network != "" && i.NetworkAWSID == "" {
		i.NetworkAWSID = templSubnetID(i.Network)
	}

	if len(i.SecurityGroups) > len(i.SecurityGroupAWSIDs) {
		for _, sg := range i.SecurityGroups {
			i.SecurityGroupAWSIDs = append(i.SecurityGroupAWSIDs, templSecurityGroupID(sg))
		}
	}

	if len(i.SecurityGroupAWSIDs) > len(i.SecurityGroups) {
		for _, sgid := range i.SecurityGroupAWSIDs {
			sg := g.GetComponents().ByProviderID(sgid)
			if sg != nil {
				i.SecurityGroups = append(i.SecurityGroups, sg.GetName())
			}
		}
	}

	for x := 0; x < len(i.Volumes); x++ {
		if i.Volumes[x].Volume == "" && i.Volumes[x].VolumeAWSID != "" {
			v := g.GetComponents().ByProviderID(i.Volumes[x].VolumeAWSID)
			if v != nil {
				i.Volumes[x].Volume = v.GetName()
			}
		}

		if i.Volumes[x].VolumeAWSID == "" && i.Volumes[x].Volume != "" {
			i.Volumes[x].VolumeAWSID = templEBSVolumeID(i.Volumes[x].Volume)
		}
	}

	if i.IAMInstanceProfile == nil && i.IAMInstanceProfileARN != nil {
		p := g.GetComponents().ByProviderID(*i.IAMInstanceProfileARN)
		if p != nil {
			name := p.GetName()
			i.IAMInstanceProfile = &name
		}
	}

	if i.IAMInstanceProfileARN == nil && i.IAMInstanceProfile != nil {
		tpl := templIAMInstanceProfileARN(*i.IAMInstanceProfile)
		i.IAMInstanceProfileARN = &tpl
	}

	i.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (i *Instance) Dependencies() []string {
	var deps []string

	for _, sg := range i.SecurityGroups {
		deps = append(deps, TYPESECURITYGROUP+TYPEDELIMITER+sg)
	}

	for _, ebs := range i.Volumes {
		deps = append(deps, TYPEEBSVOLUME+TYPEDELIMITER+ebs.Volume)
	}

	if i.IAMInstanceProfile != nil {
		deps = append(deps, TYPEIAMINSTANCEPROFILE+TYPEDELIMITER+*i.IAMInstanceProfile)
	}

	deps = append(deps, TYPENETWORK+TYPEDELIMITER+i.Network)

	return deps
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (i *Instance) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (i *Instance) Validate() error {
	if i.Name == "" {
		return errors.New("Instance name should not be null")
	}

	if i.Type == "" {
		return errors.New("Instance type should not be null")
	}

	if i.Image == "" {
		return errors.New("Instance image should not be null")
	}

	if i.Network == "" {
		return errors.New("Instance network should not be null")
	}

	if len(i.SecurityGroups) != len(i.SecurityGroupAWSIDs) {
		return errors.New("Instance security groups are incorrect")
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (i *Instance) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (i *Instance) SetDefaultVariables() {
	i.ComponentType = TYPEINSTANCE
	i.ComponentID = TYPEINSTANCE + TYPEDELIMITER + i.Name
	i.ProviderType = PROVIDERTYPE
	i.DatacenterName = DATACENTERNAME
	i.DatacenterType = DATACENTERTYPE
	i.DatacenterRegion = DATACENTERREGION
	i.AccessKeyID = ACCESSKEYID
	i.SecretAccessKey = SECRETACCESSKEY
}

func hasVolume(vols []InstanceVolume, volume string) bool {
	for _, v := range vols {
		if v.Volume == volume {
			return true
		}
	}

	return false
}
