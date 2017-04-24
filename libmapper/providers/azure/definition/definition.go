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
	Name              string             `json:"name" yaml:"name"`
	Datacenter        string             `json:"datacenter" yaml:"datacenter"`
	ResourceGroups    []ResourceGroup    `json:"resource_groups,omitempty" yaml:"resource_groups,omitempty"`
	NetworkInterfaces []NetworkInterface `json:"network_interfaces,omitempty" yaml:"network_interfaces,omitempty"`
	PublicIPs         []PublicIP         `json:"public_ips,omitempty" yaml:"public_ips,omitempty"`
	SecurityGroups    []SecurityGroup    `json:"security_groups,omitempty" yaml:"security_groups,omitempty"`
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
