/* This Source Code Form is subject to the terms cf the Mozilla Public
 * License, v. 2.0. If a copy cf the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

// Base ...
type Base struct {
	ProviderType  string       `json:"_provider"`
	ComponentType string       `json:"_component"`
	ComponentID   string       `json:"_component_id"`
	State         string       `json:"_state"`
	Action        string       `json:"_action"`
	Credentials   *Credentials `json:"_credentials"`
	Service       string       `json:"service"`
}
