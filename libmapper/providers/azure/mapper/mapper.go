/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"errors"
	"log"

	"github.com/ernestio/definition-mapper/libmapper"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	def "github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	"github.com/mitchellh/mapstructure"
	graph "gopkg.in/r3labs/graph.v2"
)

// SUPPORTEDCOMPONENTS represents all component types supported by ernest
var SUPPORTEDCOMPONENTS = []string{"network_interface", "public_ip", "resource_group", "security_group", "sql_database", "sql_server", "storage_account", "storage_container", "subnet", "virtual_machine", "virtual_network"}

// Mapper : implements the generic mapper structure
type Mapper struct{}

// New : returns a new azure mapper
func New() libmapper.Mapper {
	return &Mapper{}
}

// ConvertDefinition : converts the input yaml definition to a graph format
func (m Mapper) ConvertDefinition(gd libmapper.Definition) (*graph.Graph, error) {
	g := graph.New()

	d, ok := gd.(*def.Definition)
	if ok != true {
		return g, errors.New("Could not convert generic definition into azure format")
	}

	// Map basic component values from definition
	err := mapComponents(d, g)
	if err != nil {
		return g, err
	}

	for _, c := range g.Components {
		c.Rebuild(g)

		// Validate Components
		err := c.Validate()
		if err != nil {
			return g, err
		}

		// Build internal & template values
		for _, dep := range c.Dependencies() {
			if g.HasComponent(dep) != true {
				return g, errors.New("Component '" + c.GetID() + "': Could not resolve component dependency '" + dep + "'")
			}
		}

		// Build dependencies
		for _, dep := range c.Dependencies() {
			_ = g.Connect(dep, c.GetID())
		}
	}

	return g, nil
}

// ConvertGraph : converts the service graph into an input yaml format
func (m Mapper) ConvertGraph(g *graph.Graph) (libmapper.Definition, error) {
	var d def.Definition

	for i := len(g.Components) - 1; i >= 0; i-- {
		c := g.Components[i]
		c.Rebuild(g)

		// remove any components that were determined to not be apart of the service
		if c.IsStateful() != true {
			g.Components = append(g.Components[:i], g.Components[i+1:]...)
			continue
		}

		for _, dep := range c.Dependencies() {
			if g.HasComponent(dep) != true {
				return g, errors.New("Component '" + c.GetID() + "': Could not resolve component dependency '" + dep + "'")
			}
		}

		err := c.Validate()
		if err != nil {
			return d, err
		}
	}

	d.ResourceGroups = MapDefinitionResourceGroups(g)

	for i := 0; i < len(d.ResourceGroups); i++ {
		d.ResourceGroups[i].NetworkInterfaces = MapDefinitionNetworkInterfaces(g, &d.ResourceGroups[i])
		d.ResourceGroups[i].VirtualNetworks = MapDefinitionVirtualNetworks(g, &d.ResourceGroups[i])
		for x := 0; x < len(d.ResourceGroups[i].VirtualNetworks); x++ {
			d.ResourceGroups[i].VirtualNetworks[x].Subnets = MapDefinitionSubnets(g, &d.ResourceGroups[i], &d.ResourceGroups[i].VirtualNetworks[x])
		}
		d.ResourceGroups[i].VirtualMachines = MapDefinitionVirtualMachines(g, &d.ResourceGroups[i])
		d.ResourceGroups[i].PublicIPs = MapDefinitionPublicIPs(g, &d.ResourceGroups[i])
		d.ResourceGroups[i].SecurityGroups = MapDefinitionSecurityGroups(g, &d.ResourceGroups[i])
		d.ResourceGroups[i].SQLServers = MapDefinitionSQLServers(g, &d.ResourceGroups[i])

		for x := 0; x < len(d.ResourceGroups[i].SQLServers); x++ {
			d.ResourceGroups[i].SQLServers[x].Databases = MapDefinitionSQLDatabases(g, &d.ResourceGroups[i], &d.ResourceGroups[i].SQLServers[x])
		}
		d.ResourceGroups[i].StorageAccounts = MapDefinitionStorageAccounts(g, &d.ResourceGroups[i])
		for x := 0; x < len(d.ResourceGroups[i].StorageAccounts); x++ {
			d.ResourceGroups[i].StorageAccounts[x].Containers = MapDefinitionStorageContainers(g, &d.ResourceGroups[i], &d.ResourceGroups[i].StorageAccounts[x])
		}
	}

	return &d, nil
}

// LoadDefinition : returns an azure type definition
func (m Mapper) LoadDefinition(gd map[string]interface{}) (libmapper.Definition, error) {
	var d def.Definition

	err := d.LoadMap(gd)

	return &d, err
}

// LoadGraph : returns a generic interal graph
func (m Mapper) LoadGraph(gg map[string]interface{}) (*graph.Graph, error) {
	g := graph.New()

	_ = g.Load(gg)

	for i := 0; i < len(g.Components); i++ {
		gc := g.Components[i].(*graph.GenericComponent)

		var c graph.Component

		switch gc.GetType() {
		case "resource_group":
			c = &components.ResourceGroup{}
		case "network_interface":
			c = &components.NetworkInterface{}
		case "public_ip":
			c = &components.PublicIP{}
		case "security_group":
			c = &components.SecurityGroup{}
		case "virtual_network":
			c = &components.VirtualNetwork{}
		case "virtual_machine":
			c = &components.VirtualMachine{}
		case "subnet":
			c = &components.Subnet{}
		case "sql_server":
			c = &components.SQLServer{}
		case "sql_database":
			c = &components.SQLDatabase{}
		case "storage_account":
			c = &components.StorageAccount{}
		case "storage_container":
			c = &components.StorageContainer{}
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
		(*gc)["Base"] = gc
		(*gc)["Event"] = gc

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

	filter := make(map[string]string)

	if len(params) > 0 {
		filter["ernest.service"] = params[0]
	}

	for _, ctype := range SUPPORTEDCOMPONENTS {
		q := MapQuery(ctype+"s", filter)
		if err := g.AddComponent(q); err != nil {
			log.Println(err.Error())
		}
	}

	return g
}

// ProviderCredentials : maps aws credentials to a generic component
func (m Mapper) ProviderCredentials(details map[string]interface{}) graph.Component {
	credentials := make(graph.GenericComponent)

	credentials["_action"] = "none"
	credentials["_component"] = "credentials"
	credentials["_component_id"] = "credentials::azure"
	credentials["_provider"] = details["type"]
	credentials["name"] = details["name"]
	credentials["region"] = details["region"]
	credentials["azure_client_id"] = details["azure_client_id"]
	credentials["azure_client_secret"] = details["azure_client_secret"]
	credentials["azure_subscription_id"] = details["azure_subscription_id"]
	credentials["azure_tenant_id"] = details["azure_tenant_id"]
	credentials["azure_environment"] = details["azure_environment"]

	return &credentials
}

// mapComponents : Map basic component values from definition
func mapComponents(d *def.Definition, g *graph.Graph) error {
	for _, rg := range MapResourceGroups(d) {
		if err := g.AddComponent(rg); err != nil {
			return err
		}
	}

	for _, vn := range MapVirtualNetworks(d) {
		if err := g.AddComponent(vn); err != nil {
			return err
		}
	}

	for _, vm := range MapVirtualMachines(d) {
		if err := g.AddComponent(vm); err != nil {
			return err
		}
	}

	for _, subnet := range MapSubnets(d) {
		if err := g.AddComponent(subnet); err != nil {
			return err
		}
	}

	for _, ni := range MapNetworkInterfaces(d) {
		if err := g.AddComponent(ni); err != nil {
			return err
		}
	}

	for _, ip := range MapPublicIPs(d) {
		if err := g.AddComponent(ip); err != nil {
			return err
		}
	}

	for _, sg := range MapSecurityGroups(d) {
		if err := g.AddComponent(sg); err != nil {
			return err
		}
	}

	for _, ss := range MapSQLServers(d) {
		if err := g.AddComponent(ss); err != nil {
			return err
		}
	}

	for _, sd := range MapSQLDatabases(d) {
		if err := g.AddComponent(sd); err != nil {
			return err
		}
	}

	for _, sa := range MapStorageAccounts(d) {
		if err := g.AddComponent(sa); err != nil {
			return err
		}
	}

	for _, sc := range MapStorageContainers(d) {
		if err := g.AddComponent(sc); err != nil {
			return err
		}
	}

	return nil
}

func mapTags(name, service string) map[string]string {
	tags := make(map[string]string)

	tags["Name"] = name
	tags["ernest.service"] = service

	return tags
}

func mapTagsServiceOnly(service string) map[string]string {
	tags := make(map[string]string)

	tags["ernest.service"] = service

	return tags
}
