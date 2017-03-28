/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"strings"

	"github.com/ernestio/definition-mapper/libmapper/providers/aws/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapS3Buckets : Maps the s3 buckets from a given input payload.
func MapS3Buckets(d *definition.Definition) []*components.S3Bucket {
	var s3buckets []*components.S3Bucket

	for _, s3 := range d.S3Buckets {
		s := &components.S3Bucket{
			Name:           s3.Name,
			ACL:            s3.ACL,
			BucketLocation: s3.BucketLocation,
			Tags:           mapTagsServiceOnly(d.Name),
		}

		for _, grantee := range s3.Grantees {
			s.Grantees = append(s.Grantees, components.S3Grantee{
				ID:          grantee.ID,
				Type:        grantee.Type,
				Permissions: strings.ToUpper(grantee.Permissions),
			})
		}

		s.SetDefaultVariables()

		s3buckets = append(s3buckets, s)
	}

	return s3buckets
}

// MapDefinitionS3Buckets : Maps the s3 buckets from the internal format to the input definition format.
func MapDefinitionS3Buckets(g *graph.Graph) []definition.S3 {
	var s3buckets []definition.S3

	for _, gs3 := range g.GetComponents().ByType("s3") {
		s3 := gs3.(*components.S3Bucket)

		s := definition.S3{
			Name:           s3.Name,
			ACL:            s3.ACL,
			BucketLocation: s3.BucketLocation,
		}

		for _, grantee := range s3.Grantees {
			s.Grantees = append(s.Grantees, definition.S3Grantee{
				ID:          grantee.ID,
				Type:        grantee.Type,
				Permissions: strings.ToLower(grantee.Permissions),
			})
		}

		s3buckets = append(s3buckets, s)
	}

	return s3buckets
}
