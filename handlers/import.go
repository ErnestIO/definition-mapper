/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package handlers

import (
	"strings"

	"github.com/ernestio/definition-mapper/libmapper/providers"
	aws "github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
	azure "github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	"github.com/ernestio/definition-mapper/request"
	"github.com/r3labs/graph"
	yaml "gopkg.in/yaml.v2"
)

type build struct {
	ID         string       `json:"id"`
	Definition string       `json:"definition"`
	Mapping    *graph.Graph `json:"mapping"`
}

// Import : handles a import request
func Import(r *request.Request) (*graph.Graph, error) {
	p := r.Provider()

	m := providers.NewMapper(p)

	filters := r.ImportFilters()

	ig := m.CreateImportGraph(filters)
	ig.ID = r.ID
	ig.Name = r.Name

	c := m.ProviderCredentials(r.Credentials)
	err := ig.AddComponent(c)
	if err != nil {
		return nil, err
	}

	g, err := ig.Diff(graph.New())
	if err != nil {
		return nil, err
	}

	g.ID = r.ID
	g.Name = r.Name

	return g, nil
}

// ImportComplete : handles the conversion of an import graph to a definition
func ImportComplete(ig map[string]interface{}) (interface{}, error) {
	g := graph.New()
	err := g.Load(ig)
	if err != nil {
		return nil, err
	}

	credentials := g.GetComponents().ByType("credentials")
	provider := credentials[0].GetProvider()

	m := providers.NewMapper(provider)

	d, err := m.ConvertGraph(g)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(g.Name, "/")

	switch provider {
	case "aws":
		def := d.(*aws.Definition)
		def.Name = parts[1]
		def.Project = parts[0]
	case "azure":
		def := d.(*azure.Definition)
		def.Name = parts[1]
		def.Project = parts[0]
	}

	data, err := yaml.Marshal(d)
	if err != nil {
		return nil, err
	}

	b := build{
		ID:         g.ID,
		Definition: string(data),
		Mapping:    g,
	}

	return &b, err
}
