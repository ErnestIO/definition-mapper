/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// StorageAccount ...
type StorageAccount struct {
	ID                   string             `json:"id" yaml:"id"`
	Name                 string             `json:"name" yaml:"name"`
	AccountType          string             `json:"account_type" yaml:"account_type"`
	AccountKind          string             `json:"account_kind" yaml:"account_kind"`
	EnableBlobEncryption bool               `json:"enable_blob_encryption" yaml:"enable_blob_encryption"`
	Tags                 map[string]string  `json:"tags" yaml:"tags"`
	Containers           []StorageContainer `json:"containers" yaml:"containers"`
}
