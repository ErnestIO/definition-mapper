/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

import (
	"errors"
	"fmt"
	"strings"

	"github.com/r3labs/diff"
	"github.com/r3labs/graph"
)

var (
	// CNAME ...
	CNAME = "CNAME"

	// DNSTYPES ...
	DNSTYPES = []string{"A", "AAAA", "CNAME", "MX", "PTR", "TXT", "SRV", "SPF", "NAPTR", "NS", "SOA"}
)

// Record stores the entries for a zone
type Record struct {
	Entry         string   `json:"entry" diff:"entry,identifier"`
	Type          string   `json:"type" diff:"type"`
	Instances     []string `json:"instances,omitempty" diff:"instances"`
	Loadbalancers []string `json:"loadbalancers,omitempty" diff:"loadbalancers"`
	RDSClusters   []string `json:"rds_clusters,omitempty" diff:"rds_clusters"`
	RDSInstances  []string `json:"rds_instances,omitempty" diff:"rds_instances"`
	Values        []string `json:"values" diff:"values"`
	TTL           int64    `json:"ttl" diff:"ttl"`
}

// Route53Zone holds all information about a dns zone
type Route53Zone struct {
	ProviderType     string            `json:"_provider" diff:"-"`
	ComponentType    string            `json:"_component" diff:"-"`
	ComponentID      string            `json:"_component_id" diff:"component_id,identifier"`
	State            string            `json:"_state" diff:"-"`
	Action           string            `json:"_action" diff:"-"`
	HostedZoneID     string            `json:"hosted_zone_id" diff:"-"`
	Name             string            `json:"name" diff:"-"`
	Private          bool              `json:"private" diff:"-"`
	Records          []Record          `json:"records" diff:"records"`
	Vpc              string            `json:"vpc" diff:"-"`
	VpcID            string            `json:"vpc_id" diff:"-"`
	Tags             map[string]string `json:"tags" diff:"-"`
	DatacenterType   string            `json:"datacenter_type" diff:"-"`
	DatacenterName   string            `json:"datacenter_name" diff:"-"`
	DatacenterRegion string            `json:"datacenter_region" diff:"-"`
	AccessKeyID      string            `json:"aws_access_key_id" diff:"-"`
	SecretAccessKey  string            `json:"aws_secret_access_key" diff:"-"`
	Service          string            `json:"service" diff:"-"`
}

// GetID : returns the component's ID
func (z *Route53Zone) GetID() string {
	return z.ComponentID
}

// GetName returns a components name
func (z *Route53Zone) GetName() string {
	return z.Name
}

// GetProvider : returns the provider type
func (z *Route53Zone) GetProvider() string {
	return z.ProviderType
}

// GetProviderID returns a components provider id
func (z *Route53Zone) GetProviderID() string {
	return z.HostedZoneID
}

// GetType : returns the type of the component
func (z *Route53Zone) GetType() string {
	return z.ComponentType
}

// GetState : returns the state of the component
func (z *Route53Zone) GetState() string {
	return z.State
}

// SetState : sets the state of the component
func (z *Route53Zone) SetState(s string) {
	z.State = s
}

// GetAction : returns the action of the component
func (z *Route53Zone) GetAction() string {
	return z.Action
}

// SetAction : Sets the action of the component
func (z *Route53Zone) SetAction(s string) {
	z.Action = s
}

// GetGroup : returns the components group
func (z *Route53Zone) GetGroup() string {
	return ""
}

// GetTags returns a components tags
func (z *Route53Zone) GetTags() map[string]string {
	return z.Tags
}

// GetTag returns a components tag
func (z *Route53Zone) GetTag(tag string) string {
	return z.Tags[tag]
}

// Diff : diff's the component against another component of the same type
func (z *Route53Zone) Diff(c graph.Component) (diff.Changelog, error) {
	cz, ok := c.(*Route53Zone)
	if ok {
		return diff.Diff(cz, z)
	}

	return diff.Changelog{}, nil
}

// Update : updates the provider returned values of a component
func (z *Route53Zone) Update(c graph.Component) {
	cz, ok := c.(*Route53Zone)
	if ok {
		z.HostedZoneID = cz.HostedZoneID
	}

	z.SetDefaultVariables()
}

// Rebuild : rebuilds the component's internal state, such as templated values
func (z *Route53Zone) Rebuild(g *graph.Graph) {
	for i := 0; i < len(z.Records); i++ {
		if len(z.Records[i].Values) < 1 {
			// rebuild instance values
			for _, name := range z.Records[i].Instances {
				for _, gi := range g.GetComponents().ByType(TYPEINSTANCE).ByName(name) {
					instance := gi.(*Instance)

					if z.Private {
						z.Records[i].Values = append(z.Records[i].Values, templInstancePrivateIP(instance.Name))
					} else {
						z.Records[i].Values = append(z.Records[i].Values, templInstancePublicIP(instance.Name))
					}
				}
			}

			// rebuild loadbalancer values
			for _, name := range z.Records[i].Loadbalancers {
				z.Records[i].Values = append(z.Records[i].Values, templELBDNS(name))
			}

			// rebuild rds cluster values
			for _, name := range z.Records[i].RDSClusters {
				z.Records[i].Values = append(z.Records[i].Values, templRDSClusterDNS(name))
			}

			// rebuild rds instance values
			for _, name := range z.Records[i].RDSInstances {
				z.Records[i].Values = append(z.Records[i].Values, templRDSInstanceDNS(name))
			}
		}

		if len(z.Records[i].Instances) < 1 &&
			len(z.Records[i].Loadbalancers) < 1 &&
			len(z.Records[i].RDSClusters) < 1 &&
			len(z.Records[i].RDSInstances) < 1 {
			for x, v := range z.Records[i].Values {

				// rebuild instance names
				for _, gi := range g.GetComponents().ByType(TYPEINSTANCE) {
					in := gi.(*Instance)
					if in.IP == v {
						z.Records[i].Instances = append(z.Records[i].Instances, in.Name)
						z.Records[i].Values[x] = templInstancePrivateIP(in.Name)
					}
					if in.PublicIP == v {
						z.Records[i].Instances = append(z.Records[i].Instances, in.Name)
						z.Records[i].Values[x] = templInstancePublicIP(in.Name)
					}
					if in.ElasticIP == v {
						z.Records[i].Instances = append(z.Records[i].Instances, in.Name)
						z.Records[i].Values[x] = templInstanceElasticIP(in.Name)
					}
				}

				// rebuild loadbalancer names
				for _, ge := range g.GetComponents().ByType(TYPEELB) {
					elb := ge.(*ELB)
					if elb.DNSName == v {
						z.Records[i].Loadbalancers = append(z.Records[i].Loadbalancers, elb.Name)
						z.Records[i].Values[x] = templELBDNS(elb.Name)
					}
				}

				// rebuild rds cluster names
				for _, gr := range g.GetComponents().ByType(TYPERDSCLUSTER) {
					rds := gr.(*RDSCluster)
					if rds.Endpoint == v {
						z.Records[i].RDSClusters = append(z.Records[i].RDSClusters, rds.Name)
						z.Records[i].Values[x] = templRDSClusterDNS(rds.Name)
					}
				}

				// rebuild rds instance names
				for _, gr := range g.GetComponents().ByType(TYPERDSINSTANCE) {
					rds := gr.(*RDSInstance)
					if rds.Endpoint == v {
						z.Records[i].RDSInstances = append(z.Records[i].RDSInstances, rds.Name)
						z.Records[i].Values[x] = templRDSInstanceDNS(rds.Name)
					}
				}
			}
		}
	}

	if z.Private {
		if z.VpcID != "" && z.Vpc == "" {
			v := g.GetComponents().ByProviderID(z.VpcID)
			if v != nil {
				z.Vpc = v.GetName()
			}
		}
		if z.Vpc != "" && z.VpcID == "" {
			templVpcID(z.Vpc)
		}
	}

	z.SetDefaultVariables()
}

// Dependencies : returns a list of component id's upon which the component depends
func (z *Route53Zone) Dependencies() []string {
	var deps []string

	for _, record := range z.Records {
		for _, i := range record.Instances {
			deps = append(deps, TYPEINSTANCE+TYPEDELIMITER+i)
		}

		for _, l := range record.Loadbalancers {
			deps = append(deps, TYPEELB+TYPEDELIMITER+l)
		}

		for _, r := range record.RDSClusters {
			deps = append(deps, TYPERDSCLUSTER+TYPEDELIMITER+r)
		}

		for _, r := range record.RDSInstances {
			deps = append(deps, TYPERDSINSTANCE+TYPEDELIMITER+r)
		}
	}

	if z.Vpc != "" {
		deps = append(deps, TYPEVPC+TYPEDELIMITER+z.Vpc)
	}

	return deps
}

// SequentialDependencies : returns a list of origin components that restrict the execution of its dependents, allowing only one dependent component to be provisioned at a time (sequentially)
func (z *Route53Zone) SequentialDependencies() []string {
	return []string{}
}

// Validate : validates the components values
func (z *Route53Zone) Validate() error {
	if z.Name == "" {
		return errors.New("Route53 zone name should not be null")
	}

	if z.Private && z.Vpc == "" {
		return errors.New("Route53 private zone must specify a vpc!")
	}

	for _, record := range z.Records {
		if record.Entry == "" {
			return errors.New("Route53 record entry name should not be null")
		}

		if !isOneOf(DNSTYPES, record.Type) {
			return fmt.Errorf("Route53 record type '%s' is not a valid dns type. Please use one of [%s]", record.Type, strings.Join(DNSTYPES, ", "))
		}

		if record.TTL == 0 {
			return errors.New("Route53 record TTL must be greater than 0")
		}
	}

	return nil
}

// IsStateful : returns true if the component needs to be actioned to be removed.
func (z *Route53Zone) IsStateful() bool {
	return true
}

// SetDefaultVariables : sets up the default template variables for a component
func (z *Route53Zone) SetDefaultVariables() {
	z.ComponentType = TYPEROUTE53
	z.ComponentID = TYPEROUTE53 + TYPEDELIMITER + z.Name
	z.ProviderType = PROVIDERTYPE
	z.DatacenterName = DATACENTERNAME
	z.DatacenterType = DATACENTERTYPE
	z.DatacenterRegion = DATACENTERREGION
	z.AccessKeyID = ACCESSKEYID
	z.SecretAccessKey = SECRETACCESSKEY
}
