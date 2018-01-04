/* This Source Code Form is subject to the terms cf the Mozilla Public
 * License, v. 2.0. If a copy cf the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

// Base ...
type Base struct {
	ProviderType  string       `json:"_provider" diff:"-"`
	ComponentType string       `json:"_component" diff:"-"`
	ComponentID   string       `json:"_component_id" diff:"component_id,identifier"`
	State         string       `json:"_state" diff:"-"`
	Action        string       `json:"_action" diff:"-"`
	Credentials   *Credentials `json:"_credentials" diff:"-"`
	Service       string       `json:"service" diff:"-"`
}
