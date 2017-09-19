/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
	"github.com/r3labs/graph"
)

// MapIamInstanceProfiles ...
func MapIamInstanceProfiles(d *definition.Definition) []*components.IamInstanceProfile {
	var ps []*components.IamInstanceProfile

	for _, profile := range d.IamInstanceProfiles {
		cp := &components.IamInstanceProfile{
			Name:  profile.Name,
			Path:  profile.Path,
			Roles: profile.Roles,
		}

		cp.SetDefaultVariables()

		ps = append(ps, cp)
	}

	return ps
}

// MapDefinitionIamInstanceProfiles : Maps output iam instance profiles into a definition defined iam instance profiles
func MapDefinitionIamInstanceProfiles(g *graph.Graph) []definition.IamInstanceProfile {
	var profiles []definition.IamInstanceProfile

	for _, c := range g.GetComponents().ByType("iam_instance_profile") {
		r := c.(*components.IamInstanceProfile)

		profiles = append(profiles, definition.IamInstanceProfile{
			Name:  r.Name,
			Path:  r.Path,
			Roles: r.Roles,
		})
	}

	return profiles
}
