/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import "github.com/ernestio/definition-mapper/libmapper/providers/azure/components"

// MapQuery returns a new query
func MapQuery(ctype string, rg string) *components.Query {
	q := &components.Query{
		ComponentType:     ctype,
		Action:            "find",
		ResourceGroupName: rg,
	}

	q.SetDefaultVariables()

	return q
}
