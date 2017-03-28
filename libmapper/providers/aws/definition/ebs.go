/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// EBSVolume ...
type EBSVolume struct {
	Name             string  `json:"name" yaml:"name"`
	Type             string  `json:"type" yaml:"type"`
	Size             *int64  `json:"size" yaml:"size"`
	Iops             *int64  `json:"iops" yaml:"iops"`
	Count            int     `json:"count" yaml:"count"`
	Encrypted        bool    `json:"encrypted" yaml:"encrypted"`
	EncryptionKeyID  *string `json:"encryption_key_id" yaml:"encryption_key_id"`
	AvailabilityZone string  `json:"availability_zone" yaml:"availability_zone"`
}
