/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

// Definition ...
type Definition struct {
	Name                string               `json:"name" yaml:"name"`
	Datacenter          string               `json:"datacenter" yaml:"datacenter"`
	Vpcs                []Vpc                `json:"vpcs,omitempty" yaml:"vpcs,omitempty"`
	Networks            []Network            `json:"networks,omitempty" yaml:"networks,omitempty"`
	Instances           []Instance           `json:"instances,omitempty" yaml:"instances,omitempty"`
	SecurityGroups      []SecurityGroup      `json:"security_groups,omitempty" yaml:"security_groups,omitempty"`
	ELBs                []ELB                `json:"loadbalancers,omitempty" yaml:"loadbalancers,omitempty"`
	EBSVolumes          []EBSVolume          `json:"ebs_volumes,omitempty" yaml:"ebs_volumes,omitempty"`
	NatGateways         []NatGateway         `json:"nat_gateways,omitempty" yaml:"nat_gateways,omitempty"`
	RDSClusters         []RDSCluster         `json:"rds_clusters,omitempty" yaml:"rds_clusters,omitempty"`
	RDSInstances        []RDSInstance        `json:"rds_instances,omitempty" yaml:"rds_instances,omitempty"`
	Route53Zones        []Route53Zone        `json:"route53_zones,omitempty" yaml:"route53_zones,omitempty"`
	IamRoles            []IamRole            `json:"iam_roles" yaml:"iam_roles"`
	IamPolicies         []IamPolicy          `json:"iam_policies" yaml:"iam_policies"`
	IamInstanceProfiles []IamInstanceProfile `json:"iam_instance_profiles" yaml:"iam_instance_profiles"`
	S3Buckets           []S3                 `json:"s3_buckets,omitempty" yaml:"s3_buckets,omitempty"`
}

// New returns a new Definition
func New() *Definition {
	return &Definition{}
}

// LoadJSON unmarshals raw json data onto the defintion
func (d *Definition) LoadJSON(data []byte) error {
	return json.Unmarshal(data, d)
}

// LoadMap converts a generic definition from a map[string]interface into an aws definition
func (d *Definition) LoadMap(i map[string]interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   d,
		TagName:  "json",
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(i)
}
