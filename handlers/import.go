/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package handlers

import (
	"strings"

	"github.com/ernestio/definition-mapper/build"
	"github.com/ernestio/definition-mapper/libmapper/providers"
	aws "github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
	azure "github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	vcloud "github.com/ernestio/definition-mapper/libmapper/providers/vcloud/definition"
	"github.com/ernestio/definition-mapper/request"
	"github.com/r3labs/graph"
	yaml "gopkg.in/yaml.v2"
)

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
	g.UserID = r.UserID
	g.Username = r.Username

	return g, nil
}

// ImportComplete : handles the conversion of an import graph to a definition
func ImportComplete(ig map[string]interface{}) (*build.Build, error) {
	provider := getGraphProvider(ig)

	m := providers.NewMapper(provider)

	g, err := m.LoadGraph(ig)
	if err != nil {
		return nil, err
	}

	d, err := m.ConvertGraph(g)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(g.Name, "/")

	switch provider {
	case "aws", "aws-fake":
		def := d.(*aws.Definition)
		def.Name = parts[1]
		def.Project = parts[0]
	case "azure", "azure-fake":
		def := d.(*azure.Definition)
		def.Name = parts[1]
		def.Project = parts[0]
	case "vcloud", "vcloud-fake":
		def := d.(*vcloud.Definition)
		def.Name = parts[1]
		def.Project = parts[0]
	}

	data, err := yaml.Marshal(d)
	if err != nil {
		return nil, err
	}

	b := build.Build{
		ID:         g.ID,
		Definition: string(data),
		Mapping:    g,
	}

	return &b, err
}

func getGraphProvider(m map[string]interface{}) string {
	for _, c := range m["components"].([]interface{}) {
		x := c.(map[string]interface{})
		if x["_component"].(string) == "credentials" {
			return x["_provider"].(string)
		}
	}
	return ""
}
