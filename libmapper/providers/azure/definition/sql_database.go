/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// SQLDatabase ..
type SQLDatabase struct {
	ID                            string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name                          string            `json:"name,omitempty" yaml:"name,omitempty"`
	CreateMode                    string            `json:"create_mode,omitempty" yaml:"create_mode,omitempty"`
	SourceDatabaseID              string            `json:"source_database_id,omitempty" yaml:"source_database_id,omitempty"`
	RestorePointInTime            string            `json:"restore_point_in_time,omitempty" yaml:"restore_point_in_time,omitempty"`
	Edition                       string            `json:"edition,omitempty" yaml:"edition,omitempty"`
	Collation                     string            `json:"collation,omitempty" yaml:"collation,omitempty"`
	MaxSizeBytes                  string            `json:"max_size_bytes,omitempty" yaml:"max_size_bytes,omitempty"`
	RequestedServiceObjectiveID   string            `json:"requested_service_objective_id,omitempty" yaml:"requested_service_objective_id,omitempty"`
	RequestedServiceObjectiveName string            `json:"requested_service_objective_name,omitempty" yaml:"requested_service_objective_name,omitempty"`
	SourceDatabaseDeletionData    string            `json:"source_database_deletion_date,omitempty" yaml:"source_database_deletion_date,omitempty"`
	Tags                          map[string]string `json:"tags,omitempty" yaml:"tags,omitempty"`
}
