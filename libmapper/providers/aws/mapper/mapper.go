/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"errors"

	"github.com/ernestio/definition-mapper/libmapper"
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/components"
	def "github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
	"github.com/mitchellh/mapstructure"
	graph "gopkg.in/r3labs/graph.v2"
)

// SUPPORTEDCOMPONENTS represents all component types supported by ernest
var SUPPORTEDCOMPONENTS = []string{"vpc", "internet_gateway", "network", "instance", "firewall", "nat_gateway", "elb", "ebs", "s3", "route53", "rds_instance", "rds_cluster"}

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
		c.Rebuild(g)

		// Validate Components
		err := c.Validate()
		if err != nil {
			return g, err
		}

		// Build internal & template values
		for _, dep := range c.Dependencies() {
			if g.HasComponent(dep) != true {
				return g, errors.New("Could not resolve component dependency: " + dep)
			}
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
	var d def.Definition

	for _, c := range g.Components {
		c.Rebuild(g)

		for _, dep := range c.Dependencies() {
			if g.HasComponent(dep) != true {
				return g, errors.New("Could not resolve component dependency: " + dep)
			}
		}

		err := c.Validate()
		if err != nil {
			return d, err
		}
	}

	d.Vpcs = MapDefinitionVpcs(g)
	d.Networks = MapDefinitionNetworks(g)
	d.Instances = MapDefinitionInstances(g)
	d.SecurityGroups = MapDefinitionSecurityGroups(g)
	d.ELBs = MapDefinitionELBs(g)
	d.EBSVolumes = MapDefinitionEBSVolumes(g)
	d.NatGateways = MapDefinitionNats(g)
	d.RDSClusters = MapDefinitionRDSClusters(g)
	d.RDSInstances = MapDefinitionRDSInstances(g)
	d.Route53Zones = MapDefinitionRoute53Zones(g)
	d.S3Buckets = MapDefinitionS3Buckets(g)

	return d, nil
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
		case "vpc":
			c = &components.Vpc{}
		case "network":
			c = &components.Network{}
		case "internet_gateway":
			c = &components.InternetGateway{}
		case "instance":
			c = &components.Instance{}
		case "firewall":
			c = &components.SecurityGroup{}
		case "elb":
			c = &components.ELB{}
		case "ebs_volume":
			c = &components.EBSVolume{}
		case "nat":
			c = &components.NatGateway{}
		case "rds_cluster":
			c = &components.RDSCluster{}
		case "rds_instance":
			c = &components.RDSInstance{}
		case "route53":
			c = &components.Route53Zone{}
		case "s3":
			c = &components.S3Bucket{}
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
	filter := make(map[string]string)

	if len(params) > 0 {
		filter["ernest.service"] = params[0]
	}

	for _, ctype := range SUPPORTEDCOMPONENTS {
		q := MapQuery(ctype+"s", filter)
		g.AddComponent(q)
	}

	return g
}

// ProviderCredentials : maps aws credentials to a generic component
func (m Mapper) ProviderCredentials(details map[string]interface{}) graph.Component {
	credentials := make(graph.GenericComponent)

	credentials["_action"] = "none"
	credentials["_component"] = "credentials"
	credentials["_component_id"] = "credentials::aws"
	credentials["_provider"] = details["type"]
	credentials["name"] = details["name"]
	credentials["region"] = details["region"]
	credentials["aws_access_key_id"] = details["aws_access_key_id"]
	credentials["aws_secret_access_key"] = details["aws_secret_access_key"]

	return &credentials
}

func mapComponents(d *def.Definition, g *graph.Graph) error {
	// Map basic component values from definition

	for _, vpc := range MapVpcs(d) {
		err := g.AddComponent(vpc)
		if err != nil {
			return err
		}
	}

	for _, gateway := range MapInternetGateways(d) {
		err := g.AddComponent(gateway)
		if err != nil {
			return err
		}
	}

	for _, network := range MapNetworks(d) {
		err := g.AddComponent(network)
		if err != nil {
			return err
		}
	}

	for _, instance := range MapInstances(d) {
		err := g.AddComponent(instance)
		if err != nil {
			return err
		}
	}

	for _, securitygroup := range MapSecurityGroups(d) {
		err := g.AddComponent(securitygroup)
		if err != nil {
			return err
		}
	}

	for _, elb := range MapELBs(d) {
		err := g.AddComponent(elb)
		if err != nil {
			return err
		}
	}

	for _, ebs := range MapEBSVolumes(d) {
		err := g.AddComponent(ebs)
		if err != nil {
			return err
		}
	}

	for _, nat := range MapNats(d) {
		err := g.AddComponent(nat)
		if err != nil {
			return err
		}
	}

	for _, rds := range MapRDSClusters(d) {
		err := g.AddComponent(rds)
		if err != nil {
			return err
		}
	}

	for _, rds := range MapRDSInstances(d) {
		err := g.AddComponent(rds)
		if err != nil {
			return err
		}
	}

	for _, s3 := range MapS3Buckets(d) {
		err := g.AddComponent(s3)
		if err != nil {
			return err
		}
	}

	for _, route53 := range MapRoute53Zones(d) {
		err := g.AddComponent(route53)
		if err != nil {
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
