/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed With this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// S3Grantee ...
type S3Grantee struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Permissions string `json:"permissions"`
}

// S3 ...
type S3 struct {
	Name           string      `json:"name"`
	ACL            string      `json:"acl,omitempty"`
	BucketLocation string      `json:"bucket_location"`
	Grantees       []S3Grantee `json:"grantees,omitempty"`
}
