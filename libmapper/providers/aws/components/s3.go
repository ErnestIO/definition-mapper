/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"fmt"
	"strings"

	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

var (
	// S3GRANTEETYPES : s3 supported grantee types
	S3GRANTEETYPES = []string{"id", "emailaddress", "uri", "canonicaluser", "AmazonCustomerByEmail", "CanonicalUser", "Group"}
	// S3PERMISSIONTYPES : s3 supported permission types
	S3PERMISSIONTYPES = []string{"FULL_CONTROL", "WRITE", "WRITE_ACP", "READ", "READ_ACP"}
	// S3ACLTYPES : s3 supported acl types
	S3ACLTYPES = []string{"private", "public-read", "public-read-write", "aws-exec-read", "authenticated-read", "log-delivery-write"}
)

// S3Grantee ...
type S3Grantee struct {
	ID          string `json:"id" diff:"id,identifier"`
	Type        string `json:"type" diff:"type"`
	Permissions string `json:"permissions" diff:"permissions"`
}

// S3Bucket : Mapping of an s3 bucket component
type S3Bucket struct {
	ProviderType     string            `json:"_provider" diff:"-"`
	ComponentType    string            `json:"_component" diff:"-"`
	ComponentID      string            `json:"_component_id" diff:"component_id,identifier"`
	State            string            `json:"_state" diff:"-"`
	Action           string            `json:"_action" diff:"-"`
	Name             string            `json:"name" diff:"-"`
	ACL              string            `json:"acl" diff:"acl"`
	BucketLocation   string            `json:"bucket_location" diff:"-"`
	BucketURI        string            `json:"bucket_uri" diff:"-"`
	Grantees         []S3Grantee       `json:"grantees,omitempty" diff:"grantees"`
	Tags             map[string]string `json:"tags" diff:"-"`
	DatacenterType   string            `json:"datacenter_type,omitempty" diff:"-"`
	DatacenterName   string            `json:"datacenter_name,omitempty" diff:"-"`
	DatacenterRegion string            `json:"datacenter_region" diff:"-"`
	AccessKeyID      string            `json:"aws_access_key_id" diff:"-"`
	SecretAccessKey  string            `json:"aws_secret_access_key" diff:"-"`
	Service          string            `json:"service" diff:"-"`
}

// GetID : returns the component's ID
func (s3 *S3Bucket) GetID() string {
	return s3.ComponentID
}

// GetName returns a components name
func (s3 *S3Bucket) GetName() string {
	return s3.Name
}

// GetProvider : returns the provider type
func (s3 *S3Bucket) GetProvider() string {
	return s3.ProviderType
}

// GetProviderID returns a components provider id
func (s3 *S3Bucket) GetProviderID() string {
	return s3.Name
}

// GetType : returns the type of the component
func (s3 *S3Bucket) GetType() string {
	return s3.ComponentType
}

// GetState : returns the state of the component
func (s3 *S3Bucket) GetState() string {
	return s3.State
}

// SetState : sets the state of the component
func (s3 *S3Bucket) SetState(s string) {
	s3.State = s
}

// GetAction : returns the action of the component
func (s3 *S3Bucket) GetAction() string {
	return s3.Action
}

// SetAction : Sets the action of the component
func (s3 *S3Bucket) SetAction(s string) {
	s3.Action = s
}

// GetGroup : returns the components group
func (s3 *S3Bucket) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (s3 *S3Bucket) GetTags() map[string]string {
	return s3.Tags
}

// GetTag returns a components tag
func (s3 *S3Bucket) GetTag(tag string) string {
	return s3.Tags[tag]
}

// Diff : diff's the component against another component of the same type
func (s3 *S3Bucket) Diff(c graph.Component) (diff.Changelog, error) {
	cs3, ok := c.(*S3Bucket)
	if ok {
		return diff.Diff(cs3, s3)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (s3 *S3Bucket) Update(c graph.Component) {
	cs3, ok := c.(*S3Bucket)
	if ok {
		s3.BucketURI = cs3.BucketURI
	}

	s3.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (s3 *S3Bucket) Rebuild(g *graph.Graph) {
	s3.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (s3 *S3Bucket) Dependencies() []string {
	return []string{}
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (s3 *S3Bucket) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (s3 *S3Bucket) Validate() error {
	if s3.Name == "" {
		return errors.New("S3 bucket name should not be null")
	}

	if s3.BucketLocation == "" {
		return errors.New("S3 bucket location should not be null")
	}

	if s3.ACL != "" && len(s3.Grantees) > 0 {
		return errors.New("S3 bucket must specify either acl or grantees, not both")
	}

	if s3.ACL != "" && isOneOf(S3ACLTYPES, s3.ACL) == false {
		return fmt.Errorf("S3 bucket ACL (%s) is not valid. Must be one of [%s]", s3.ACL, strings.Join(S3ACLTYPES, " | "))
	}

	for _, g := range s3.Grantees {
		if isOneOf(S3GRANTEETYPES, g.Type) == false {
			return fmt.Errorf("S3 grantee type (%s) is invalid", g.Type)
		}

		if g.ID == "" {
			return fmt.Errorf("S3 grantee id should not be null")
		}

		if isOneOf(S3PERMISSIONTYPES, g.Permissions) == false {
			return fmt.Errorf("S3 grantee permissions (%s) is not valid. Must be one of [%s]", s3.ACL, strings.ToLower(strings.Join(S3PERMISSIONTYPES, " | ")))
		}
	}
	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (s3 *S3Bucket) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (s3 *S3Bucket) SetDefaultVariables() {
	s3.ComponentType = TYPES3BUCKET
	s3.ComponentID = TYPES3BUCKET + TYPEDELIMITER + s3.Name
	s3.ProviderType = PROVIDERTYPE
	s3.DatacenterName = DATACENTERNAME
	s3.DatacenterType = DATACENTERTYPE
	s3.DatacenterRegion = DATACENTERREGION
	s3.AccessKeyID = ACCESSKEYID
	s3.SecretAccessKey = SECRETACCESSKEY
}
