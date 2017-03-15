/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"errors"

	"github.com/ernestio/definition-mapper/libmapper"
	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/components"
	def "github.com/ernestio/definition-mapper/libmapper/providers/vcloud/definition"
	"github.com/mitchellh/mapstructure"
	graph "gopkg.in/r3labs/graph.v2"
)

// SUPPORTEDCOMPONENTS represents all component types supported by ernest
var SUPPORTEDCOMPONENTS = []string{"router", "network", "instance"}

// Mapper : implements the generic mapper structure
type Mapper struct{}

// New : returns a new aws mapper
func New() libmapper.Mapper {
	return &Mapper{}
}

// ConvertDefinition : converts the input yaml definition to a graph format
func (m Mapper) ConvertDefinition(gd libmapper.Definition) (*graph.Graph, error) {
	g := graph.New()

	d, ok := gd.(*def.Definition)
	if ok != true {
		return g, errors.New("Could not convert generic definition into aws format")
	}

	// Map basic component values from definition
	err := mapComponents(d, g)
	if err != nil {
		return g, err
	}

	for _, c := range g.Components {
		// Build internal & template values
		for _, dep := range c.Dependencies() {
			if g.HasComponent(dep) != true {
				return g, errors.New("Could not resolve component dependency: " + dep)
			}
		}

		c.Rebuild(g)

		// Validate Components
		err := c.Validate()
		if err != nil {
			return g, err
		}

		// Build dependencies
		for _, dep := range c.Dependencies() {
			g.Connect(dep, c.GetID())
		}
	}

	return g, nil
}

// ConvertGraph : converts the service graph into an input yaml format
func (m Mapper) ConvertGraph(g *graph.Graph) (libmapper.Definition, error) {
	return def.Definition{}, nil
}

// LoadDefinition : returns an aws type definition
func (m Mapper) LoadDefinition(gd map[string]interface{}) (libmapper.Definition, error) {
	var d def.Definition

	err := d.LoadMap(gd)

	return &d, err
}

// LoadGraph : returns a generic interal graph
func (m Mapper) LoadGraph(gg map[string]interface{}) (*graph.Graph, error) {
	g := graph.New()

	g.Load(gg)

	for i := 0; i < len(g.Components); i++ {
		gc := g.Components[i].(*graph.GenericComponent)

		var c graph.Component

		switch gc.GetType() {
		case "router":
			c = &components.Router{}
		case "network":
			c = &components.Network{}
		case "instance":
			c = &components.Instance{}
		default:
			continue
		}

		config := &mapstructure.DecoderConfig{
			Metadata: nil,
			Result:   c,
			TagName:  "json",
		}

		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			return g, err
		}

		err = decoder.Decode(gc)
		if err != nil {
			return g, err
		}

		g.Components[i] = c
	}

	return g, nil
}

// CreateImportGraph : creates a new graph with component queries used to import components from a provider
func (m Mapper) CreateImportGraph(params []string) *graph.Graph {
	g := graph.New()

	return g
}

// ProviderCredentials : maps aws credentials to a generic component
func (m Mapper) ProviderCredentials(details map[string]interface{}) graph.Component {
	credentials := make(graph.GenericComponent)

	credentials["_action"] = "none"
	credentials["_component"] = "credentials"
	credentials["_component_id"] = "credentials::vcloud"
	credentials["_provider"] = details["type"]
	credentials["name"] = details["name"]
	credentials["region"] = details["region"]
	credentials["username"] = details["username"]
	credentials["password"] = details["password"]
	credentials["vcloud_url"] = details["vcloud_url"]
	credentials["external_network"] = details["external_network"]

	return &credentials
}

func mapComponents(d *def.Definition, g *graph.Graph) error {
	// Map basic component values from definition

	for _, router := range MapRouters(d) {
		err := g.AddComponent(router)
		if err != nil {
			return err
		}
	}

	/*
		for _, network := range MapNetworks(d) {
			err := g.AddComponent(network)
			if err != nil {
				return err
			}
		}
	*/

	for _, instance := range MapInstances(d) {
		err := g.AddComponent(instance)
		if err != nil {
			return err
		}
	}

	return nil
}
