/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// RDSBackup ...
type RDSBackup struct {
	Window    string `json:"window" yaml:"window"`
	Retention *int64 `json:"retention" yaml:"retention"`
}

// RDSCluster ...
type RDSCluster struct {
	Name              string    `json:"name" yaml:"name"`
	Engine            string    `json:"engine" yaml:"engine"`
	EngineVersion     string    `json:"engine_version" yaml:"engine_version"`
	Port              *int64    `json:"port" yaml:"port"`
	AvailabilityZones []string  `json:"availability_zones" yaml:"availability_zones"`
	SecurityGroups    []string  `json:"security_groups" yaml:"security_groups"`
	Networks          []string  `json:"networks" yaml:"networks"`
	DatabaseName      string    `json:"database_name" yaml:"database_name"`
	DatabaseUsername  string    `json:"database_username" yaml:"database_username"`
	DatabasePassword  string    `json:"database_password" yaml:"database_password"`
	Backups           RDSBackup `json:"backups" yaml:"backups"`
	MaintenanceWindow string    `json:"maintenance_window" yaml:"maintenance_window"`
	ReplicationSource string    `json:"replication_source" yaml:"replication_source"`
	FinalSnapshot     bool      `json:"final_snapshot" yaml:"final_snapshot"`
}
