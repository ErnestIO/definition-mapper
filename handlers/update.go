/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package handlers

import (
	"errors"

	"github.com/ernestio/definition-mapper/libmapper/providers"
	"github.com/ernestio/definition-mapper/request"
	"github.com/r3labs/graph"
)

// Update : handles a update request
func Update(r *request.Request) (*graph.Graph, error) {
	p := r.Provider()

	m := providers.NewMapper(p)
	if m == nil {
		return nil, errors.New("could not infer environment provider type")
	}

	dg, err := r.DefinitionToGraph(m)
	if err != nil {
		return nil, err
	}

	fg, err := r.FromMapping(m)
	if err != nil {
		return nil, err
	}

	for _, c := range dg.Components {
		oc := fg.Component(c.GetID())
		if oc != nil {
			c.Update(oc)
		}
	}

	return dg.Diff(fg)
}
