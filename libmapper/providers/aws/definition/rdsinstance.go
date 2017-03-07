/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package definition

// RDSStorage ...
type RDSStorage struct {
	Type string `json:"type"`
	Size *int64 `json:"size"`
	Iops *int64 `json:"iops"`
}

// RDSInstance ...
type RDSInstance struct {
	Name              string     `json:"name"`
	Size              string     `json:"size"`
	Engine            string     `json:"engine"`
	EngineVersion     string     `json:"engine_version"`
	Port              *int64     `json:"port"`
	Cluster           string     `json:"cluster"`
	Public            bool       `json:"public"`
	MultiAZ           bool       `json:"multi_az"`
	PromotionTier     *int64     `json:"promotion_tier"`
	Storage           RDSStorage `json:"storage"`
	AvailabilityZone  string     `json:"availability_zone"`
	SecurityGroups    []string   `json:"security_groups"`
	Networks          []string   `json:"networks"`
	DatabaseName      string     `json:"database_name"`
	DatabaseUsername  string     `json:"database_username"`
	DatabasePassword  string     `json:"database_password"`
	AutoUpgrade       bool       `json:"auto_upgrade"`
	Backups           RDSBackup  `json:"backups"`
	MaintenanceWindow string     `json:"maintenance_window"`
	FinalSnapshot     bool       `json:"final_snapshot"`
	ReplicationSource string     `json:"replication_source"`
	License           string     `json:"license"`
	Timezone          string     `json:"timezone"`
}
