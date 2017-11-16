/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

const (
	TYPEDELIMITER = "::"
	TYPEROUTER    = "router"
	TYPENETWORK   = "network"
	TYPEINSTANCE  = "instance"

	GROUPINSTANCE = "ernest.instance_group"

	PROVIDERTYPE       = `$(components.#[_component_id="credentials::vcloud"]._provider)`
	DATACENTERNAME     = `$(components.#[_component_id="credentials::vcloud"].vdc)`
	DATACENTERTYPE     = `$(components.#[_component_id="credentials::vcloud"]._provider)`
	DATACENTERUSERNAME = `$(components.#[_component_id="credentials::vcloud"].username)`
	DATACENTERPASSWORD = `$(components.#[_component_id="credentials::vcloud"].password)`
	DATACENTERREGION   = `$(components.#[_component_id="credentials::vcloud"].region)`
	VCLOUDURL          = `$(components.#[_component_id="credentials::vcloud"].vcloud_url)`
)
