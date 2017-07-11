/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// RDSStorage ...
type RDSStorage struct {
	Type string `json:"type" yaml:"type"`
	Size *int64 `json:"size" yaml:"size"`
	Iops *int64 `json:"iops" yaml:"iops"`
}

// RDSInstance ...
type RDSInstance struct {
	Name              string     `json:"name" yaml:"name"`
	Size              string     `json:"size" yaml:"size"`
	Engine            string     `json:"engine" yaml:"engine"`
	EngineVersion     string     `json:"engine_version" yaml:"engine_version"`
	Port              *int64     `json:"port" yaml:"port"`
	Cluster           string     `json:"cluster" yaml:"cluster"`
	Public            bool       `json:"public" yaml:"public"`
	MultiAZ           bool       `json:"multi_az" yaml:"multi_az"`
	PromotionTier     *int64     `json:"promotion_tier" yaml:"promotion_tier"`
	Storage           RDSStorage `json:"storage" yaml:"storage"`
	AvailabilityZone  string     `json:"availability_zone" yaml:"availability_zone"`
	SecurityGroups    []string   `json:"security_groups" yaml:"security_groups"`
	Networks          []string   `json:"networks" yaml:"networks"`
	DatabaseName      string     `json:"database_name" yaml:"database_name"`
	DatabaseUsername  string     `json:"database_username" yaml:"database_username"`
	DatabasePassword  string     `json:"database_password" yaml:"database_password"`
	AutoUpgrade       bool       `json:"auto_upgrade" yaml:"auto_upgrade"`
	Backups           RDSBackup  `json:"backups" yaml:"backups"`
	MaintenanceWindow string     `json:"maintenance_window" yaml:"maintenance_window"`
	FinalSnapshot     bool       `json:"final_snapshot" yaml:"final_snapshot"`
	ReplicationSource string     `json:"replication_source" yaml:"replication_source"`
	License           string     `json:"license" yaml:"license"`
	Timezone          string     `json:"timezone" yaml:"timezone"`
}
