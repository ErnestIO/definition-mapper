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

// Delete : handles a delete request
func Delete(r *request.Request) (*graph.Graph, error) {
	var g *graph.Graph

	p := r.Provider()

	m := providers.NewMapper(p)
	if m == nil {
		return nil, errors.New("could not infer environment provider type")
	}

	original, err := r.FromMapping(m)
	if err != nil {
		return nil, err
	}

	creds := m.ProviderCredentials(r.Credentials)
	original.UpdateComponent(creds)
	for i := range original.Components {
		original.Components[i].Rebuild(original)
	}

	empty := graph.New()

	if r.Changelog {
		g, err = empty.DiffWithChangelog(original)
	} else {
		g, err = empty.Diff(original)
	}

	if err != nil {
		return nil, err
	}

	g.ID = r.ID
	g.Name = r.Name
	g.UserID = r.UserID
	g.Username = r.Username

	return g, nil
}
