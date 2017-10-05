/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package request

import (
	"strings"

	"github.com/ernestio/definition-mapper/libmapper"
	"github.com/r3labs/graph"
)

// Request :
type Request struct {
	ID          string                 `json:"id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Filters     []string               `json:"filters,omitempty"`
	Definition  map[string]interface{} `json:"definition,omitempty"`
	From        map[string]interface{} `json:"from,omitempty"`
	To          map[string]interface{} `json:"to,omitempty"`
	Credentials map[string]interface{} `json:"credentials,omitempty"`
}

// DefinitionToGraph : converts a Defintiion to a graph
func (r *Request) DefinitionToGraph(m libmapper.Mapper) (*graph.Graph, error) {
	d, err := m.LoadDefinition(r.Definition)
	if err != nil {
		return nil, err
	}

	g, err := m.ConvertDefinition(d)
	if err != nil {
		return nil, err
	}

	// set graph ID and credentials
	g.ID = r.ID
	g.Name = r.Name

	err = g.AddComponent(m.ProviderCredentials(r.Credentials))
	if err != nil {
		return nil, err
	}

	return g, nil
}

// ToMapping : loads the "to" graph mapping as a graph
func (r *Request) ToMapping(m libmapper.Mapper) (*graph.Graph, error) {
	return m.LoadGraph(r.To)
}

// FromMapping : loads the "from" graph mapping as a graph
func (r *Request) FromMapping(m libmapper.Mapper) (*graph.Graph, error) {
	return m.LoadGraph(r.From)
}

// Provider : returns the provider/env type
func (r *Request) Provider() string {
	p, _ := r.Credentials["type"].(string)
	return p
}

// ImportFilters : returns the collection of import filters used on an import
func (r *Request) ImportFilters() []string {
	switch r.Provider() {
	case "azure", "azure-fake":
		return r.Filters
	default:
		return []string{env(r.Name)}
	}
}

func env(e string) string {
	return strings.Split(e, "/")[1]
}
