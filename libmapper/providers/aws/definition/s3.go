/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed With this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// S3Grantee ...
type S3Grantee struct {
	ID          string `json:"id" yaml:"id"`
	Type        string `json:"type" yaml:"type"`
	Permissions string `json:"permissions" yaml:"permissions"`
}

// S3 ...
type S3 struct {
	Name           string      `json:"name" yaml:"name"`
	ACL            string      `json:"acl,omitempty" yaml:"acl,omitempty"`
	BucketLocation string      `json:"bucket_location" yaml:"bucket_location"`
	Grantees       []S3Grantee `json:"grantees,omitempty" yaml:"grantees,omitempty"`
}
