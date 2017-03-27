/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// Record stores the entries for a zone
type Record struct {
	Entry         string   `json:"entry" yaml:"entry"`
	Type          string   `json:"type" yaml:"type"`
	Instances     []string `json:"instances,omitempty" yaml:"instances,omitempty"`
	Loadbalancers []string `json:"loadbalancers,omitempty" yaml:"loadbalancers,omitempty"`
	RDSClusters   []string `json:"rds_clusters,omitempty" yaml:"rds_clusters,omitempty"`
	RDSInstances  []string `json:"rds_instances,omitempty" yaml:"rds_instances,omitempty"`
	Values        []string `json:"values,omitempty" yaml:"values,omitempty"`
	TTL           int64    `json:"ttl" yaml:"ttl"`
}

// Route53Zone ...
type Route53Zone struct {
	Name    string   `json:"name" yaml:"name"`
	Private bool     `json:"private" yaml:"private"`
	Records []Record `json:"records" yaml:"records"`
	Vpc     string   `json:"vpc" yaml:"vpc"`
}
