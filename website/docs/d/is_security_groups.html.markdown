---
layout: "ibm"
page_title: "IBM : ibm_is_security_groups"
description: |-
  Get information about SecurityGroupCollection
subcategory: "VPC infrastructure"
---

# ibm_is_security_groups

Provides a read-only data source for SecurityGroupCollection. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.
For more information, about security group, see API Docs(https://cloud.ibm.com/docs/vpc?topic=vpc-using-security-groups).

## Example Usage

```terraform
data "ibm_is_security_groups" "example" {
}
```

OR with Filters:

Filter with VPC name

```terraform
data "ibm_is_security_groups" "example" {
  vpc_name = ibm_is_vpc.example.name
}
```

Filter with VPC ID

```terraform
data "ibm_is_security_groups" "example" {
  vpc_id = ibm_is_vpc.example.id
}
```

Filter with VPC CRN
```terraform
data "ibm_is_security_groups" "example" {
  vpc_crn = ibm_is_vpc.example.crn
}
```

Filter with Resource Group ID

```terraform
data "ibm_is_security_groups" "example" {
  resource_group= data.ibm_resource_group.default.id
}
```


## Attribute Reference

In addition to all argument references listed, you can access the following attribute references after your data source is created.

- `id` - The unique identifier of the SecurityGroupCollection.
- `vpc_name` - Filters the collection to resources in the VPC with the exact specified name
- `vpc_id` - Filters the collection to resources in the VPC with the specified identifier
- `vpc_crn` - Filters the collection to resources in the VPC with the specified CRN
- `resource_group` -  Filters the collection to resources in the resource group with the specified identifier
- `security_groups` - (List) Collection of security groups.
	
	Nested scheme for `security_groups`:
	- `access_tags`  - (List) Access management tags associated for the security group.
	- `created_at` - (String) The date and time that this security group was created.
	- `crn` - (String) The security group's CRN.
	- `href` - (String) The security group's canonical URL.
	- `id` - (String) The unique identifier for this security group.
	- `name` - (String) The user-defined name for this security group. Names must be unique within the VPC the security group resides in.
	- `resource_group` - (List) The resource group object, for this security group.
		
		Nested scheme for `resource_group`:
		- `href` - (String) The URL for this resource group.
		- `id` - (String) The unique identifier for this resource group.
		- `name` - (String) The user-defined name for this resource group.
	- `rules` - (List) The rules for this security group. If no rules exist, all traffic will be denied.
		
		Nested scheme for `rules`:
		- `code` - (Integer) The ICMP traffic code to allow.
		- `direction` - (String) The direction of traffic to enforce, either `inbound` or `outbound`.
		- `href` - (String) The URL for this security group rule.
		- `id` - (String) The unique identifier for this security group rule.
		- `name` - (String) The name for this security group rule. The name must not be used by another rule in the security group.
		- `ip_version` - (String) The IP version to enforce. The format of `local.address`, `local.cidr_block`, `remote.address`, or `remote.cidr_block` must match this property, if they are used. Alternatively, if `remote` references a security group, then this rule only applies to IP addresses in that group matching this IP version. Supported values are `ipv4` and `ipv6`.
		- `local` - (List) The local IP address or range of local IP addresses to which this rule will allow inbound traffic (or from which, for outbound traffic). A CIDR block of `0.0.0.0/0` allows traffic to all local IPv4 addresses; `::/0` allows traffic to all local IPv6 addresses.

			Nested scheme for `local`:
			- `address` - (String) The IP address.
			- `cidr_block` - (String) The CIDR block.
		- `port_max` - (Integer) The inclusive upper bound of TCP/UDP port range.
		- `port_min` - (Integer) The inclusive lower bound of TCP/UDP port range.
		- `protocol` - (String) The protocol to enforce. Supported values for IPv4 rules: `icmp`, `tcp`, `udp`, `any`, `icmp_tcp_udp`, and others. Supported additional values for IPv6 rules: `ipv6_icmp` (protocol 58), `ipv6_hop_opt` (protocol 0), `ipv6_route` (protocol 43), `ipv6_frag` (protocol 44), `ipv6_no_next` (protocol 59), `ipv6_dest_opts` (protocol 60), `ipv6_mobility` (protocol 135).
		- `remote` - (List) The IP addresses or security groups from which this rule allows traffic (or to which, for outbound rules). Can be specified as an IP address, a CIDR block, or a security group. A CIDR block of `0.0.0.0/0` allows traffic from any IPv4 source; `::/0` allows traffic from any IPv6 source.

			Nested scheme for `remote`:
			- `address` - (String) The IP address.
			- `cidr_block` - (String) The CIDR block.
			- `crn` - (String) The security group's CRN.
			- `deleted` - (List) If present, this property indicates the referenced resource has been deleted and provides some supplementary information.

			Nested scheme for `deleted`:
			- `more_info` - (String) Link to documentation about deleted resources.
			- `href` - (String) The security group's canonical URL.
			- `id` - (String) The unique identifier for this security group.
			- `name` - (String) The user-defined name for this security group. Names must be unique within the VPC the security group resides in.
		- `type` - (Integer) The ICMP traffic type to allow.
	- `targets` - (List) The targets for this security group.

	Nested scheme for `targets`:
		- `crn` - (String) The load balancer's CRN.
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted and provides some supplementary information.

	Nested scheme for `deleted`:
		- `more_info` - (String) Link to documentation about deleted resources.
		- `href` - (String) The URL for this network interface.
		- `id` - (String) The unique identifier for this network interface.
		- `name` - (String) The user-defined name for this network interface.
		- `resource_type` - (String) The resource type.
	- `vpc` - (List) The VPC this security group is a part of.

	Nested scheme for `vpc`:
		- `crn` - (String) The CRN for this VPC.
		- `deleted` - (List) If present, this property indicates the referenced resource has been deleted and provides some supplementary information.

	Nested scheme for `deleted`:
		- `more_info` - (String) Link to documentation about deleted resources.
		- `href` - (String) The URL for this VPC.
		- `id` - (String) The unique identifier for this VPC.
		- `name` - (String) The unique user-defined name for this VPC.
