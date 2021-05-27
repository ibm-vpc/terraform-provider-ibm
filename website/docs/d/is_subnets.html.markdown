---
subcategory: "VPC infrastructure"
layout: "ibm"
page_title: "IBM : Subnets"
description: |-
  Manages IBM Cloud Infrastructure Subnets.
---

# ibm\_is_subnets

Import the details of an existing IBM Cloud Infrastructure subnets as a read-only data source. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax.


## Example Usage

```terraform

data "ibm_is_subnets" "ds_subnets" {
}

```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `subnets` - List of all subnets in the IBM Cloud Infrastructure.
  * `name` - The name for this subnet.
  * `id` - The unique identifier for this subnet.
  * `ipv4_cidr_block` - The IPv4 CIDR block for this subnet.
  * `ipv6_cidr_block` - The IPv6 CIDR block for this subnet when used.
  * `status` - The status of this subnet.
  * `crn` - The CRN for this image.
  * `available_ipv4_address_count` - Amount of addresses available within this subnet.
  * `total_ipv4_address_count` - Amount of addresses used within this subnet.
  * `network_acl` - Security group attached to this subnet.
  * `public_gateway` - Public gateway attached to this subnet.
  * `resource_group` - Resource group where this subnet is created.
  * `vpc` - VPC where this subnet is created.
  * `zone` - Zone where this subnet is created.
