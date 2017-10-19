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

// Diff : handles a diff request
func Diff(r *request.Request) (*graph.Graph, error) {
	p := r.Provider()

	m := providers.NewMapper(p)
	if m == nil {
		return nil, errors.New("could not infer environment provider type")
	}

	fg, err := r.FromMapping(m)
	if err != nil {
		return nil, err
	}

	tg, err := r.ToMapping(m)
	if err != nil {
		return nil, err
	}

	g, err := tg.Diff(fg)
	if err != nil {
		return nil, err
	}

	g.ID = r.ID
	g.Name = r.Name

	return g, err
}
