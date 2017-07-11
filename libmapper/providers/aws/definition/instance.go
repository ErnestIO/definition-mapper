/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// InstanceVolume ...
type InstanceVolume struct {
	Volume string `json:"volume" yaml:"volume"`
	Device string `json:"device" yaml:"device"`
}

// Instance ...
type Instance struct {
	Name           string           `json:"name" yaml:"name"`
	Type           string           `json:"type" yaml:"type"`
	Image          string           `json:"image" yaml:"image"`
	Count          int              `json:"count" yaml:"count"`
	Network        string           `json:"network" yaml:"network"`
	StartIP        string           `json:"start_ip" yaml:"start_ip"`
	KeyPair        string           `json:"key_pair" yaml:"key_pair"`
	ElasticIP      bool             `json:"elastic_ip" yaml:"elastic_ip"`
	SecurityGroups []string         `json:"security_groups" yaml:"security_groups"`
	Volumes        []InstanceVolume `json:"volumes" yaml:"volumes"`
	UserData       string           `json:"user_data" yaml:"user_data"`
	IamProfile     *string          `json:"iam_profile" yaml:"iam_profile"`
}
