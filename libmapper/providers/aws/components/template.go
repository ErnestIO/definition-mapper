/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package components

const (
	TYPEDELIMITER     = "::"
	TYPEVPC           = "vpc"
	TYPENETWORK       = "network"
	TYPEINSTANCE      = "instance"
	TYPEELB           = "elb"
	TYPEEBSVOLUME     = "ebs_volume"
	TYPESECURITYGROUP = "security_group"
	TYPENATGATEWAY    = "nat"
	TYPERDSCLUSTER    = "rds_cluster"
	TYPERDSINSTANCE   = "rds_instance"
	TYPEROUTE53       = "route_53"
	TYPES3BUCKET      = "s3"

	GROUPINSTANCE  = "ernest.instance_group"
	GROUPEBSVOLUME = "ernest.volume_group"

	PROVIDERTYPE     = `$(components.#[_component_id="credentials::aws"]._provider)`
	DATACENTERNAME   = `$(components.#[_component_id="credentials::aws"].name)`
	DATACENTERTYPE   = `$(components.#[_component_id="credentials::aws"]._provider)`
	ACCESSKEYID      = `$(components.#[_component_id="credentials::aws"].aws_access_key_id)`
	SECRETACCESSKEY  = `$(components.#[_component_id="credentials::aws"].aws_secret_access_key)`
	DATACENTERREGION = `$(components.#[_component_id="credentials::aws"].region)`
)

func templVpcID(vpc string) string {
	return `$(components.#[_component_id="` + TYPEVPC + TYPEDELIMITER + vpc + `"].vpc_aws_id)`
}

func templSecurityGroupID(sg string) string {
	return `$(components.#[_component_id="` + TYPESECURITYGROUP + TYPEDELIMITER + sg + `"].security_group_aws_id)`
}

func templSubnetID(nw string) string {
	return `$(components.#[_component_id="` + TYPENETWORK + TYPEDELIMITER + nw + `"].network_aws_id)`
}

func templInstanceID(in string) string {
	return `$(components.#[_component_id="` + TYPEINSTANCE + TYPEDELIMITER + in + `"].instance_aws_id)`
}

func templEBSVolumeID(ebs string) string {
	return `$(components.#[_component_id="` + TYPEEBSVOLUME + TYPEDELIMITER + ebs + `"].volume_aws_id)`
}

func templInstancePublicIP(in string) string {
	return `$(components.#[_component_id="` + TYPEINSTANCE + TYPEDELIMITER + in + `"].public_ip)`
}

func templInstanceElasticIP(in string) string {
	return `$(components.#[_component_id="` + TYPEINSTANCE + TYPEDELIMITER + in + `"].elastic_ip)`
}

func templInstancePrivateIP(in string) string {
	return `$(components.#[_component_id="` + TYPEINSTANCE + TYPEDELIMITER + in + `"].ip)`
}

func templELBDNS(elb string) string {
	return `$(components.#[_component_id="` + TYPEELB + TYPEDELIMITER + elb + `"].dns_name)`
}

func templRDSClusterDNS(rds string) string {
	return `$(components.#[_component_id="` + TYPERDSCLUSTER + TYPEDELIMITER + `"].endpoint)`
}

func templRDSInstanceDNS(rds string) string {
	return `$(components.#[_component_id="` + TYPERDSINSTANCE + TYPEDELIMITER + `"].endpoint)`
}
