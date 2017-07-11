/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"encoding/json"
	"sort"

	"github.com/ernestio/definition-mapper/libmapper/providers/aws/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapIamPolicies ...
func MapIamPolicies(d *definition.Definition) []*components.IamPolicy {
	var ps []*components.IamPolicy

	for _, policy := range d.IamPolicies {
		cp := &components.IamPolicy{
			Name:           policy.Name,
			Path:           policy.Path,
			Description:    policy.Description,
			PolicyDocument: policy.PolicyDocumentRaw,
		}

		if len(policy.PolicyDocument) > 0 {
			data, _ := json.Marshal(policy.PolicyDocument)
			cp.PolicyDocument = string(data)
		}

		cp.SetDefaultVariables()

		ps = append(ps, cp)
	}

	return ps
}

// MapDefinitionIamPolicies : Maps output iam policies into a definition defined iam policies
func MapDefinitionIamPolicies(g *graph.Graph) []definition.IamPolicy {
	var policys []definition.IamPolicy
	var referenced []string

	for _, c := range g.GetComponents().ByType("iam_role") {
		role := c.(*components.IamRole)
		referenced = append(referenced, role.AssumePolicyDocument)
	}

	for _, c := range g.GetComponents().ByType("iam_policy") {
		var policyDoc map[string]interface{}

		r := c.(*components.IamPolicy)

		if sort.SearchStrings(referenced, r.Name) == -1 {
			g.DeleteComponent(c)
			continue
		}

		_ = json.Unmarshal([]byte(r.PolicyDocument), &policyDoc)

		policys = append(policys, definition.IamPolicy{
			Name:           r.Name,
			Path:           r.Path,
			Description:    r.Description,
			PolicyDocument: policyDoc,
		})
	}

	return policys
}
