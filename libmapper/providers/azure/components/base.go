/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

// Base : Shared internal component fields
type Base struct {
	ProviderType     string `json:"_provider" diff:"-"`
	ComponentID      string `json:"_component_id" diff:"_component_id,immutable"`
	ComponentType    string `json:"_component" diff:"-"`
	State            string `json:"_state" diff:"-"`
	Action           string `json:"_action" diff:"-"`
	DatacenterName   string `json:"datacenter_name" diff:"-"`
	DatacenterType   string `json:"datacenter_type" diff:"-"`
	DatacenterRegion string `json:"datacenter_region" diff:"-"`
}
