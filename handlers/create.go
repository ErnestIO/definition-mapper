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

// Create : handles a create request
func Create(r *request.Request) (*graph.Graph, error) {
	var g *graph.Graph

	p := r.Provider()

	m := providers.NewMapper(p)
	if m == nil {
		return nil, errors.New("could not infer environment provider type")
	}

	dg, err := r.DefinitionToGraph(m)
	if err != nil {
		return nil, err
	}

	if r.Changelog {
		g, err = dg.DiffWithChangelog(graph.New())
	} else {
		g, err = dg.Diff(graph.New())
	}

	if err != nil {
		return nil, err
	}

	g.ID = r.ID
	g.Name = r.Name
	g.UserID = r.UserID
	g.Username = r.Username

	return g, err
}
