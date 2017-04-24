/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

const (
	TYPEDELIMITER        = "::"
	TYPENETWORKINTERFACE = "network_interface"
	TYPEPUBLICIP         = "public_ip"
	TYPERESOURCEGROUP    = "resource_group"
	TYPESECURITYGROUP    = "security_group"
	TYPESQLDATABASE      = "sql_database"
	TYPESQLSERVER        = "sql_server"
	TYPESTORAGEACCOUNT   = "storage_account"
	TYPESTORAGECONTAINER = "storage_container"
	TYPESUBNET           = "subnet"
	TYPEVIRTUALMACHINE   = "virtual_machine"
	TYPEVIRTUALNETWORK   = "virtual_network"

	GROUPINSTANCE  = "ernest.instance_group"
	GROUPEBSVOLUME = "ernest.volume_group"

	PROVIDERTYPE     = `$(components.#[_component_id="credentials::azure"]._provider)`
	DATACENTERNAME   = `$(components.#[_component_id="credentials::azure"].name)`
	DATACENTERTYPE   = `$(components.#[_component_id="credentials::azure"]._provider)`
	DATACENTERREGION = `$(components.#[_component_id="credentials::azure"].region)`
	CLIENTID         = `$(components.#[_component_id="credentials::azure"].azure_client_id)`
	CLIENTSECRET     = `$(components.#[_component_id="credentials::azure"].azure_client_secret)`
	TENANTID         = `$(components.#[_component_id="credentials::azure"].azure_tenant_id)`
	SUBSCRIPTIONID   = `$(components.#[_component_id="credentials::azure"].azure_subscription_id)`
	ENVIRONMENT      = `$(components.#[_component_id="credentials::azure"].azure_environment`
)
