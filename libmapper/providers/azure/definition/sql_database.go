/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// Database ..
type Database struct {
	Name                          string            `json:"name" yaml:"name"`
	CreateMode                    string            `json:"create_mode" yaml:"create_mode"`
	SourceDatabaseID              string            `json:"source_database_id" yaml:"source_database_id"`
	RestorePointInTime            string            `json:"restore_point_in_time" yaml:"restore_point_in_time"`
	Edition                       string            `json:"edition" yaml:"edition"`
	Collation                     string            `json:"collation" yaml:"collation"`
	MaxSizeBytes                  string            `json:"max_size_bytes" yaml:"max_size_bytes"`
	RequestedServiceObjectiveID   string            `json:"requested_service_objective_id" yaml:"requested_service_objective_id"`
	RequestedServiceObjectiveName string            `json:"requested_service_objective_name" yaml:"requested_service_objective_name"`
	SourceDatabaseDeletionData    string            `json:"source_database_deletion_date" yaml:"source_database_deletion_date"`
	ElasticPoolName               string            `json:"elastic_pool_name" yaml:"elastic_pool_name"`
	Tags                          map[string]string `json:"tags" yaml:"tags"`
}
