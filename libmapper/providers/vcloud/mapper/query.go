/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import "github.com/ernestio/definition-mapper/libmapper/providers/vcloud/components"

// MapQuery returns a new query
func MapQuery(ctype string, values map[string]string) *components.Query {
	q := &components.Query{
		Tags: values,
	}

	q.Action = "find"
	q.ComponentType = ctype

	q.SetDefaultVariables()

	return q
}
