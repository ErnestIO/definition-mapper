/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// StorageAccount ...
type StorageAccount struct {
	Name                 string            `json:"name" yaml:"name"`
	AccountType          string            `json:"account_type" yaml:"account_type"`
	AccessTier           string            `json:"access_tier" yaml:"access_tier"`
	AccountKind          string            `json:"account_kind" yaml:"account_kind"`
	EnableBlobEncryption bool              `json:"enable_blob_encryption" yaml:"enable_blob_encryption"`
	Tags                 map[string]string `json:"tags" yaml:"tags"`
	Containers           []Container       `json:"containers" yaml:"containers"`
}

// Container ...
type Container struct {
	Name       string `json:"name" yaml:"name"`
	AccessType string `json:"access_type" yaml:"access_type"`
}
