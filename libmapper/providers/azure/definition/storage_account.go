/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// StorageAccount ...
type StorageAccount struct {
	ID                   string             `json:"id,omitempty" yaml:"id,omitempty"`
	Name                 string             `json:"name,omitempty" yaml:"name,omitempty"`
	AccountType          string             `json:"account_type,omitempty" yaml:"account_type,omitempty"`
	AccountKind          string             `json:"account_kind,omitempty" yaml:"account_kind,omitempty"`
	EnableBlobEncryption bool               `json:"enable_blob_encryption,omitempty" yaml:"enable_blob_encryption,omitempty"`
	Tags                 map[string]string  `json:"tags,omitempty" yaml:"tags,omitempty"`
	Containers           []StorageContainer `json:"containers,omitempty" yaml:"containers,omitempty"`
}
